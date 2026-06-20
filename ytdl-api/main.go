package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"ytdl-api/database"
	"ytdl-api/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.Connect()
	defer db.Close()

	playlistHandler := &handlers.PlaylistHandler{DB: db}
	downloadHandler := &handlers.DownloadHandler{DB: db}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/playlists", playlistHandler.FetchPlaylists)
	mux.HandleFunc("GET /api/downloads", downloadHandler.FetchDownloads)

	port := os.Getenv("API_PORT")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
