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
	handler := server.corsMiddleware(server.router)
	httpServer := &http.Server{
		Addr:         ":" + server.cfg.ApiPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	slog.Info("Server starting", "port", server.cfg.ApiPort)
	return httpServer.ListenAndServe()
}

func (server *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
