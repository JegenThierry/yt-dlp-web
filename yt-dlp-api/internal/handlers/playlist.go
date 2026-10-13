package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"ytdl-api/internal/models"
	"ytdl-api/internal/repository/db"
	"ytdl-api/pkg/ytdlp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PlaylistHandler struct {
	Queries *db.Queries
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

	writeJSON(w, http.StatusOK, models.Playlist{
		ID:         playlist.ID,
		URL:        playlist.Url,
		Title:      playlist.Title.String,
		Entries:    toPlaylistEntryModels(entries),
		EntryCount: int64(len(entries)),
		CreatedAt:  playlist.CreatedAt.Time,
		UpdatedAt:  playlist.UpdatedAt.Time,
	})
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
			ID:         p.ID,
			URL:        p.Url,
			Title:      p.Title.String,
			EntryCount: p.EntryCount,
			CreatedAt:  p.CreatedAt.Time,
			UpdatedAt:  p.UpdatedAt.Time,
		})
	}

	writeJSON(w, http.StatusOK, playlistModels)
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
		_, err := h.Queries.CreatePlaylistEntry(r.Context(), db.CreatePlaylistEntryParams{
			PlaylistID: playlist.ID,
			VideoID:    entry.ID,
			Title:      pgtype.Text{String: entry.Title, Valid: true},
			Artist:     pgtype.Text{String: entry.Uploader, Valid: entry.Uploader != ""},
			Duration:   pgtype.Int4{Int32: int32(math.Min(math.Max(entry.Duration, 0), math.MaxInt32)), Valid: true},
			Thumbnail:  pgtype.Text{String: entry.Thumbnail, Valid: entry.Thumbnail != ""},
			Status:     db.NullEntryStatus{EntryStatus: db.EntryStatusNew, Valid: true},
			FilePath:   pgtype.Text{Valid: false},
			Format:     db.NullMediaFormat{MediaFormat: db.MediaFormat(req.TargetFormat), Valid: req.TargetFormat != ""},
		})

		if err != nil {
			slog.Error("Failed to create playlist entry", "error", err)
			continue
		}
		slog.Info("Playlist entry created", "title", entry.Title, "url", entry.URL)
	}

	entries, err := h.Queries.FetchPlaylistEntries(r.Context(), playlist.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch playlist entries")
		return
	}

	writeJSON(w, http.StatusCreated, models.Playlist{
		ID:         playlist.ID,
		URL:        playlist.Url,
		Title:      playlist.Title.String,
		Entries:    toPlaylistEntryModels(entries),
		EntryCount: int64(len(entries)),
		CreatedAt:  playlist.CreatedAt.Time,
		UpdatedAt:  playlist.UpdatedAt.Time,
	})
}

func toPlaylistEntryModels(entries []db.PlaylistEntry) []models.PlaylistEntry {
	modelsEntries := make([]models.PlaylistEntry, 0, len(entries))
	for _, e := range entries {
		modelsEntries = append(modelsEntries, models.PlaylistEntry{
			ID:         e.ID,
			PlaylistID: e.PlaylistID,
			VideoID:    e.VideoID,
			Title:      e.Title.String,
			Artist:     e.Artist.String,
			Duration:   int64(e.Duration.Int32),
			Thumbnail:  e.Thumbnail.String,
			Status:     models.DownloadStatus(e.Status.EntryStatus),
			FilePath:   e.FilePath.String,
			Format:     models.AudioFormat(e.Format.MediaFormat),
			CreatedAt:  e.CreatedAt.Time,
			UpdatedAt:  e.UpdatedAt.Time,
		})
	}
	return modelsEntries
}
