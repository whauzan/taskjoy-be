package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(pool *pgxpool.Pool, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		pool:   pool,
		logger: logger,
	}
}

// HealthData represents the health check response data
type HealthData struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Time     string `json:"time"`
}

// Check handles health check requests
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Check database connection
	dbStatus := "healthy"
	if err := h.pool.Ping(ctx); err != nil {
		h.logger.ErrorContext(ctx, "database health check failed", "error", err)
		dbStatus = "unhealthy"
	}

	status := "healthy"
	statusCode := http.StatusOK

	if dbStatus == "unhealthy" {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	healthData := HealthData{
		Status:   status,
		Database: dbStatus,
		Time:     time.Now().UTC().Format(time.RFC3339),
	}

	// Return health data with envelope
	JSON(w, statusCode, healthData)
}
