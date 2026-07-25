package models

type PlaylistEntry struct {
	ID         int64          `json:"id"`
	PlaylistID int64          `json:"playlist_id"`
	VideoID    string         `json:"video_id"`
	Title      string         `json:"title"`
	Artist     string         `json:"artist"`
	Duration   int64          `json:"duration"`
	Thumbnail  string         `json:"thumbnail"`
	Status     DownloadStatus `json:"status"`
	FilePath   string         `json:"file_path"`
	Format     AudioFormat    `json:"format"`
}
