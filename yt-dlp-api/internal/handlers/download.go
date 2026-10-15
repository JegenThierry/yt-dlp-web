package handlers

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ytdl-api/internal/downloads"
	"ytdl-api/internal/models"
	"ytdl-api/internal/repository/db"

	"github.com/jackc/pgx/v5"
)

type DownloadHandler struct {
	Queries *db.Queries
	Manager *downloads.Manager
}

func (h *DownloadHandler) DownloadPlaylist(w http.ResponseWriter, r *http.Request) {
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

	opts := decodeFilterOptions(r, &models.DownloadPlaylistRequest{})

	slog.Info("Starting playlist download", "playlist_id", id, "entries", len(entries))

	go h.Manager.DownloadPlaylist(context.Background(), id, opts)

	writeJSON(w, http.StatusAccepted, models.Playlist{
		ID:         playlist.ID,
		URL:        playlist.Url,
		Title:      playlist.Title.String,
		Entries:    toPlaylistEntryModels(entries),
		EntryCount: int64(len(entries)),
		CreatedAt:  playlist.CreatedAt.Time,
		UpdatedAt:  playlist.UpdatedAt.Time,
	})
}

func (h *DownloadHandler) DownloadEntry(w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.ParseInt(r.PathValue("entryId"), 10, 64)
	if err != nil || entryID < 1 {
		writeError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	entry, err := h.Queries.FetchPlaylistEntryByID(r.Context(), entryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "playlist entry not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist entry")
		return
	}

	opts := decodeFilterOptions(r, &models.DownloadEntryRequest{})

	slog.Info("Starting entry download", "entry_id", entryID, "video_id", entry.VideoID)

	go h.Manager.DownloadEntry(context.Background(), entryID, opts)

	writeJSON(w, http.StatusAccepted, models.PlaylistEntry{
		ID:         entry.ID,
		PlaylistID: entry.PlaylistID,
		VideoID:    entry.VideoID,
		Title:      entry.Title.String,
		Artist:     entry.Artist.String,
		Duration:   int64(entry.Duration.Int32),
		Thumbnail:  entry.Thumbnail.String,
		Status:     models.DownloadStatus(entry.Status.EntryStatus),
		FilePath:   entry.FilePath.String,
		Format:     models.AudioFormat(entry.Format.MediaFormat),
		CreatedAt:  entry.CreatedAt.Time,
		UpdatedAt:  entry.UpdatedAt.Time,
	})
}

func (h *DownloadHandler) GetPlaylistEvents(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	events := h.Manager.Hub.Subscribe(id)
	defer h.Manager.Hub.Unsubscribe(id, events)

	slog.Debug("SSE subscriber connected", "playlist_id", id)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			slog.Debug("SSE subscriber disconnected", "playlist_id", id)
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			payload, err := json.Marshal(ev)
			if err != nil {
				slog.Error("Failed to marshal progress event", "error", err)
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}

type filterOptionsRequest interface {
	FilterOptions() *models.FilterOptions
}

func decodeFilterOptions(r *http.Request, req filterOptionsRequest) models.FilterOptions {
	if r.Body == nil {
		return models.FilterOptions{}
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return models.FilterOptions{}
	}

	if f := req.FilterOptions(); f != nil {
		return *f
	}

	return models.FilterOptions{}
}

// DownloadEntryFile handles GET /api/playlists/entries/{entryId}/file. It
// streams the entry's cached audio file to the browser as an attachment.
func (h *DownloadHandler) DownloadEntryFile(w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.ParseInt(r.PathValue("entryId"), 10, 64)
	if err != nil || entryID < 1 {
		writeError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	entry, err := h.Queries.FetchPlaylistEntryByID(r.Context(), entryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "playlist entry not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist entry")
		return
	}

	filePath := entry.FilePath.String
	if filePath == "" {
		writeError(w, http.StatusUnprocessableEntity, "entry has no downloaded file")
		return
	}

	h.serveAttachment(w, r, filePath, entryFilename(entry))
}

// DownloadPlaylistZip handles GET /api/playlists/{id}/download/zip. It streams
// a zip archive of every completed, on-disk entry file to the browser.
func (h *DownloadHandler) DownloadPlaylistZip(w http.ResponseWriter, r *http.Request) {
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

	var files []string
	for _, entry := range entries {
		if entry.Status.EntryStatus != db.EntryStatusCompleted || entry.FilePath.String == "" {
			continue
		}
		files = append(files, entry.FilePath.String)
	}

	if len(files) == 0 {
		writeError(w, http.StatusNotFound, "no downloaded files available for this playlist")
		return
	}

	zipName := sanitizeFilename(playlist.Title.String) + ".zip"
	if zipName == ".zip" {
		zipName = fmt.Sprintf("playlist-%d.zip", id)
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipName))

	zw := zip.NewWriter(w)
	for _, entry := range entries {
		if entry.Status.EntryStatus != db.EntryStatusCompleted || entry.FilePath.String == "" {
			continue
		}
		if err := addFileToZip(zw, entry.FilePath.String, entryFilename(entry)); err != nil {
			slog.Error("Failed to add file to zip", "entry_id", entry.ID, "error", err)
		}
	}
	if err := zw.Close(); err != nil {
		slog.Error("Failed to close zip writer", "error", err)
	}
}

// serveAttachment streams a single file to the browser as an attachment under a
// human-readable filename, honoring range requests through http.ServeContent.
// The on-disk path may use a hash name; filename is the user-facing name.
func (h *DownloadHandler) serveAttachment(w http.ResponseWriter, r *http.Request, filePath, filename string) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			writeError(w, http.StatusNotFound, "file no longer exists on server")
			return
		}
		slog.Error("Failed to open file", "file_path", filePath, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		slog.Error("Failed to stat file", "file_path", filePath, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	if contentType := mime.TypeByExtension(filepath.Ext(filename)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	http.ServeContent(w, r, filename, info.ModTime(), file)
}

// addFileToZip appends a single file as an entry to the zip writer, using the
// human-readable name for the archive entry.
func addFileToZip(zw *zip.Writer, filePath, name string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = name
	header.Method = zip.Deflate

	entry, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, file)
	return err
}

// entryFilename returns the user-facing filename for an entry's audio file: the
// "author - title" form with the file extension, falling back to the base name
// when metadata is missing.
func entryFilename(entry db.PlaylistEntry) string {
	ext := filepath.Ext(entry.FilePath.String)
	base := ""
	switch {
	case entry.Artist.String != "" && entry.Title.String != "":
		base = fmt.Sprintf("%s - %s", entry.Artist.String, entry.Title.String)
	case entry.Title.String != "":
		base = entry.Title.String
	default:
		return sanitizeFilename(filepath.Base(entry.FilePath.String))
	}
	return sanitizeFilename(base) + ext
}

// sanitizeFilename strips characters that are unsafe for HTTP content
// disposition filenames and trims surrounding whitespace.
func sanitizeFilename(name string) string {
	name = strings.Map(func(r rune) rune {
		switch r {
		case '"', '\\', '/', ':', '*', '?', '<', '>', '|':
			return '_'
		default:
			return r
		}
	}, name)
	return strings.TrimSpace(name)
}
