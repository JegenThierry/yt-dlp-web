package ytdlp

import (
	"encoding/json"
	"os/exec"
)

type YtDlpMetadata struct {
	Title   string `json:"title"`
	Entries []struct {
		ID        string  `json:"id"`
		Title     string  `json:"title"`
		Artist    string  `json:"artist"`
		Uploader  string  `json:"uploader"`
		Thumbnail string  `json:"thumbnail"`
		Duration  float64 `json:"duration"`
		URL       string  `json:"url"`
	} `json:"entries"`
}

func GetPlaylist(url string) (*YtDlpMetadata, error) {
	out, err := exec.Command(string(YT_DLP), string(FLAT_PLAYLIST), string(DUMP_SINGLE_JSON), url).Output()

	if err != nil {
		return nil, err
	}

	var metadata YtDlpMetadata
	err = json.Unmarshal(out, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
