package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whauzan/todo-api/internal/config"
	"github.com/whauzan/todo-api/internal/handler"
	"github.com/whauzan/todo-api/internal/middleware"
	"github.com/whauzan/todo-api/internal/pkg/jwt"
	"github.com/whauzan/todo-api/internal/pkg/password"
	"github.com/whauzan/todo-api/internal/repository/postgres"
	"github.com/whauzan/todo-api/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger
	logger := setupLogger(cfg)
	logger.Info("starting todo-api", "env", cfg.Env, "port", cfg.Port)

	// Setup database connection
	pool, err := setupDatabase(cfg, logger)
	if err != nil {
		logger.Error("failed to setup database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize dependencies
	tokenManager := jwt.NewTokenManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	hasher := password.NewHasher()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pool)
	todoRepo := postgres.NewTodoRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenManager, hasher, logger)
	todoService := service.NewTodoService(todoRepo, logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, logger)
	todoHandler := handler.NewTodoHandler(todoService, logger)
	healthHandler := handler.NewHealthHandler(pool, logger)

	// Initialize middleware
	authMiddleware := middleware.NewAuth(tokenManager, logger)
	loggingMiddleware := middleware.NewLogging(logger)
	requestIDMiddleware := middleware.NewRequestID()
	recoverMiddleware := middleware.NewRecover(logger)

	// Setup router
	r := setupRouter(cfg, authHandler, todoHandler, healthHandler, authMiddleware, loggingMiddleware, requestIDMiddleware, recoverMiddleware)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}

// setupLogger creates and configures the logger
func setupLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// setupDatabase creates and configures the database connection pool
func setupDatabase(cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established")

	return pool, nil
}

// setupRouter configures and returns the HTTP router
func setupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	todoHandler *handler.TodoHandler,
	healthHandler *handler.HealthHandler,
	authMiddleware *middleware.Auth,
	loggingMiddleware *middleware.Logging,
	requestIDMiddleware *middleware.RequestID,
	recoverMiddleware *middleware.Recover,
) *chi.Mux {
	r := chi.NewRouter()

	// Apply global middleware
	r.Use(recoverMiddleware.Handle)
	r.Use(requestIDMiddleware.Handle)
	r.Use(loggingMiddleware.Log)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", healthHandler.Check)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
		})

		// Todo routes (protected)
		r.Route("/todos", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Get("/", todoHandler.List)
			r.Post("/", todoHandler.Create)
			r.Get("/{id}", todoHandler.GetByID)
			r.Patch("/{id}", todoHandler.Update)
			r.Delete("/{id}", todoHandler.Delete)
		})
	})

	return r
}
