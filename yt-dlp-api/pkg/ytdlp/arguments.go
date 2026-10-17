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
	AUDIO_QUALITY        CommandArguments = "--audio-quality"
	EXTRACT_AUDIO        CommandArguments = "--extract-audio"
	EMBED_THUMBNAIL      CommandArguments = "--embed-thumbnail"
	EMBED_METADATA       CommandArguments = "--embed-metadata"
	OUTPUT               CommandArguments = "--output"
	PATHS                CommandArguments = "--paths"
	PRINT                CommandArguments = "--print"
	NEWLINE              CommandArguments = "--newline"
	PROGRESS             CommandArguments = "--progress"
	PROGRESS_TEMPLATE    CommandArguments = "--progress-template"
	COOKIES              CommandArguments = "--cookies"
	COOKIES_FROM_BROWSER CommandArguments = "--cookies-from-browser"
	FLAT_PLAYLIST        CommandArguments = "--flat-playlist"
	DUMP_JSON            CommandArguments = "--dump-json"
	DUMP_SINGLE_JSON     CommandArguments = "--dump-single-json"
	ENCODING             CommandArguments = "--encoding"
	END_OF_OPTIONS       CommandArguments = "--"
)

const (
	FFMPEG_NO_STDIN        CommandArguments = "-nostdin"
	FFMPEG_HIDE_BANNER     CommandArguments = "-hide_banner"
	FFMPEG_LOG_LEVEL       CommandArguments = "-loglevel"
	FFMPEG_OVERWRITE       CommandArguments = "-y"
	FFMPEG_INPUT           CommandArguments = "-i"
	FFMPEG_MAP             CommandArguments = "-map"
	FFMPEG_CODEC           CommandArguments = "-c"
	FFMPEG_METADATA        CommandArguments = "-metadata"
	FFMPEG_STREAM_METADATA CommandArguments = "-metadata:s:a"
)

const (
	TASKKILL_TREE  CommandArguments = "/T"
	TASKKILL_FORCE CommandArguments = "/F"
	TASKKILL_PID   CommandArguments = "/PID"
)

type CommandLineTools string

const (
	YT_DLP   CommandLineTools = "yt-dlp"
	FFMPEG   CommandLineTools = "ffmpeg"
	TASKKILL CommandLineTools = "taskkill"
)
