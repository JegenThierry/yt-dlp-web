-- name: FetchAllPlaylists :many
SELECT p.*, COUNT(e.id) AS entry_count
FROM playlists p
LEFT JOIN playlist_entries e ON e.playlist_id = p.id
GROUP BY p.id
ORDER BY p.created_at DESC;

-- name: FetchPlaylistByID :one
SELECT * FROM playlists WHERE id = $1;

-- name: CreatePlaylist :one
INSERT INTO playlists (url, title, created_at, updated_at) 
VALUES ($1, $2, NOW(), NOW()) 
RETURNING *;