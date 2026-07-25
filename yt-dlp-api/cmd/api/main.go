package main

import (
	"context"
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
	if err != nil {
		slog.Error("Error loading .env file", err.Error())
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", err.Error())
		return err
	}

	db, err := repository.Connect(context.Background(), cfg)
	if err != nil {
		slog.Error("Error connecting to database", err.Error())
		return err
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return server.NewServer(cfg, db, logger).Run()
}
