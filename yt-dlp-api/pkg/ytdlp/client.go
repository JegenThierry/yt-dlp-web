package ytdlp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const maxThumbnailWidth = 480

type Playlist struct {
	Title     string
	SourceURL string
	Videos    []Video
	Skipped   int
}

type Video struct {
	ID        string
	Title     string
	Artist    string
	Duration  int32
	Thumbnail string
}

type thumbnail struct {
	URL   string `json:"url"`
	Width int    `json:"width"`
}

type rawEntry struct {
	Type         string      `json:"_type"`
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Uploader     string      `json:"uploader"`
	Channel      string      `json:"channel"`
	Duration     float64     `json:"duration"`
	Thumbnail    string      `json:"thumbnail"`
	Thumbnails   []thumbnail `json:"thumbnails"`
	IEKey        string      `json:"ie_key"`
	ExtractorKey string      `json:"extractor_key"`
	WebpageURL   string      `json:"webpage_url"`
	Entries      []rawEntry  `json:"entries"`
}

var placeholderTitles = map[string]bool{
	"[Private video]":     true,
	"[Deleted video]":     true,
	"[Unavailable video]": true,
}

var nestedListKeys = map[string]bool{
	"YoutubeTab":      true,
	"YoutubePlaylist": true,
}

var hiddenVideosPattern = regexp.MustCompile(`(?i)(\d+) unavailable videos? (?:is|are) hidden`)

func GetPlaylist(ctx context.Context, input string, cookies Cookies) (*Playlist, error) {
	args := []string{
		string(ENCODING), "utf-8",
		string(FLAT_PLAYLIST),
		string(DUMP_SINGLE_JSON),
	}
	args = append(args, cookies.args()...)
	args = append(args, string(END_OF_OPTIONS), input)

	cmd := exec.CommandContext(ctx, string(YT_DLP), args...)
	cmd.Cancel = func() error {
		return interrupt(cmd)
	}
	cmd.WaitDelay = 10 * time.Second

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	errorLine, hidden := scanMetadataStderr(stderr.String())
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil, newMetadataError(err, errorLine)
	}
	if err != nil {
		return nil, fmt.Errorf("run yt-dlp: %w", err)
	}

	var raw rawEntry
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("parse yt-dlp output: %w", err)
	}

	playlist := raw.toPlaylist()
	playlist.Skipped += hidden
	return playlist, nil
}

func scanMetadataStderr(stderr string) (errorLine string, hidden int) {
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ERROR:") {
			errorLine = line
			continue
		}
		if match := hiddenVideosPattern.FindStringSubmatch(line); match != nil {
			if n, err := strconv.Atoi(match[1]); err == nil {
				hidden = max(hidden, n)
			}
		}
	}
	return errorLine, hidden
}

func (r rawEntry) toPlaylist() *Playlist {
	playlist := &Playlist{
		Title:     strings.TrimSpace(r.Title),
		SourceURL: r.WebpageURL,
	}
	if playlist.Title == "" {
		playlist.Title = r.ID
	}

	entries := r.Entries
	if r.Type == "video" || (r.Type == "" && r.Entries == nil && r.ID != "") {
		entries = []rawEntry{r}
	}

	playlist.collect(entries, make(map[string]bool))
	return playlist
}

func (p *Playlist) collect(entries []rawEntry, seen map[string]bool) {
	for _, entry := range entries {
		if entry.Type == "playlist" || len(entry.Entries) > 0 {
			p.collect(entry.Entries, seen)
			continue
		}
		if entry.ID == "" || nestedListKeys[entry.IEKey] {
			continue
		}
		if placeholderTitles[strings.TrimSpace(entry.Title)] {
			p.Skipped++
			continue
		}
		if seen[entry.ID] {
			continue
		}
		seen[entry.ID] = true
		p.Videos = append(p.Videos, entry.toVideo())
	}
}

func (e rawEntry) toVideo() Video {
	artist := strings.TrimSpace(e.Uploader)
	if artist == "" {
		artist = strings.TrimSpace(e.Channel)
	}
	artist = strings.TrimSpace(strings.TrimSuffix(artist, " - Topic"))

	return Video{
		ID:        e.ID,
		Title:     strings.TrimSpace(e.Title),
		Artist:    artist,
		Duration:  int32(math.Min(math.Max(e.Duration, 0), math.MaxInt32)),
		Thumbnail: e.pickThumbnail(),
	}
}

func (e rawEntry) pickThumbnail() string {
	var best thumbnail
	for _, t := range e.Thumbnails {
		if t.URL != "" && t.Width > best.Width && t.Width <= maxThumbnailWidth {
			best = t
		}
	}
	if best.URL != "" {
		return best.URL
	}
	if e.IEKey == "Youtube" || e.ExtractorKey == "Youtube" {
		return "https://i.ytimg.com/vi/" + url.PathEscape(e.ID) + "/mqdefault.jpg"
	}
	return e.Thumbnail
}
