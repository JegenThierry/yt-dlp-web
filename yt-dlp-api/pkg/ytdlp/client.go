package ytdlp

import (
	"encoding/json"
	"os/exec"
)

type YtDlpEntry struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Uploader   string `json:"uploader"`
	Thumbnail  string `json:"thumbnail"`
	Thumbnails []struct {
		URL string `json:"url"`
	} `json:"thumbnails"`
	Duration float64 `json:"duration"`
	URL      string  `json:"url"`
}

type YtDlpMetadata struct {
	Title   string       `json:"title"`
	Entries []YtDlpEntry `json:"entries"`
}

func (e YtDlpEntry) GetThumbnail() string {
	if e.Thumbnail != "" {
		return e.Thumbnail
	}
	if len(e.Thumbnails) > 0 {
		return e.Thumbnails[len(e.Thumbnails)-1].URL
	}
	return ""
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
