-- name: FetchAllPlaylists :many
SELECT * FROM playlists ORDER BY created_at DESC;

-- name: FetchPlaylistByID :one
SELECT * FROM playlists 
INNER JOIN playlist_entries ON playlists.id = playlist_entries.playlist_id 
WHERE playlists.id = $1;

-- name: CreatePlaylist :one
INSERT INTO playlists (url, title, created_at, updated_at) 
VALUES ($1, $2, NOW(), NOW()) 
RETURNING *;