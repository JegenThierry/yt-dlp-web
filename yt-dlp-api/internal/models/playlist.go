package models

import (
	"time"
)

type Playlist struct {
	ID        int64           `json:"id"`
	URL       string          `json:"url"`
	Title     string          `json:"title"`
	Entries   []PlaylistEntry `json:"entries"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreatePlaylistRequest struct {
	URL          string      `json:"url"`
	TargetFormat AudioFormat `json:"target_format"`
}
