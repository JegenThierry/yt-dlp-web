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

type DownloadStatus string

const (
	StatusPending     DownloadStatus = "pending"
	StatusDownloading DownloadStatus = "downloading"
	StatusCompleted   DownloadStatus = "completed"
	StatusFailed      DownloadStatus = "failed"
)
