package models

import "time"

type PlaylistEntry struct {
	ID         int64          `json:"id"`
	PlaylistID int64          `json:"playlist_id"`
	VideoID    string         `json:"video_id"`
	Title      string         `json:"title"`
	Artist     string         `json:"artist"`
	Genre      string         `json:"genre"`
	Year       int64          `json:"year"`
	Duration   int64          `json:"duration"`
	Thumbnail  string         `json:"thumbnail"`
	Status     DownloadStatus `json:"status"`
	FilePath   string         `json:"file_path"`
	Format     AudioFormat    `json:"format"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type PlaylistEntryMetadataUpdate struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Genre  string `json:"genre"`
	Year   int64  `json:"year"`
}

type DownloadEntryRequest struct {
	Filters *FilterOptions `json:"filters,omitempty"`
}

func (r *DownloadEntryRequest) FilterOptions() *FilterOptions { return r.Filters }
