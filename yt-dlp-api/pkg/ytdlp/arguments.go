package ytdlp

type AudioCodecs string

const (
	AAC    AudioCodecs = "aac"
	ALAC   AudioCodecs = "alac"
	FLAC   AudioCodecs = "flac"
	M4A    AudioCodecs = "m4a"
	MP3    AudioCodecs = "mp3"
	OPUS   AudioCodecs = "opus"
	VORBIS AudioCodecs = "vorbis"
	WAV    AudioCodecs = "wav"
)

type AudioQuality string

const (
	QualityBest     AudioQuality = "best"
	QualityHigh     AudioQuality = "high"
	QualityStandard AudioQuality = "standard"
	QualityLow      AudioQuality = "low"
)

func toAudioQualityValue(q AudioQuality) string {
	switch q {
	case QualityHigh:
		return "192K"
	case QualityStandard:
		return "128K"
	case QualityLow:
		return "64K"
	case QualityBest, "":
		return "0"
	default:
		return "0"
	}
}

type CommandArguments string

const (
	FORMAT               CommandArguments = "--format"
	AUDIO_FORMAT         CommandArguments = "--audio-format"
	EXTRACT_AUDIO        CommandArguments = "--extract-audio"
	EMBED_THUMBNAIL      CommandArguments = "--embed-thumbnail"
	EMBED_METADATA       CommandArguments = "--embed-metadata"
	METADATA             CommandArguments = "--metadata"
	OUTPUT               CommandArguments = "--output"
	COOKIES              CommandArguments = "--cookies"
	COOKIES_FROM_BROWSER CommandArguments = "--cookies-from-browser"
	FLAT_PLAYLIST        CommandArguments = "--flat-playlist"
	DUMP_JSON            CommandArguments = "--dump-json"
	DUMP_SINGLE_JSON     CommandArguments = "--dump-single-json"
)

type CommandLineTools string

const (
	YT_DLP CommandLineTools = "yt-dlp"
)
