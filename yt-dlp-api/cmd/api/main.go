package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"ytdl-api/internal/config"
	"ytdl-api/internal/repository"
	"ytdl-api/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		slog.Error("Error loading .env file", "error", err)
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		return err
	}

	db, err := repository.Connect(context.Background(), cfg)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		return err
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := server.NewServer(cfg, db, logger)

	if err := srv.DownloadManager.ResetInterrupted(context.Background()); err != nil {
		slog.Error("Error resetting interrupted downloads", "error", err)
		return err
	}

	return srv.Run()
}
