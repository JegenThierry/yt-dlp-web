-- name: FetchAllPlaylists :many
SELECT * FROM playlists ORDER BY created_at DESC;

-- name: CreatePlaylist :one
INSERT INTO playlists (url, title, created_at, updated_at) 
VALUES ($1, $2, NOW(), NOW()) 
RETURNING *;

-- name: FetchPlaylistByID :one
SELECT * FROM playlists 
INNER JOIN playlist_entries ON playlists.id = playlist_entries.playlist_id 
WHERE playlists.id = $1;

-- name: CreatePlaylistEntry :one
INSERT INTO playlist_entries (playlist_id, video_id, title, artist, duration, thumbnail, status, file_path, format, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING *;