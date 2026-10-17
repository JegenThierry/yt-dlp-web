ALTER TABLE playlist_entries
    ADD COLUMN source_title TEXT NULL,
    ADD COLUMN source_artist TEXT NULL;

UPDATE playlist_entries
SET source_title = title,
    source_artist = artist;
