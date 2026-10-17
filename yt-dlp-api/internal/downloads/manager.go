package downloads

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ytdl-api/internal/models"
	"ytdl-api/internal/repository/db"
	"ytdl-api/pkg/ytdlp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrAlreadyDownloading = errors.New("entry is already downloading")
	ErrRunInProgress      = errors.New("playlist download is already running")
	ErrNothingToDownload  = errors.New("nothing left to download")
	ErrNotDownloading     = errors.New("nothing is downloading")
	ErrPlaylistBusy       = errors.New("playlist is downloading")
)

const progressInterval = time.Second

type job struct {
	entry   db.PlaylistEntry
	opts    ytdlp.FilterOptions
	ctx     context.Context
	cancel  context.CancelFunc
	started bool
}

type run struct {
	cancel context.CancelFunc
}

type Manager struct {
	Queries            *db.Queries
	DownloadDir        string
	Hub                *Hub
	CookiesFromBrowser string
	CookiesFile        string

	partialDir string
	taggedDir  string

	mu     sync.Mutex
	jobs   map[int64]*job
	runs   map[int64]*run
	hashes map[string]chan struct{}
}

func NewManager(queries *db.Queries, downloadDir, cookiesFromBrowser, cookiesFile string) *Manager {
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		slog.Error("Failed to create download directory", "dir", downloadDir, "error", err)
	}

	partialDir := absPath(filepath.Join(downloadDir, ".partial"))
	taggedDir := absPath(filepath.Join(downloadDir, "tagged"))

	return &Manager{
		Queries:            queries,
		DownloadDir:        downloadDir,
		Hub:                NewHub(),
		CookiesFromBrowser: cookiesFromBrowser,
		CookiesFile:        cookiesFile,
		partialDir:         partialDir,
		taggedDir:          taggedDir,
		jobs:               make(map[int64]*job),
		runs:               make(map[int64]*run),
		hashes:             make(map[string]chan struct{}),
	}
}

func absPath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return path
}

func (m *Manager) ResetInterrupted(ctx context.Context) error {
	if err := os.RemoveAll(m.partialDir); err != nil {
		return fmt.Errorf("remove partial downloads: %w", err)
	}

	reset, err := m.Queries.ResetInterruptedEntries(ctx)
	if err != nil {
		return fmt.Errorf("reset interrupted entries: %w", err)
	}
	if reset > 0 {
		slog.Info("Reset interrupted downloads", "entries", reset)
	}
	return nil
}

func (m *Manager) StartEntry(ctx context.Context, entryID int64, opts models.FilterOptions) error {
	entry, err := m.Queries.FetchPlaylistEntryByID(ctx, entryID)
	if err != nil {
		return err
	}

	playlist, err := m.Queries.FetchPlaylistByID(ctx, entry.PlaylistID)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, busy := m.jobs[entry.ID]; busy {
		return ErrAlreadyDownloading
	}

	j := m.addJob(context.Background(), entry, m.toYtDlpOptions(entry, playlist, opts))
	m.setStatus(entry, models.StatusPending, "")

	slog.Info("Starting entry download", "entry_id", entry.ID, "video_id", entry.VideoID)
	go m.run(j)
	return nil
}

func (m *Manager) StartPlaylist(ctx context.Context, playlistID int64, opts models.FilterOptions, retryUnavailable bool) (int, error) {
	playlist, err := m.Queries.FetchPlaylistByID(ctx, playlistID)
	if err != nil {
		return 0, err
	}

	entries, err := m.Queries.FetchPlaylistEntries(ctx, playlistID)
	if err != nil {
		return 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, running := m.runs[playlistID]; running {
		return 0, ErrRunInProgress
	}

	runCtx, cancelRun := context.WithCancel(context.Background())
	var queue []*job
	for _, entry := range entries {
		if entry.Status.EntryStatus == db.EntryStatusCompleted {
			continue
		}
		if !retryUnavailable && isUnavailable(entry) {
			continue
		}
		if _, busy := m.jobs[entry.ID]; busy {
			continue
		}
		queue = append(queue, m.addJob(runCtx, entry, m.toYtDlpOptions(entry, playlist, opts)))
	}

	if len(queue) == 0 {
		cancelRun()
		return 0, ErrNothingToDownload
	}

	r := &run{cancel: cancelRun}
	m.runs[playlistID] = r
	m.setStatuses(playlistID, queue, models.StatusPending)

	slog.Info("Starting playlist download", "playlist_id", playlistID, "queued", len(queue))
	go func() {
		defer func() {
			m.mu.Lock()
			if m.runs[playlistID] == r {
				delete(m.runs, playlistID)
			}
			m.mu.Unlock()
			cancelRun()
		}()

		for _, j := range queue {
			m.run(j)
		}
	}()

	return len(queue), nil
}

func (m *Manager) CancelEntry(entryID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	j, ok := m.jobs[entryID]
	if !ok {
		return ErrNotDownloading
	}

	j.cancel()
	if !j.started {
		delete(m.jobs, entryID)
		m.setStatus(j.entry, models.StatusNew, "")
	}

	slog.Info("Cancelled entry download", "entry_id", entryID)
	return nil
}

func (m *Manager) CancelPlaylist(playlistID int64) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r, ok := m.runs[playlistID]; ok {
		r.cancel()
		delete(m.runs, playlistID)
	}

	cancelled := 0
	var queued []*job
	for id, j := range m.jobs {
		if j.entry.PlaylistID != playlistID {
			continue
		}
		j.cancel()
		cancelled++
		if !j.started {
			delete(m.jobs, id)
			queued = append(queued, j)
		}
	}

	if cancelled == 0 {
		return 0, ErrNotDownloading
	}

	if len(queued) > 0 {
		m.setStatuses(playlistID, queued, models.StatusNew)
	}

	slog.Info("Cancelled playlist download", "playlist_id", playlistID, "entries", cancelled)
	return cancelled, nil
}

func (m *Manager) DeletePlaylist(ctx context.Context, playlistID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, running := m.runs[playlistID]; running {
		return ErrPlaylistBusy
	}
	for _, j := range m.jobs {
		if j.entry.PlaylistID == playlistID {
			return ErrPlaylistBusy
		}
	}

	deleted, err := m.Queries.DeletePlaylist(ctx, playlistID)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return pgx.ErrNoRows
	}

	slog.Info("Deleted playlist", "playlist_id", playlistID)
	return nil
}

func (m *Manager) MarkFileMissing(entry db.PlaylistEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, busy := m.jobs[entry.ID]; busy {
		return
	}

	slog.Warn("Downloaded file is missing", "entry_id", entry.ID, "file_path", entry.FilePath.String)
	m.setStatus(entry, models.StatusNew, "")
}

func isUnavailable(entry db.PlaylistEntry) bool {
	return entry.Status.EntryStatus == db.EntryStatusFailed &&
		entry.ErrorCode.EntryErrorCode == db.EntryErrorCodeUnavailable
}

func (m *Manager) addJob(parent context.Context, entry db.PlaylistEntry, opts ytdlp.FilterOptions) *job {
	ctx, cancel := context.WithCancel(parent)
	j := &job{entry: entry, opts: opts, ctx: ctx, cancel: cancel}
	m.jobs[entry.ID] = j
	return j
}

func (m *Manager) run(j *job) {
	defer m.release(j)

	m.mu.Lock()
	if j.ctx.Err() != nil {
		m.mu.Unlock()
		return
	}
	j.started = true
	m.mu.Unlock()

	m.download(j)
}

func (m *Manager) release(j *job) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.jobs[j.entry.ID] == j {
		delete(m.jobs, j.entry.ID)
	}
	j.cancel()
}

func (m *Manager) download(j *job) {
	entry := j.entry
	hash := ytdlp.ContentHash(entry.VideoID, j.opts)

	unlock, err := m.lockHash(j.ctx, hash)
	if err != nil {
		slog.Info("Download cancelled", "entry_id", entry.ID, "video_id", entry.VideoID)
		m.setStatus(entry, models.StatusNew, "")
		return
	}
	defer unlock()

	expectedPath := filepath.Join(m.DownloadDir, hash+"."+expectedExtension(j.opts.Format))
	if _, err := os.Stat(expectedPath); err == nil {
		slog.Info("Reusing cached file", "entry_id", entry.ID, "video_id", entry.VideoID, "file_path", expectedPath)
		m.setStatus(entry, models.StatusCompleted, expectedPath)
		return
	}

	m.setStatus(entry, models.StatusDownloading, "")

	tempDir := filepath.Join(m.partialDir, hash)
	defer os.RemoveAll(tempDir)

	filePath, err := ytdlp.Download(j.ctx, entry.VideoID, m.DownloadDir, tempDir, hash, j.opts, m.progressReporter(entry))
	switch {
	case j.ctx.Err() != nil:
		slog.Info("Download cancelled", "entry_id", entry.ID, "video_id", entry.VideoID)
		m.setStatus(entry, models.StatusNew, "")
	case err != nil:
		code, message := failureReason(err)
		slog.Error("Download failed", "entry_id", entry.ID, "video_id", entry.VideoID, "error_code", code, "error", err)
		m.setFailed(entry, code, message)
	default:
		slog.Info("Download completed", "entry_id", entry.ID, "file_path", filePath)
		m.setStatus(entry, models.StatusCompleted, filePath)
	}
}

func (m *Manager) lockHash(ctx context.Context, hash string) (func(), error) {
	for {
		m.mu.Lock()
		busy, taken := m.hashes[hash]
		if !taken {
			done := make(chan struct{})
			m.hashes[hash] = done
			m.mu.Unlock()
			return func() {
				m.mu.Lock()
				delete(m.hashes, hash)
				m.mu.Unlock()
				close(done)
			}, nil
		}
		m.mu.Unlock()

		select {
		case <-busy:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (m *Manager) progressReporter(entry db.PlaylistEntry) func(ytdlp.Progress) {
	var last time.Time
	var lastPhase ytdlp.Phase
	return func(p ytdlp.Progress) {
		now := time.Now()
		if p.Phase == lastPhase && now.Sub(last) < progressInterval {
			return
		}
		last, lastPhase = now, p.Phase

		m.Hub.Broadcast(entry.PlaylistID, ProgressEvent{
			PlaylistID: entry.PlaylistID,
			EntryID:    entry.ID,
			Status:     models.StatusDownloading,
			Phase:      models.DownloadPhase(p.Phase),
			Percent:    p.Percent,
			ETA:        p.ETA,
		})
	}
}

func expectedExtension(format ytdlp.AudioCodecs) string {
	switch format {
	case ytdlp.MP3:
		return "mp3"
	case ytdlp.M4A, ytdlp.AAC, ytdlp.ALAC:
		return "m4a"
	case ytdlp.OPUS:
		return "opus"
	case ytdlp.VORBIS:
		return "ogg"
	case ytdlp.FLAC:
		return "flac"
	case ytdlp.WAV:
		return "wav"
	default:
		return string(format)
	}
}

func failureReason(err error) (models.EntryErrorCode, string) {
	var downloadErr *ytdlp.DownloadError
	if errors.As(err, &downloadErr) {
		return models.EntryErrorCode(downloadErr.Code), downloadErr.Message
	}
	return models.ErrorUnknown, err.Error()
}

func (m *Manager) setStatus(entry db.PlaylistEntry, status models.DownloadStatus, filePath string) {
	m.saveStatus(entry, ProgressEvent{Status: status, FilePath: filePath})
}

func (m *Manager) setFailed(entry db.PlaylistEntry, code models.EntryErrorCode, message string) {
	m.saveStatus(entry, ProgressEvent{Status: models.StatusFailed, ErrorCode: code, ErrorMessage: message})
}

func (m *Manager) saveStatus(entry db.PlaylistEntry, ev ProgressEvent) {
	ev.PlaylistID = entry.PlaylistID
	ev.EntryID = entry.ID

	_, err := m.Queries.UpdatePlaylistEntryStatus(context.Background(), db.UpdatePlaylistEntryStatusParams{
		ID:           entry.ID,
		Status:       db.NullEntryStatus{EntryStatus: db.EntryStatus(ev.Status), Valid: true},
		FilePath:     pgtype.Text{String: ev.FilePath, Valid: ev.FilePath != ""},
		ErrorMessage: pgtype.Text{String: ev.ErrorMessage, Valid: ev.ErrorMessage != ""},
		ErrorCode:    db.NullEntryErrorCode{EntryErrorCode: db.EntryErrorCode(ev.ErrorCode), Valid: ev.ErrorCode != ""},
	})
	if err != nil {
		slog.Error("Failed to update playlist entry status", "entry_id", entry.ID, "error", err)
	}

	m.Hub.Broadcast(entry.PlaylistID, ev)
}

func (m *Manager) setStatuses(playlistID int64, jobs []*job, status models.DownloadStatus) {
	ids := make([]int64, len(jobs))
	for i, j := range jobs {
		ids[i] = j.entry.ID
	}

	err := m.Queries.SetPlaylistEntriesStatus(context.Background(), db.SetPlaylistEntriesStatusParams{
		Status: db.NullEntryStatus{EntryStatus: db.EntryStatus(status), Valid: true},
		Ids:    ids,
	})
	if err != nil {
		slog.Error("Failed to update playlist entry statuses", "playlist_id", playlistID, "error", err)
	}

	for _, id := range ids {
		m.Hub.Broadcast(playlistID, ProgressEvent{
			PlaylistID: playlistID,
			EntryID:    id,
			Status:     status,
		})
	}
}

func (m *Manager) toYtDlpOptions(entry db.PlaylistEntry, playlist db.Playlist, opts models.FilterOptions) ytdlp.FilterOptions {
	format := ytdlp.AudioCodecs(opts.Format)
	if format == "" {
		format = ytdlp.AudioCodecs(entry.Format.MediaFormat)
	}
	if format == "" {
		format = ytdlp.MP3
	}

	quality := ytdlp.AudioQuality(playlist.Quality)
	if opts.Quality.Valid() {
		quality = ytdlp.AudioQuality(opts.Quality)
	}
	if quality == "" {
		quality = ytdlp.QualityBest
	}

	return ytdlp.FilterOptions{
		Format:  format,
		Quality: quality,
		Cookies: ytdlp.Cookies{File: m.CookiesFile, FromBrowser: m.CookiesFromBrowser},
	}
}
