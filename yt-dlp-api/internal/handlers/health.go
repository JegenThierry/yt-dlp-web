package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const healthTimeout = 3 * time.Second

type HealthHandler struct {
	Pool *pgxpool.Pool
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthTimeout)
	defer cancel()

	if err := h.Pool.Ping(ctx); err != nil {
		slog.Error("Health check failed", "error", err)
		writeError(w, http.StatusServiceUnavailable, "database unreachable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
