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

type AudioQuality string

const (
	QualityBest     AudioQuality = "best"
	QualityHigh     AudioQuality = "high"
	QualityStandard AudioQuality = "standard"
	QualityLow      AudioQuality = "low"
)

type DownloadStatus string

const (
	StatusNew         DownloadStatus = "new"
	StatusPending     DownloadStatus = "pending"
	StatusDownloading DownloadStatus = "downloading"
	StatusCompleted   DownloadStatus = "completed"
	StatusFailed      DownloadStatus = "failed"
)
