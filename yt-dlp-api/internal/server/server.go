package server

import (
	"log/slog"
	"net/http"
	"time"
	"ytdl-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	cfg    config.Config
	db     *pgxpool.Pool
	logger *slog.Logger
	router *http.ServeMux
}

func NewServer(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *Server {
	server := &Server{
		cfg:    cfg,
		db:     db,
		logger: logger,
		router: http.NewServeMux(),
	}
	server.registerRoutes()
	return server
}

func (server *Server) Run() error {
	httpServer := &http.Server{
		Addr:         ":" + server.cfg.ApiPort,
		Handler:      server.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	slog.Info("Server starting", "port", server.cfg.ApiPort)
	return httpServer.ListenAndServe()
}
