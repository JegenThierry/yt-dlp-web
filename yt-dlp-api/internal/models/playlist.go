package models

type Playlist struct {
	ID           int64           `json:"id"`
	URL          string          `json:"url"`
	Title        string          `json:"title"`
	TargetFormat AudioFormat     `json:"target_format"`
	Entries      []PlaylistEntry `json:"entries"`
}

type CreatePlaylistRequest struct {
	URL          string      `json:"url"`
	TargetFormat AudioFormat `json:"target_format"`
}
