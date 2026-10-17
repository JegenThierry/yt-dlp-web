package ytdlp

import (
	"fmt"
	"regexp"
	"strings"
)

type ErrorCode string

const (
	ErrorUnavailable    ErrorCode = "unavailable"
	ErrorSignInRequired ErrorCode = "sign_in_required"
	ErrorBlocked        ErrorCode = "blocked"
	ErrorNetwork        ErrorCode = "network"
	ErrorPostprocessing ErrorCode = "postprocessing"
	ErrorUnknown        ErrorCode = "unknown"
)

const maxErrorMessageLength = 500

type DownloadError struct {
	Code    ErrorCode
	Message string
	err     error
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("yt-dlp download failed: %v: %s", e.err, e.Message)
}

func (e *DownloadError) Unwrap() error {
	return e.err
}

func newDownloadError(err error, errorLine string) *DownloadError {
	message := cleanErrorLine(errorLine)
	if message == "" {
		message = fmt.Sprintf("yt-dlp stopped without an error message (%v)", err)
	}
	return &DownloadError{Code: classifyError(message), Message: message, err: err}
}

var errorLinePrefix = regexp.MustCompile(`^ERROR:\s*(?:\[[^\]]+\]\s*(?:[^\s:]+:\s+)?)?`)

func cleanErrorLine(line string) string {
	message := errorLinePrefix.ReplaceAllString(strings.TrimSpace(line), "")
	if i := strings.Index(message, "; please report this issue"); i >= 0 {
		message = message[:i]
	}
	message = strings.Join(strings.Fields(message), " ")

	if runes := []rune(message); len(runes) > maxErrorMessageLength {
		message = string(runes[:maxErrorMessageLength-1]) + "…"
	}
	return message
}

var errorRules = []struct {
	code     ErrorCode
	patterns []string
}{
	{ErrorBlocked, []string{"http error 429", "too many requests", "http error 403", "rate-limit", "rate limit", "try again later"}},
	{ErrorUnavailable, []string{"private video", "members-only", "join this channel", "copyright", "has been terminated", "has been removed", "no longer available", "not available in your country", "not made this video available"}},
	{ErrorSignInRequired, []string{"sign in to confirm", "age-restricted", "age restricted", "inappropriate for some users", "--cookies"}},
	{ErrorNetwork, []string{"unable to download webpage", "unable to download api page", "timed out", "connection reset", "connection refused", "connection aborted", "remote end closed", "name resolution", "getaddrinfo", "network is unreachable", "urlopen error"}},
	{ErrorPostprocessing, []string{"postprocessing", "ffmpeg", "ffprobe"}},
	{ErrorUnavailable, []string{"video unavailable", "this video is not available", "does not exist"}},
}

func classifyError(message string) ErrorCode {
	lower := strings.ToLower(message)
	for _, rule := range errorRules {
		if containsAny(lower, rule.patterns) {
			return rule.code
		}
	}
	return ErrorUnknown
}

type MetadataErrorKind string

const (
	MetadataUnsupported    MetadataErrorKind = "unsupported"
	MetadataNotFound       MetadataErrorKind = "not_found"
	MetadataSignInRequired MetadataErrorKind = "sign_in_required"
	MetadataBlocked        MetadataErrorKind = "blocked"
	MetadataFailed         MetadataErrorKind = "failed"
)

type MetadataError struct {
	Kind    MetadataErrorKind
	Message string
	err     error
}

func (e *MetadataError) Error() string {
	return fmt.Sprintf("yt-dlp metadata failed: %v: %s", e.err, e.Message)
}

func (e *MetadataError) Unwrap() error {
	return e.err
}

func newMetadataError(err error, errorLine string) *MetadataError {
	message := cleanErrorLine(errorLine)
	if message == "" {
		message = fmt.Sprintf("yt-dlp stopped without an error message (%v)", err)
	}
	return &MetadataError{Kind: classifyMetadataError(message), Message: message, err: err}
}

var (
	unsupportedPatterns = []string{"is not a valid url", "unsupported url"}
	notFoundPatterns    = []string{"http error 404", "http error 400", "does not exist", "private"}
)

func classifyMetadataError(message string) MetadataErrorKind {
	lower := strings.ToLower(message)
	if containsAny(lower, unsupportedPatterns) {
		return MetadataUnsupported
	}

	code := classifyError(message)
	switch {
	case code == ErrorBlocked:
		return MetadataBlocked
	case code == ErrorUnavailable || containsAny(lower, notFoundPatterns):
		return MetadataNotFound
	case code == ErrorSignInRequired:
		return MetadataSignInRequired
	default:
		return MetadataFailed
	}
}

func containsAny(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	return false
}
