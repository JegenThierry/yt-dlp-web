package models

type Download struct {
	ID          int64 `json:"id"`
	RefPlaylist int64 `json:"ref_playlist"`
}
