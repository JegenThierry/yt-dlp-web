-- name: FetchPlaylistEntries :many
SELECT * FROM playlist_entries WHERE playlist_id = $1 ORDER BY id;

-- name: FetchPlaylistEntryByID :one
SELECT * FROM playlist_entries WHERE id = $1;

-- name: UpdatePlaylistEntryStatus :one
UPDATE playlist_entries
SET status = $2, file_path = $3, error_message = $4, error_code = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePlaylistEntryMetadata :one
UPDATE playlist_entries
SET title = $3, artist = $4, genre = $5, year = $6, updated_at = NOW()
WHERE id = $1 AND playlist_id = $2
RETURNING *;

-- name: CreatePlaylistEntry :execrows
INSERT INTO playlist_entries (playlist_id, video_id, title, artist, genre, year, duration, thumbnail, status, file_path, format, source_title, source_artist, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $3, $4, NOW(), NOW())
ON CONFLICT (playlist_id, video_id) DO NOTHING;

-- name: FetchPlaylistFormat :one
SELECT format FROM playlist_entries
WHERE playlist_id = $1 AND format IS NOT NULL
ORDER BY id
LIMIT 1;

-- name: SetPlaylistEntriesStatus :exec
UPDATE playlist_entries
SET status = @status, file_path = NULL, error_message = NULL, error_code = NULL, updated_at = NOW()
WHERE id = ANY(@ids::bigint[]);

-- name: ResetInterruptedEntries :execrows
UPDATE playlist_entries
SET status = 'new', file_path = NULL, error_message = NULL, error_code = NULL, updated_at = NOW()
WHERE status IN ('downloading', 'pending');
