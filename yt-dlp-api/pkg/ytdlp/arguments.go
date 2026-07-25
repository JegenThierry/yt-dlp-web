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

type CommandArguments string

const (
	FORMAT           CommandArguments = "--format"
	AUDIO_FORMAT     CommandArguments = "--audio-format"
	EMBED_THUMBNAIL  CommandArguments = "--embed-thumbnail"
	EMBED_METADATA   CommandArguments = "--embed-metadata"
	OUTPUT           CommandArguments = "--output"
	FLAT_PLAYLIST    CommandArguments = "--flat-playlist"
	DUMP_JSON        CommandArguments = "--dump-json"
	DUMP_SINGLE_JSON CommandArguments = "--dump-single-json"
)
