package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"ytdl-api/internal/models"
	"ytdl-api/internal/repository/db"
	"ytdl-api/pkg/ytdlp"

	"github.com/jackc/pgx/v5/pgtype"
)

type PlaylistHandler struct {
	Queries *db.Queries
}

func (h *PlaylistHandler) FetchPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.Queries.FetchAllPlaylists(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlists")
		return
	}

	writeJSON(w, http.StatusOK, playlists)
}

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePlaylistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	slog.Info("Fetching playlist metadata for URL", "url", req.URL)

	metadata, err := ytdlp.GetPlaylist(req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist metadata")
		return
	}

	playlist, err := h.Queries.CreatePlaylist(r.Context(), db.CreatePlaylistParams{
		Title: pgtype.Text{String: metadata.Title, Valid: true},
		Url:   req.URL,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create playlist")
		return
	}

	for _, entry := range metadata.Entries {
		thumbnail := entry.GetThumbnail()
		_, err := h.Queries.CreatePlaylistEntry(r.Context(), db.CreatePlaylistEntryParams{
			PlaylistID: playlist.ID,
			VideoID:    entry.ID,
			Title:      pgtype.Text{String: entry.Title, Valid: true},
			Artist:     pgtype.Text{String: entry.Uploader, Valid: entry.Uploader != ""},
			Duration:   pgtype.Int4{Int32: int32(entry.Duration), Valid: true},
			Thumbnail:  pgtype.Text{String: thumbnail, Valid: thumbnail != ""},
			Status:     db.NullEntryStatus{EntryStatus: db.EntryStatusNew, Valid: true},
			FilePath:   pgtype.Text{Valid: false},
			Format:     db.NullMediaFormat{MediaFormat: db.MediaFormat(req.TargetFormat), Valid: req.TargetFormat != ""},
		})

		if err != nil {
			slog.Error("Failed to create playlist entry", "error", err)
			continue
		}
		slog.Info("Playlist entry created", "title", entry.Title, "url", entry.URL, "thumbnail", thumbnail)
	}

	writeJSON(w, http.StatusCreated, playlist)
}
