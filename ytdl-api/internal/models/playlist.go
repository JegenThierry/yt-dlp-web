package models

type Playlist struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Entries []Single `json:"entries"`
}
