-- name: FetchAllPlaylists :many
SELECT p.*,
       COUNT(e.id) AS entry_count,
       COUNT(e.id) FILTER (WHERE e.status = 'completed') AS completed_count,
       COUNT(e.id) FILTER (WHERE e.status = 'failed' AND e.error_code IS DISTINCT FROM 'unavailable') AS failed_count,
       COUNT(e.id) FILTER (WHERE e.status = 'failed' AND e.error_code = 'unavailable') AS unavailable_count,
       COUNT(e.id) FILTER (WHERE e.status IN ('pending', 'downloading')) AS active_count,
       COALESCE(SUM(e.duration), 0)::bigint AS total_duration,
       GREATEST(p.updated_at, MAX(e.updated_at))::timestamptz AS last_activity_at
FROM playlists p
LEFT JOIN playlist_entries e ON e.playlist_id = p.id
GROUP BY p.id
ORDER BY p.created_at DESC;

-- name: FetchPlaylistByID :one
SELECT * FROM playlists WHERE id = $1;

-- name: FetchPlaylistBySourceID :one
SELECT * FROM playlists WHERE source_id = $1;

-- name: CreatePlaylist :one
INSERT INTO playlists (url, title, quality, source_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (source_id) DO NOTHING
RETURNING *;

-- name: SetPlaylistSourceID :exec
UPDATE playlists
SET source_id = @source_id
WHERE playlists.id = @id
  AND playlists.source_id IS NULL
  AND NOT EXISTS (SELECT 1 FROM playlists other WHERE other.source_id = @source_id);

-- name: TouchPlaylist :one
UPDATE playlists SET updated_at = NOW() WHERE id = $1
RETURNING *;

-- name: DeletePlaylist :execrows
DELETE FROM playlists WHERE id = $1;
