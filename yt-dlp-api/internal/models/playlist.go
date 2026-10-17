package models

import (
	"time"
)

type Playlist struct {
	ID               int64           `json:"id"`
	URL              string          `json:"url"`
	Title            string          `json:"title"`
	Quality          AudioQuality    `json:"quality"`
	Entries          []PlaylistEntry `json:"entries"`
	EntryCount       int64           `json:"entry_count"`
	CompletedCount   int64           `json:"completed_count"`
	FailedCount      int64           `json:"failed_count"`
	UnavailableCount int64           `json:"unavailable_count"`
	ActiveCount      int64           `json:"active_count"`
	TotalDuration    int64           `json:"total_duration"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	LastActivityAt   time.Time       `json:"last_activity_at"`
}

type PlaylistImport struct {
	Playlist
	Existing     bool  `json:"existing"`
	AddedCount   int64 `json:"added_count"`
	SkippedCount int64 `json:"skipped_count"`
}

type CreatePlaylistRequest struct {
	URL          string       `json:"url"`
	TargetFormat AudioFormat  `json:"target_format"`
	Quality      AudioQuality `json:"quality"`
}

type DownloadPlaylistRequest struct {
	Filters          *FilterOptions `json:"filters,omitempty"`
	RetryUnavailable bool           `json:"retry_unavailable"`
}

func (r *DownloadPlaylistRequest) FilterOptions() *FilterOptions { return r.Filters }
