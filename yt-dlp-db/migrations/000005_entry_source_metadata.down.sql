ALTER TABLE playlist_entries
    DROP COLUMN IF EXISTS source_artist,
    DROP COLUMN IF EXISTS source_title;
