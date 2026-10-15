package ytdlp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FilterOptions struct {
	Format             AudioCodecs
	Quality            AudioQuality
	EmbedMetadata      bool
	CookiesFromBrowser string
	CookiesFile        string
	Metadata           Metadata
}

type Metadata struct {
	Title  string
	Artist string
	Genre  string
	Year   int64
}

func (m Metadata) empty() bool {
	return m.Title == "" && m.Artist == "" && m.Genre == "" && m.Year == 0
}

func (m Metadata) metadataArgs() []string {
	var args []string
	if m.Title != "" {
		args = append(args, string(METADATA), "title="+m.Title)
	}
	if m.Artist != "" {
		args = append(args, string(METADATA), "artist="+m.Artist)
	}
	if m.Genre != "" {
		args = append(args, string(METADATA), "genre="+m.Genre)
	}
	if m.Year != 0 {
		args = append(args, string(METADATA), fmt.Sprintf("release_year=%d", m.Year))
	}
	return args
}

func extension(format AudioCodecs) string {
	switch format {
	case MP3:
		return "mp3"
	case M4A, AAC, ALAC:
		return "m4a"
	case OPUS:
		return "opus"
	case VORBIS:
		return "ogg"
	case FLAC:
		return "flac"
	case WAV:
		return "wav"
	default:
		return string(format)
	}
}

// ContentHash returns the stable filename base (SHA-256 hex) for a combination
// of source video and conversion settings. Two downloads that share the same
// video, format, quality and metadata-embedding flag produce the same hash and
// therefore reuse the same cached file on disk. Editable metadata fields
// (title/artist/genre/year) are intentionally excluded so metadata edits do not
// force a re-download.
func ContentHash(videoID string, opts FilterOptions) string {
	format := opts.Format
	if format == "" {
		format = MP3
	}

	quality := opts.Quality
	if quality == "" {
		quality = QualityBest
	}

	h := sha256.New()
	fmt.Fprintf(h, "%s|%s|%s|%t", videoID, format, quality, opts.EmbedMetadata)
	return hex.EncodeToString(h.Sum(nil))
}

// expectedExtension maps a resolved audio format to the file extension yt-dlp
// produces after extraction. It mirrors the extension() switch.
func expectedExtension(format AudioCodecs) string {
	return extension(format)
}

func Download(input, outputDir, hash string, opts FilterOptions) (string, error) {
	format := opts.Format
	if format == "" {
		format = MP3
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	args := []string{
		string(FORMAT), "bestaudio/best",
		string(EXTRACT_AUDIO),
		string(AUDIO_FORMAT), string(format),
	}

	if opts.Quality != "" {
		args = append(args, "--audio-quality", toAudioQualityValue(opts.Quality))
	}

	if opts.EmbedMetadata {
		args = append(args, string(EMBED_METADATA))
	}

	if !opts.Metadata.empty() {
		args = append(args, string(EMBED_METADATA))
		args = append(args, opts.Metadata.metadataArgs()...)
	}

	if opts.CookiesFile != "" {
		args = append(args, string(COOKIES), opts.CookiesFile)
	} else if opts.CookiesFromBrowser != "" {
		args = append(args, string(COOKIES_FROM_BROWSER), opts.CookiesFromBrowser)
	}

	outputTemplate := filepath.Join(outputDir, fmt.Sprintf("%s.%s", hash, expectedExtension(format)))
	args = append(args, string(OUTPUT), outputTemplate, "--print", "after_move:filepath", input)

	out, err := exec.Command(string(YT_DLP), args...).Output()
	if err != nil {
		return "", fmt.Errorf("yt-dlp download failed: %w", err)
	}

	filePath := strings.TrimSpace(string(out))
	if filePath == "" {
		filePath = outputTemplate
	}

	return filePath, nil
}
