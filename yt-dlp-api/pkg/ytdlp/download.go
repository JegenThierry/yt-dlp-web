package ytdlp

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Cookies struct {
	File        string
	FromBrowser string
}

func (c Cookies) args() []string {
	if c.File != "" {
		return []string{string(COOKIES), c.File}
	}
	if c.FromBrowser != "" {
		return []string{string(COOKIES_FROM_BROWSER), c.FromBrowser}
	}
	return nil
}

type FilterOptions struct {
	Format  AudioCodecs
	Quality AudioQuality
	Cookies Cookies
}

func (o FilterOptions) withDefaults() FilterOptions {
	if o.Format == "" {
		o.Format = MP3
	}
	if o.Quality == "" {
		o.Quality = QualityBest
	}
	return o
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

func ContentHash(videoID string, opts FilterOptions) string {
	opts = opts.withDefaults()

	h := sha256.New()
	fmt.Fprintf(h, "%s|%s|%s|true", videoID, opts.Format, opts.Quality)
	return hex.EncodeToString(h.Sum(nil))
}

func expectedExtension(format AudioCodecs) string {
	return extension(format)
}

type Phase string

const (
	PhaseDownloading Phase = "downloading"
	PhaseConverting  Phase = "converting"
)

type Progress struct {
	Phase   Phase
	Percent *float64
	ETA     *int64
}

const (
	downloadLinePrefix    = "ytdl-api:download"
	postprocessLinePrefix = "ytdl-api:postprocess"
	downloadTemplate      = "download:" + downloadLinePrefix + " %(progress.downloaded_bytes)s %(progress.total_bytes)s %(progress.total_bytes_estimate)s %(progress.eta)s"
	postprocessTemplate   = "postprocess:" + postprocessLinePrefix + " %(progress.status)s %(progress.postprocessor)s"
)

func Download(ctx context.Context, input, outputDir, tempDir, hash string, opts FilterOptions, onProgress func(Progress)) (string, error) {
	opts = opts.withDefaults()

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	args := []string{
		string(ENCODING), "utf-8",
		string(FORMAT), "bestaudio/best",
		string(EXTRACT_AUDIO),
		string(AUDIO_FORMAT), string(opts.Format),
		string(AUDIO_QUALITY), toAudioQualityValue(opts.Quality),
		string(EMBED_METADATA),
	}

	args = append(args, opts.Cookies.args()...)

	filename := fmt.Sprintf("%s.%s", hash, expectedExtension(opts.Format))
	args = append(args,
		string(PATHS), "home:"+outputDir,
		string(PATHS), "temp:"+tempDir,
		string(OUTPUT), filename,
		string(NEWLINE),
		string(PROGRESS),
		string(PROGRESS_TEMPLATE), downloadTemplate,
		string(PROGRESS_TEMPLATE), postprocessTemplate,
		string(PRINT), "after_move:filepath",
		string(END_OF_OPTIONS), input,
	)

	cmd := exec.CommandContext(ctx, string(YT_DLP), args...)
	cmd.Cancel = func() error {
		return interrupt(cmd)
	}
	cmd.WaitDelay = 10 * time.Second

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("yt-dlp stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("yt-dlp stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start yt-dlp: %w", err)
	}

	var reportMu sync.Mutex
	report := func(p Progress) {
		if onProgress == nil {
			return
		}
		reportMu.Lock()
		defer reportMu.Unlock()
		onProgress(p)
	}

	var filePath, lastError string
	var downloadStarted atomic.Bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		scanLines(stdout, func(line string) {
			if rest, ok := strings.CutPrefix(line, downloadLinePrefix); ok {
				downloadStarted.Store(true)
				report(parseDownloadProgress(rest))
				return
			}
			if line != "" {
				filePath = line
			}
		})
	}()
	go func() {
		defer wg.Done()
		scanLines(stderr, func(line string) {
			if rest, ok := strings.CutPrefix(line, postprocessLinePrefix); ok {
				if downloadStarted.Load() && strings.HasPrefix(strings.TrimSpace(rest), "started") {
					report(Progress{Phase: PhaseConverting})
				}
				return
			}
			if strings.HasPrefix(line, "ERROR:") {
				lastError = line
			}
		})
	}()
	wg.Wait()
	waitErr := cmd.Wait()

	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if waitErr != nil {
		return "", newDownloadError(waitErr, lastError)
	}

	if filePath == "" {
		filePath = filepath.Join(outputDir, filename)
	}

	return filePath, nil
}

func scanLines(r io.Reader, handle func(line string)) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		handle(strings.TrimSpace(scanner.Text()))
	}
	io.Copy(io.Discard, r)
}

func parseDownloadProgress(fields string) Progress {
	progress := Progress{Phase: PhaseDownloading}

	values := strings.Fields(fields)
	if len(values) != 4 {
		return progress
	}

	downloaded, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return progress
	}

	total, err := strconv.ParseFloat(values[1], 64)
	if err != nil || total <= 0 {
		total, err = strconv.ParseFloat(values[2], 64)
	}
	if err == nil && total > 0 {
		percent := math.Round(math.Min(downloaded/total, 1)*1000) / 10
		progress.Percent = &percent
	}

	if eta, err := strconv.ParseFloat(values[3], 64); err == nil && eta >= 0 {
		seconds := int64(math.Round(eta))
		progress.ETA = &seconds
	}

	return progress
}
