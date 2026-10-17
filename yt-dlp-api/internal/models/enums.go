package models

type AudioFormat string

const (
	MP3    AudioFormat = "mp3"
	AAC    AudioFormat = "aac"
	FLAC   AudioFormat = "flac"
	WAV    AudioFormat = "wav"
	M4A    AudioFormat = "m4a"
	OPUS   AudioFormat = "opus"
	VORBIS AudioFormat = "vorbis"
	ALAC   AudioFormat = "alac"
)

func (f AudioFormat) Valid() bool {
	switch f {
	case MP3, AAC, FLAC, WAV, M4A, OPUS, VORBIS, ALAC:
		return true
	default:
		return false
	}
}

type AudioQuality string

const (
	QualityBest     AudioQuality = "best"
	QualityHigh     AudioQuality = "high"
	QualityStandard AudioQuality = "standard"
	QualityLow      AudioQuality = "low"
)

func (q AudioQuality) Valid() bool {
	switch q {
	case QualityBest, QualityHigh, QualityStandard, QualityLow:
		return true
	default:
		return false
	}
}

type DownloadStatus string

const (
	StatusNew         DownloadStatus = "new"
	StatusPending     DownloadStatus = "pending"
	StatusDownloading DownloadStatus = "downloading"
	StatusCompleted   DownloadStatus = "completed"
	StatusFailed      DownloadStatus = "failed"
)

type DownloadPhase string

const (
	PhaseDownloading DownloadPhase = "downloading"
	PhaseConverting  DownloadPhase = "converting"
)

type EntryErrorCode string

const (
	ErrorUnavailable    EntryErrorCode = "unavailable"
	ErrorSignInRequired EntryErrorCode = "sign_in_required"
	ErrorBlocked        EntryErrorCode = "blocked"
	ErrorNetwork        EntryErrorCode = "network"
	ErrorPostprocessing EntryErrorCode = "postprocessing"
	ErrorUnknown        EntryErrorCode = "unknown"
)
