package handlers

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ytdl-api/internal/downloads"
	"ytdl-api/internal/models"
	"ytdl-api/internal/repository/db"

	"github.com/jackc/pgx/v5"
)

const (
	sseRetry        = 2 * time.Second
	ssePingInterval = 15 * time.Second
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

	req := &models.DownloadPlaylistRequest{}
	opts := decodeFilterOptions(r, req)

	_, err = h.Manager.StartPlaylist(r.Context(), id, opts, req.RetryUnavailable)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, "playlist not found")
	case errors.Is(err, downloads.ErrRunInProgress):
		writeError(w, http.StatusConflict, "this playlist is already downloading")
	case errors.Is(err, downloads.ErrNothingToDownload):
		writeError(w, http.StatusConflict, "nothing left to download in this playlist")
	case err != nil:
		slog.Error("Failed to start playlist download", "playlist_id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to start playlist download")
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *DownloadHandler) DownloadEntry(w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.ParseInt(r.PathValue("entryId"), 10, 64)
	if err != nil || entryID < 1 {
		writeError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	opts := decodeFilterOptions(r, &models.DownloadEntryRequest{})

	err = h.Manager.StartEntry(r.Context(), entryID, opts)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, "playlist entry not found")
	case errors.Is(err, downloads.ErrAlreadyDownloading):
		writeError(w, http.StatusConflict, "this song is already downloading")
	case err != nil:
		slog.Error("Failed to start entry download", "entry_id", entryID, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to start entry download")
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *DownloadHandler) CancelPlaylistDownload(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	if _, err := h.Manager.CancelPlaylist(id); err != nil {
		writeError(w, http.StatusConflict, "nothing is downloading in this playlist")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DownloadHandler) CancelEntryDownload(w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.ParseInt(r.PathValue("entryId"), 10, 64)
	if err != nil || entryID < 1 {
		writeError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	if err := h.Manager.CancelEntry(entryID); err != nil {
		writeError(w, http.StatusConflict, "this song is not downloading")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	events := h.Manager.Hub.Subscribe(id)
	defer h.Manager.Hub.Unsubscribe(id, events)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "retry: %d\n\n", sseRetry.Milliseconds())
	flusher.Flush()

	slog.Debug("SSE subscriber connected", "playlist_id", id)

	ping := time.NewTicker(ssePingInterval)
	defer ping.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			slog.Debug("SSE subscriber disconnected", "playlist_id", id)
			return
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
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

	if entry.FilePath.String == "" {
		writeError(w, http.StatusUnprocessableEntity, "entry has no downloaded file")
		return
	}

	if _, err := os.Stat(entry.FilePath.String); errors.Is(err, fs.ErrNotExist) {
		h.Manager.MarkFileMissing(entry)
		writeError(w, http.StatusGone, "file was removed from the server, download it again")
		return
	}

	filePath := h.servedFile(r.Context(), entry)
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("Failed to open file", "file_path", filePath, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer file.Close()

	serveAttachment(w, r, file, entryFilename(entry))
}

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

	type zipItem struct {
		entry db.PlaylistEntry
		name  string
	}

	width := max(len(strconv.Itoa(len(entries))), 2)
	var items []zipItem
	missing := 0
	for i, entry := range entries {
		if entry.Status.EntryStatus != db.EntryStatusCompleted || entry.FilePath.String == "" {
			continue
		}
		if _, err := os.Stat(entry.FilePath.String); errors.Is(err, fs.ErrNotExist) {
			h.Manager.MarkFileMissing(entry)
			missing++
			continue
		}
		items = append(items, zipItem{
			entry: entry,
			name:  fmt.Sprintf("%0*d - %s", width, i+1, entryFilename(entry)),
		})
	}

	if len(items) == 0 {
		if missing > 0 {
			writeError(w, http.StatusGone, "the downloaded files were removed from the server, download them again")
			return
		}
		writeError(w, http.StatusNotFound, "no downloaded files available for this playlist")
		return
	}

	zipName := sanitizeFilename(playlist.Title.String) + ".zip"
	if zipName == ".zip" {
		zipName = fmt.Sprintf("playlist-%d.zip", id)
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", contentDisposition(zipName))

	zw := zip.NewWriter(w)
	for _, item := range items {
		if err := addFileToZip(zw, h.servedFile(r.Context(), item.entry), item.name); err != nil {
			slog.Error("Failed to add file to zip", "entry_id", item.entry.ID, "error", err)
		}
	}
	if err := zw.Close(); err != nil {
		slog.Error("Failed to close zip writer", "error", err)
	}
}

func (h *DownloadHandler) servedFile(ctx context.Context, entry db.PlaylistEntry) string {
	path, err := h.Manager.TaggedFile(ctx, entry)
	if err != nil {
		if ctx.Err() == nil {
			slog.Error("Failed to tag file, serving it untagged", "entry_id", entry.ID, "error", err)
		}
		return entry.FilePath.String
	}
	return path
}

func serveAttachment(w http.ResponseWriter, r *http.Request, file *os.File, filename string) {
	info, err := file.Stat()
	if err != nil {
		slog.Error("Failed to stat file", "file_path", file.Name(), "error", err)
		writeError(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	w.Header().Set("Content-Disposition", contentDisposition(filename))
	if contentType := mime.TypeByExtension(filepath.Ext(filename)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	http.ServeContent(w, r, filename, info.ModTime(), file)
}

func contentDisposition(filename string) string {
	encoded := strings.ReplaceAll(url.QueryEscape(filename), "+", "%20")
	return fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, asciiFilename(filename), encoded)
}

func asciiFilename(name string) string {
	ext := filepath.Ext(name)
	base := strings.Map(func(r rune) rune {
		if r < 0x20 || r > 0x7e || r == '"' || r == '\\' || r == '%' {
			return -1
		}
		return r
	}, strings.TrimSuffix(name, ext))

	base = strings.Join(strings.Fields(base), " ")
	if base == "" {
		base = "download"
	}
	return base + ext
}

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
	header.Method = zip.Store

	entry, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, file)
	return err
}

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
