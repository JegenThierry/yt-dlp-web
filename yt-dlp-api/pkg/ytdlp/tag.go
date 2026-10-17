package ytdlp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Metadata struct {
	Title  string
	Artist string
	Genre  string
	Year   int64
}

func (m Metadata) Empty() bool {
	return m.Title == "" && m.Artist == "" && m.Genre == "" && m.Year == 0
}

func (m Metadata) Hash() string {
	h := sha256.New()
	fmt.Fprintf(h, "%q|%q|%q|%d", m.Title, m.Artist, m.Genre, m.Year)
	return hex.EncodeToString(h.Sum(nil))
}

func (m Metadata) tags() []string {
	var tags []string
	if m.Title != "" {
		tags = append(tags, "title="+m.Title)
	}
	if m.Artist != "" {
		tags = append(tags, "artist="+m.Artist)
	}
	if m.Genre != "" {
		tags = append(tags, "genre="+m.Genre)
	}
	if m.Year != 0 {
		tags = append(tags, "date="+strconv.FormatInt(m.Year, 10))
	}
	return tags
}

func Tag(ctx context.Context, source, target string, meta Metadata) error {
	args := []string{
		string(FFMPEG_NO_STDIN),
		string(FFMPEG_HIDE_BANNER),
		string(FFMPEG_LOG_LEVEL), "error",
		string(FFMPEG_OVERWRITE),
		string(FFMPEG_INPUT), source,
		string(FFMPEG_MAP), "0",
		string(FFMPEG_CODEC), "copy",
	}
	for _, tag := range meta.tags() {
		args = append(args, string(FFMPEG_METADATA), tag, string(FFMPEG_STREAM_METADATA), tag)
	}
	args = append(args, target)

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, string(FFMPEG), args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg tagging failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
