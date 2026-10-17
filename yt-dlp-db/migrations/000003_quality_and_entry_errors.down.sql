ALTER TABLE playlist_entries
    DROP COLUMN IF EXISTS error_code,
    DROP COLUMN IF EXISTS error_message;

ALTER TABLE playlists
    DROP COLUMN IF EXISTS quality;

DROP TYPE IF EXISTS entry_error_code;

DROP TYPE IF EXISTS audio_quality;
