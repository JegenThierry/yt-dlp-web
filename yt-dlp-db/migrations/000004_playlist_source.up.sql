ALTER TABLE playlists
    ADD COLUMN source_id TEXT NULL;

UPDATE playlists p
SET source_id = s.source_id
FROM (
    SELECT DISTINCT ON (c.source_id) c.id, c.source_id
    FROM (
        SELECT id,
               created_at,
               CASE
                   WHEN url ~ '[?&]list=[A-Za-z0-9_-]+'
                       THEN 'https://www.youtube.com/playlist?list=' || substring(url FROM '[?&]list=([A-Za-z0-9_-]+)')
                   WHEN url ~ '[?&]v=[A-Za-z0-9_-]{11}'
                       THEN 'https://www.youtube.com/watch?v=' || substring(url FROM '[?&]v=([A-Za-z0-9_-]{11})')
                   WHEN url ~* '^https?://youtu\.be/[A-Za-z0-9_-]{11}'
                       THEN 'https://www.youtube.com/watch?v=' || substring(url FROM '(?i)^https?://youtu\.be/([A-Za-z0-9_-]{11})')
               END AS source_id
        FROM playlists
        WHERE url ~* '^https?://([a-z0-9-]+\.)*(youtube\.com|youtu\.be)/'
    ) c
    WHERE c.source_id IS NOT NULL
    ORDER BY c.source_id, c.created_at, c.id
) s
WHERE p.id = s.id;

CREATE UNIQUE INDEX playlists_source_id_key ON playlists (source_id);

UPDATE playlist_entries e
SET thumbnail = 'https://i.ytimg.com/vi/' || e.video_id || '/mqdefault.jpg'
FROM playlists p
WHERE p.id = e.playlist_id
  AND (e.thumbnail IS NULL OR e.thumbnail = '')
  AND p.url ~* '^https?://([a-z0-9-]+\.)*(youtube\.com|youtu\.be)/';
