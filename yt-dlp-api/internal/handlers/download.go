package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
