-- name: FetchPlaylistEntries :many
SELECT * FROM playlist_entries WHERE playlist_id = $1 ORDER BY id;

-- name: FetchPlaylistEntryByID :one
SELECT * FROM playlist_entries WHERE id = $1;

-- name: UpdatePlaylistEntryStatus :one
UPDATE playlist_entries
SET status = $2, file_path = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePlaylistEntryMetadata :one
UPDATE playlist_entries
SET title = $2, artist = $3, genre = $4, year = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreatePlaylistEntry :one
INSERT INTO playlist_entries (playlist_id, video_id, title, artist, genre, year, duration, thumbnail, status, file_path, format, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
RETURNING *;