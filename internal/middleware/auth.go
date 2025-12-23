package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/whauzan/todo-api/internal/pkg/apperror"
	"github.com/whauzan/todo-api/internal/pkg/jwt"
)

// ContextKey is a custom type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey ContextKey = "user_email"
)

// Auth is a middleware that validates JWT tokens
type Auth struct {
	tokenManager *jwt.TokenManager
	logger       *slog.Logger
}

// NewAuth creates a new Auth middleware
func NewAuth(tokenManager *jwt.TokenManager, logger *slog.Logger) *Auth {
	return &Auth{
		tokenManager: tokenManager,
		logger:       logger,
	}
}

// Response is the standard envelope for error responses
type Response struct {
	Success bool       `json:"success"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo contains structured error information
type ErrorInfo struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// Authenticate validates the JWT token and adds user info to context
func (a *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.writeError(w, r, apperror.ErrUnauthorized)
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.writeError(w, r, apperror.NewAppError(
				apperror.CodeUnauthorized,
				"Invalid authorization header format",
				http.StatusUnauthorized,
				nil,
			))
			return
		}

		token := parts[1]

		// Validate the token
		claims, err := a.tokenManager.ValidateToken(token)
		if err != nil {
			a.logger.WarnContext(r.Context(), "invalid token", "error", err)
			a.writeError(w, r, apperror.NewAppError(
				apperror.CodeUnauthorized,
				"Invalid or expired token",
				http.StatusUnauthorized,
				err,
			))
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, apperror.ErrUnauthorized
	}
	return userID, nil
}

// GetUserEmail extracts the user email from the request context
func GetUserEmail(ctx context.Context) (string, error) {
	email, ok := ctx.Value(UserEmailKey).(string)
	if !ok {
		return "", apperror.ErrUnauthorized
	}
	return email, nil
}

// writeError writes an error response in envelope format
func (a *Auth) writeError(w http.ResponseWriter, r *http.Request, appErr *apperror.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.ErrorContext(r.Context(), "failed to encode error response", "error", err)
	}
}
