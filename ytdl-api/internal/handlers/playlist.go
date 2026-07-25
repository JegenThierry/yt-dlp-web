package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"ytdl-api/internal/models"
)

type PlaylistHandler struct {
	DB *pgxpool.Pool
}

func (h *PlaylistHandler) FetchPlaylists(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, url, title, target_format FROM playlists ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list playlists")
		return
	}
	defer rows.Close()

	playlists := []models.Playlist{}
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.URL, &p.Title, &p.TargetFormat); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan playlist")
			return
		}
		playlists = append(playlists, p)
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

	log.Printf("Downloading playlist from URL: %s with target format: %s", req.URL, req.TargetFormat)
}
