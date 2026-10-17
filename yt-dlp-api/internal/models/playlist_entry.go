package models

import "time"

type PlaylistEntry struct {
	ID           int64          `json:"id"`
	PlaylistID   int64          `json:"playlist_id"`
	VideoID      string         `json:"video_id"`
	Title        string         `json:"title"`
	Artist       string         `json:"artist"`
	SourceTitle  string         `json:"source_title"`
	SourceArtist string         `json:"source_artist"`
	Genre        string         `json:"genre"`
	Year         int64          `json:"year"`
	Duration     int64          `json:"duration"`
	Thumbnail    string         `json:"thumbnail"`
	Status       DownloadStatus `json:"status"`
	FilePath     string         `json:"file_path"`
	Format       AudioFormat    `json:"format"`
	ErrorMessage string         `json:"error_message,omitempty"`
	ErrorCode    EntryErrorCode `json:"error_code,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type PlaylistEntryMetadataUpdate struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Genre  string `json:"genre"`
	Year   *int64 `json:"year"`
}

type MetadataFieldError struct {
	EntryID int64  `json:"entry_id"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type MetadataValidationError struct {
	Error  string               `json:"error"`
	Fields []MetadataFieldError `json:"fields"`
}

type DownloadEntryRequest struct {
	Filters *FilterOptions `json:"filters,omitempty"`
}

func (r *DownloadEntryRequest) FilterOptions() *FilterOptions { return r.Filters }
