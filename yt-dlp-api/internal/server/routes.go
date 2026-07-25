package server

import "ytdl-api/internal/handlers"

func (Server *Server) registerRoutes() {
	// Handlers
	playlistHandler := handlers.PlaylistHandler{DB: Server.db}

	// Server.router...
	Server.router.HandleFunc("POST /playlists", playlistHandler.CreatePlaylist)
}
