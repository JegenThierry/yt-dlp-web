package server

import (
	"ytdl-api/internal/handlers"
	"ytdl-api/internal/repository/db"
)

func (Server *Server) registerRoutes() {
	queries := db.New(Server.db)

	// Handlers
	playlistHandler := handlers.PlaylistHandler{Queries: queries}

	// Server.router...
	Server.router.HandleFunc("POST /api/playlists", playlistHandler.CreatePlaylist)
	Server.router.HandleFunc("GET /api/playlists", playlistHandler.FetchPlaylists)
}
