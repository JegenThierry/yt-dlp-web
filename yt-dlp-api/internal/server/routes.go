package server

import (
	"ytdl-api/internal/handlers"
	"ytdl-api/internal/repository/db"
)

func (Server *Server) registerRoutes() {
	queries := db.New(Server.db)

	// Handlers
	playlistHandler := handlers.PlaylistHandler{Queries: queries}
	downloadHandler := handlers.DownloadHandler{Queries: queries, Manager: Server.DownloadManager}

	// Playlist Routes
	Server.router.HandleFunc("POST /api/playlists", playlistHandler.CreatePlaylist)
	Server.router.HandleFunc("GET /api/playlists", playlistHandler.FetchPlaylists)
	Server.router.HandleFunc("GET /api/playlists/{id}", playlistHandler.GetPlaylist)
	Server.router.HandleFunc("POST /api/playlists/{id}/download", downloadHandler.DownloadPlaylist)
	Server.router.HandleFunc("POST /api/playlists/entries/{entryId}/download", downloadHandler.DownloadEntry)
	Server.router.HandleFunc("GET /api/playlists/entries/{entryId}/file", downloadHandler.DownloadEntryFile)
	Server.router.HandleFunc("GET /api/playlists/{id}/download/zip", downloadHandler.DownloadPlaylistZip)
	Server.router.HandleFunc("PATCH /api/playlists/entries/{entryId}/metadata", playlistHandler.UpdateEntryMetadata)
	Server.router.HandleFunc("GET /api/playlists/{id}/events", downloadHandler.GetPlaylistEvents)
}
