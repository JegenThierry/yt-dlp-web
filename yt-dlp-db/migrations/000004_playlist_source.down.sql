UPDATE playlist_entries
SET thumbnail = NULL
WHERE thumbnail = 'https://i.ytimg.com/vi/' || video_id || '/mqdefault.jpg';

DROP INDEX IF EXISTS playlists_source_id_key;

ALTER TABLE playlists
    DROP COLUMN IF EXISTS source_id;
