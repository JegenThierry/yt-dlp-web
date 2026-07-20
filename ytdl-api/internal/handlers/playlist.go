package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"ytdl-api/models"
)

type PlaylistHandler struct {
	DB *pgxpool.Pool
}

func (h *PlaylistHandler) FetchPlaylists(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, name FROM playlists ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list playlists")
		return
	}
	defer rows.Close()

	playlists := []models.Playlist{}
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan playlist")
			return
		}
		playlists = append(playlists, p)
	}

	writeJSON(w, http.StatusOK, playlists)
}
