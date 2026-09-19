-- name: FetchPlaylistEntries :many
SELECT * FROM playlist_entries WHERE playlist_id = $1 ORDER BY id;

-- name: CreatePlaylistEntry :one
INSERT INTO playlist_entries (playlist_id, video_id, title, artist, duration, thumbnail, status, file_path, format, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING *;