package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ytdl-api/internal/downloads"
	"ytdl-api/internal/models"
	"ytdl-api/internal/repository/db"
	"ytdl-api/pkg/ytdlp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const playlistReadTimeout = 5 * time.Minute

type PlaylistHandler struct {
	Queries *db.Queries
	Pool    *pgxpool.Pool
	Manager *downloads.Manager
	Cookies ytdlp.Cookies
}

func (h *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	playlist, err := h.Queries.FetchPlaylistByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "playlist not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist")
		return
	}

	entries, err := h.Queries.FetchPlaylistEntries(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist entries")
		return
	}

	slog.Debug("Fetched playlist", "id", id, "entries", len(entries))

	writeJSON(w, http.StatusOK, toPlaylistModel(playlist, entries))
}

func (h *PlaylistHandler) FetchPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.Queries.FetchAllPlaylists(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlists")
		return
	}

	playlistModels := make([]models.Playlist, 0, len(playlists))
	for _, p := range playlists {
		playlistModels = append(playlistModels, models.Playlist{
			ID:               p.ID,
			URL:              p.Url,
			Title:            p.Title.String,
			Quality:          models.AudioQuality(p.Quality),
			EntryCount:       p.EntryCount,
			CompletedCount:   p.CompletedCount,
			FailedCount:      p.FailedCount,
			UnavailableCount: p.UnavailableCount,
			ActiveCount:      p.ActiveCount,
			TotalDuration:    p.TotalDuration,
			CreatedAt:        p.CreatedAt.Time,
			UpdatedAt:        p.UpdatedAt.Time,
			LastActivityAt:   p.LastActivityAt.Time,
		})
	}

	writeJSON(w, http.StatusOK, playlistModels)
}

func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	err = h.Manager.DeletePlaylist(r.Context(), id)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, "playlist not found")
	case errors.Is(err, downloads.ErrPlaylistBusy):
		writeError(w, http.StatusConflict, "This playlist is downloading. Cancel the download before deleting it.")
	case err != nil:
		slog.Error("Failed to delete playlist", "playlist_id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete playlist")
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	source, problem := validateSourceURL(req.URL)
	if problem != "" {
		writeError(w, http.StatusBadRequest, problem)
		return
	}
	req.URL = source

	if req.Quality == "" {
		req.Quality = models.QualityBest
	}
	if !req.Quality.Valid() {
		writeError(w, http.StatusBadRequest, "quality must be one of best, high, standard or low")
		return
	}
	if req.TargetFormat != "" && !req.TargetFormat.Valid() {
		writeError(w, http.StatusBadRequest, "target_format must be one of mp3, aac, flac, wav, m4a, opus, vorbis or alac")
		return
	}

	info, ok := h.readPlaylist(w, r, req.URL)
	if !ok {
		return
	}

	result, err := h.importNew(context.WithoutCancel(r.Context()), req, info)
	if err != nil {
		slog.Error("Failed to save playlist", "url", req.URL, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save playlist")
		return
	}

	slog.Info("Playlist imported", "id", result.ID, "existing", result.Existing, "added", result.AddedCount, "skipped", result.SkippedCount)

	status := http.StatusCreated
	if result.Existing {
		status = http.StatusOK
	}
	writeJSON(w, status, result)
}

func (h *PlaylistHandler) RefreshPlaylist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	playlist, err := h.Queries.FetchPlaylistByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "playlist not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist")
		return
	}

	info, ok := h.readPlaylist(w, r, playlist.Url)
	if !ok {
		return
	}

	result, err := h.importInto(context.WithoutCancel(r.Context()), playlist, info)
	if err != nil {
		slog.Error("Failed to refresh playlist", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save new songs")
		return
	}

	slog.Info("Playlist refreshed", "id", id, "added", result.AddedCount, "skipped", result.SkippedCount)
	writeJSON(w, http.StatusOK, result)
}

func validateSourceURL(raw string) (string, string) {
	source := strings.TrimSpace(raw)
	if source == "" {
		return "", "Paste a playlist or video link."
	}

	parsed, err := url.Parse(source)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", "Enter a full link that starts with https://"
	}
	return source, ""
}

func (h *PlaylistHandler) readPlaylist(w http.ResponseWriter, r *http.Request, source string) (*ytdlp.Playlist, bool) {
	ctx, cancel := context.WithTimeout(r.Context(), playlistReadTimeout)
	defer cancel()

	slog.Info("Fetching playlist metadata", "url", source)

	info, err := ytdlp.GetPlaylist(ctx, source, h.Cookies)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("Playlist read cancelled by the client", "url", source)
			return nil, false
		}
		status, message := playlistReadError(err)
		slog.Warn("Failed to read playlist", "url", source, "status", status, "error", err)
		writeError(w, status, message)
		return nil, false
	}

	if len(info.Videos) == 0 {
		writeError(w, http.StatusUnprocessableEntity, noVideosMessage(info.Skipped))
		return nil, false
	}

	return info, true
}

func playlistReadError(err error) (int, string) {
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout, "Reading this link took longer than 5 minutes, so it was stopped. Try again later."
	}

	var metadataErr *ytdlp.MetadataError
	if !errors.As(err, &metadataErr) {
		return http.StatusBadGateway, "Couldn't run yt-dlp on the server. Check the API log."
	}

	switch metadataErr.Kind {
	case ytdlp.MetadataUnsupported:
		return http.StatusUnprocessableEntity, "yt-dlp can't read this link. Paste the link of a playlist or video page."
	case ytdlp.MetadataNotFound:
		return http.StatusNotFound, "Nothing was found at this link. The playlist or video may be private, deleted or mistyped."
	case ytdlp.MetadataSignInRequired:
		return http.StatusForbidden, "This link needs a signed-in account. Set YT_COOKIES_FILE or YT_COOKIES_FROM_BROWSER for the API and try again."
	case ytdlp.MetadataBlocked:
		return http.StatusBadGateway, "YouTube is blocking or rate-limiting the server right now. Try again in a while."
	default:
		return http.StatusBadGateway, "yt-dlp couldn't read this link: " + metadataErr.Message
	}
}

func noVideosMessage(skipped int) string {
	switch skipped {
	case 0:
		return "This link has no downloadable videos."
	case 1:
		return "This link has no downloadable videos. 1 unavailable video was skipped."
	default:
		return fmt.Sprintf("This link has no downloadable videos. %d unavailable videos were skipped.", skipped)
	}
}

func (h *PlaylistHandler) importNew(ctx context.Context, req models.CreatePlaylistRequest, info *ytdlp.Playlist) (models.PlaylistImport, error) {
	return inTx(ctx, h, func(q *db.Queries) (models.PlaylistImport, error) {
		sourceID := pgtype.Text{String: info.SourceURL, Valid: info.SourceURL != ""}

		existing := false
		playlist, err := q.CreatePlaylist(ctx, db.CreatePlaylistParams{
			Url:      req.URL,
			Title:    pgtype.Text{String: info.Title, Valid: true},
			Quality:  db.AudioQuality(req.Quality),
			SourceID: sourceID,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			existing = true
			playlist, err = q.FetchPlaylistBySourceID(ctx, sourceID)
		}
		if err != nil {
			return models.PlaylistImport{}, fmt.Errorf("create playlist: %w", err)
		}

		format := db.NullMediaFormat{MediaFormat: db.MediaFormat(req.TargetFormat), Valid: req.TargetFormat != ""}
		if existing {
			format, err = playlistFormat(ctx, q, playlist.ID, format)
			if err != nil {
				return models.PlaylistImport{}, err
			}
		}

		return importEntries(ctx, q, playlist, format, info, existing)
	})
}

func (h *PlaylistHandler) importInto(ctx context.Context, playlist db.Playlist, info *ytdlp.Playlist) (models.PlaylistImport, error) {
	return inTx(ctx, h, func(q *db.Queries) (models.PlaylistImport, error) {
		if info.SourceURL != "" {
			err := q.SetPlaylistSourceID(ctx, db.SetPlaylistSourceIDParams{
				ID:       playlist.ID,
				SourceID: pgtype.Text{String: info.SourceURL, Valid: true},
			})
			if err != nil {
				return models.PlaylistImport{}, fmt.Errorf("set source id: %w", err)
			}
		}

		format, err := playlistFormat(ctx, q, playlist.ID, db.NullMediaFormat{})
		if err != nil {
			return models.PlaylistImport{}, err
		}

		return importEntries(ctx, q, playlist, format, info, true)
	})
}

func inTx[T any](ctx context.Context, h *PlaylistHandler, fn func(q *db.Queries) (T, error)) (T, error) {
	var zero T
	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		return zero, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := fn(h.Queries.WithTx(tx))
	if err != nil {
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return zero, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func playlistFormat(ctx context.Context, q *db.Queries, playlistID int64, fallback db.NullMediaFormat) (db.NullMediaFormat, error) {
	format, err := q.FetchPlaylistFormat(ctx, playlistID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fallback, nil
	}
	if err != nil {
		return db.NullMediaFormat{}, fmt.Errorf("fetch playlist format: %w", err)
	}
	return format, nil
}

func importEntries(ctx context.Context, q *db.Queries, playlist db.Playlist, format db.NullMediaFormat, info *ytdlp.Playlist, existing bool) (models.PlaylistImport, error) {
	var added int64
	for _, video := range info.Videos {
		inserted, err := q.CreatePlaylistEntry(ctx, db.CreatePlaylistEntryParams{
			PlaylistID: playlist.ID,
			VideoID:    video.ID,
			Title:      pgtype.Text{String: video.Title, Valid: true},
			Artist:     pgtype.Text{String: video.Artist, Valid: video.Artist != ""},
			Genre:      pgtype.Text{Valid: false},
			Year:       pgtype.Int4{Valid: false},
			Duration:   pgtype.Int4{Int32: video.Duration, Valid: true},
			Thumbnail:  pgtype.Text{String: video.Thumbnail, Valid: video.Thumbnail != ""},
			Status:     db.NullEntryStatus{EntryStatus: db.EntryStatusNew, Valid: true},
			FilePath:   pgtype.Text{Valid: false},
			Format:     format,
		})
		if err != nil {
			return models.PlaylistImport{}, fmt.Errorf("insert entry %s: %w", video.ID, err)
		}
		added += inserted
	}

	if existing && added > 0 {
		touched, err := q.TouchPlaylist(ctx, playlist.ID)
		if err != nil {
			return models.PlaylistImport{}, fmt.Errorf("touch playlist: %w", err)
		}
		playlist = touched
	}

	entries, err := q.FetchPlaylistEntries(ctx, playlist.ID)
	if err != nil {
		return models.PlaylistImport{}, fmt.Errorf("fetch entries: %w", err)
	}

	return models.PlaylistImport{
		Playlist:     toPlaylistModel(playlist, entries),
		Existing:     existing,
		AddedCount:   added,
		SkippedCount: int64(info.Skipped),
	}, nil
}

func toPlaylistModel(playlist db.Playlist, entries []db.PlaylistEntry) models.Playlist {
	model := models.Playlist{
		ID:             playlist.ID,
		URL:            playlist.Url,
		Title:          playlist.Title.String,
		Quality:        models.AudioQuality(playlist.Quality),
		Entries:        toPlaylistEntryModels(entries),
		EntryCount:     int64(len(entries)),
		CreatedAt:      playlist.CreatedAt.Time,
		UpdatedAt:      playlist.UpdatedAt.Time,
		LastActivityAt: playlist.UpdatedAt.Time,
	}

	for _, e := range entries {
		switch e.Status.EntryStatus {
		case db.EntryStatusCompleted:
			model.CompletedCount++
		case db.EntryStatusFailed:
			if e.ErrorCode.EntryErrorCode == db.EntryErrorCodeUnavailable {
				model.UnavailableCount++
			} else {
				model.FailedCount++
			}
		case db.EntryStatusPending, db.EntryStatusDownloading:
			model.ActiveCount++
		}
		model.TotalDuration += int64(e.Duration.Int32)
		if e.UpdatedAt.Time.After(model.LastActivityAt) {
			model.LastActivityAt = e.UpdatedAt.Time
		}
	}

	return model
}

const (
	minYear = 1900
	maxYear = 2100
)

var errEntryNotInPlaylist = errors.New("entry not in playlist")

func (h *PlaylistHandler) UpdatePlaylistMetadata(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	var updates []models.PlaylistEntryMetadataUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(updates) == 0 {
		writeError(w, http.StatusBadRequest, "No songs to update.")
		return
	}

	if fields := validateMetadata(updates); len(fields) > 0 {
		writeJSON(w, http.StatusBadRequest, models.MetadataValidationError{
			Error:  "Some changes aren't valid. Fix the marked fields and save again.",
			Fields: fields,
		})
		return
	}

	entries, err := h.updateMetadata(r.Context(), id, updates)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, "playlist not found")
	case errors.Is(err, errEntryNotInPlaylist):
		writeError(w, http.StatusNotFound, "A song you edited is no longer in this playlist. Reload the page and try again.")
	case err != nil:
		slog.Error("Failed to update metadata", "playlist_id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save metadata")
	default:
		slog.Info("Metadata updated", "playlist_id", id, "songs", len(entries))
		writeJSON(w, http.StatusOK, toPlaylistEntryModels(entries))
	}
}

func validateMetadata(updates []models.PlaylistEntryMetadataUpdate) []models.MetadataFieldError {
	var fields []models.MetadataFieldError
	seen := make(map[int64]bool, len(updates))
	for i := range updates {
		u := &updates[i]
		u.Title = strings.TrimSpace(u.Title)
		u.Artist = strings.TrimSpace(u.Artist)
		u.Genre = strings.TrimSpace(u.Genre)

		fail := func(field, message string) {
			fields = append(fields, models.MetadataFieldError{EntryID: u.ID, Field: field, Message: message})
		}
		switch {
		case u.ID < 1:
			fail("id", "Invalid song id.")
		case seen[u.ID]:
			fail("id", "This song is listed more than once.")
		}
		seen[u.ID] = true

		if u.Title == "" {
			fail("title", "Enter a title.")
		}
		switch {
		case u.Year == nil:
			fail("year", "Year is required. Send 0 to clear it.")
		case *u.Year != 0 && (*u.Year < minYear || *u.Year > maxYear):
			fail("year", fmt.Sprintf("Use a year from %d to %d.", minYear, maxYear))
		}
	}
	return fields
}

func (h *PlaylistHandler) updateMetadata(ctx context.Context, playlistID int64, updates []models.PlaylistEntryMetadataUpdate) ([]db.PlaylistEntry, error) {
	return inTx(ctx, h, func(q *db.Queries) ([]db.PlaylistEntry, error) {
		if _, err := q.FetchPlaylistByID(ctx, playlistID); err != nil {
			return nil, err
		}

		entries := make([]db.PlaylistEntry, 0, len(updates))
		for _, u := range updates {
			entry, err := q.UpdatePlaylistEntryMetadata(ctx, db.UpdatePlaylistEntryMetadataParams{
				ID:         u.ID,
				PlaylistID: playlistID,
				Title:      pgtype.Text{String: u.Title, Valid: true},
				Artist:     pgtype.Text{String: u.Artist, Valid: u.Artist != ""},
				Genre:      pgtype.Text{String: u.Genre, Valid: u.Genre != ""},
				Year:       pgtype.Int4{Int32: int32(*u.Year), Valid: *u.Year != 0},
			})
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("entry %d: %w", u.ID, errEntryNotInPlaylist)
			}
			if err != nil {
				return nil, fmt.Errorf("update entry %d: %w", u.ID, err)
			}
			entries = append(entries, entry)
		}
		return entries, nil
	})
}

func toPlaylistEntryModels(entries []db.PlaylistEntry) []models.PlaylistEntry {
	modelsEntries := make([]models.PlaylistEntry, 0, len(entries))
	for _, e := range entries {
		modelsEntries = append(modelsEntries, models.PlaylistEntry{
			ID:           e.ID,
			PlaylistID:   e.PlaylistID,
			VideoID:      e.VideoID,
			Title:        e.Title.String,
			Artist:       e.Artist.String,
			SourceTitle:  e.SourceTitle.String,
			SourceArtist: e.SourceArtist.String,
			Genre:        e.Genre.String,
			Year:         int64(e.Year.Int32),
			Duration:     int64(e.Duration.Int32),
			Thumbnail:    e.Thumbnail.String,
			Status:       models.DownloadStatus(e.Status.EntryStatus),
			FilePath:     e.FilePath.String,
			Format:       models.AudioFormat(e.Format.MediaFormat),
			ErrorMessage: e.ErrorMessage.String,
			ErrorCode:    models.EntryErrorCode(e.ErrorCode.EntryErrorCode),
			CreatedAt:    e.CreatedAt.Time,
			UpdatedAt:    e.UpdatedAt.Time,
		})
	}
	return modelsEntries
}
