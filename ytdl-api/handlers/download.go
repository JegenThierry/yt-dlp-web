package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"ytdl-api/models"
)

type DownloadHandler struct {
	DB *pgxpool.Pool
}

func (h *DownloadHandler) FetchDownloads(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(r.Context(), "SELECT id, ref_playlist FROM downloads ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list downloads")
		return
	}
	defer rows.Close()

	downloads := []models.Download{}
	for rows.Next() {
		var d models.Download
		if err := rows.Scan(&d.ID, &d.RefPlaylist); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan download")
			return
		}
		downloads = append(downloads, d)
	}

	writeJSON(w, http.StatusOK, downloads)
}
