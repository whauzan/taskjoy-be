package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/whauzan/todo-api/internal/domain"
	"github.com/whauzan/todo-api/internal/pkg/apperror"
	"github.com/whauzan/todo-api/internal/service"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest

	// Decode request body
	if err := decodeJSON(r, &req); err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Validate request
	if err := validateStruct(&req); err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Register user
	userInfo, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return created user with envelope
	JSON(w, http.StatusCreated, userInfo)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest

	// Decode request body
	if err := decodeJSON(r, &req); err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Validate request
	if err := validateStruct(&req); err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Login user
	loginResp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return token and user info with envelope
	JSON(w, http.StatusOK, loginResp)
}

// Refresh handles JWT token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		JSONError(w, h.logger, r, apperror.ErrUnauthorized)
		return
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		JSONError(w, h.logger, r, apperror.NewAppError(
			apperror.CodeUnauthorized,
			"Invalid authorization header format",
			401,
			nil,
		))
		return
	}

	token := parts[1]

	// Refresh the token
	loginResp, err := h.authService.Refresh(r.Context(), token)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return new token and user info with envelope
	JSON(w, http.StatusOK, loginResp)
}
