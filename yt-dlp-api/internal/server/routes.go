package server

import (
	"ytdl-api/internal/handlers"
	"ytdl-api/internal/repository/db"
	"ytdl-api/pkg/ytdlp"
)

func (Server *Server) registerRoutes() {
	queries := db.New(Server.db)

	playlistHandler := handlers.PlaylistHandler{
		Queries: queries,
		Pool:    Server.db,
		Manager: Server.DownloadManager,
		Cookies: ytdlp.Cookies{File: Server.cfg.CookiesFile, FromBrowser: Server.cfg.CookiesFromBrowser},
	}
	downloadHandler := handlers.DownloadHandler{Queries: queries, Manager: Server.DownloadManager}
	healthHandler := handlers.HealthHandler{Pool: Server.db}

	Server.router.HandleFunc("GET /api/health", healthHandler.Health)
	Server.router.HandleFunc("POST /api/playlists", playlistHandler.CreatePlaylist)
	Server.router.HandleFunc("GET /api/playlists", playlistHandler.FetchPlaylists)
	Server.router.HandleFunc("GET /api/playlists/{id}", playlistHandler.GetPlaylist)
	Server.router.HandleFunc("DELETE /api/playlists/{id}", playlistHandler.DeletePlaylist)
	Server.router.HandleFunc("POST /api/playlists/{id}/refresh", playlistHandler.RefreshPlaylist)
	Server.router.HandleFunc("POST /api/playlists/{id}/download", downloadHandler.DownloadPlaylist)
	Server.router.HandleFunc("POST /api/playlists/{id}/download/cancel", downloadHandler.CancelPlaylistDownload)
	Server.router.HandleFunc("POST /api/playlists/entries/{entryId}/download", downloadHandler.DownloadEntry)
	Server.router.HandleFunc("POST /api/playlists/entries/{entryId}/download/cancel", downloadHandler.CancelEntryDownload)
	Server.router.HandleFunc("GET /api/playlists/entries/{entryId}/file", downloadHandler.DownloadEntryFile)
	Server.router.HandleFunc("GET /api/playlists/{id}/download/zip", downloadHandler.DownloadPlaylistZip)
	Server.router.HandleFunc("PATCH /api/playlists/{id}/entries/metadata", playlistHandler.UpdatePlaylistMetadata)
	Server.router.HandleFunc("GET /api/playlists/{id}/events", downloadHandler.GetPlaylistEvents)
}
