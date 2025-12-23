package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/whauzan/todo-api/internal/pkg/apperror"
)

// Recover is a middleware that recovers from panics
type Recover struct {
	logger *slog.Logger
}

// NewRecover creates a new Recover middleware
func NewRecover(logger *slog.Logger) *Recover {
	return &Recover{
		logger: logger,
	}
}

// Handle recovers from panics and logs them
func (rec *Recover) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				rec.logger.ErrorContext(r.Context(),
					"panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"path", r.URL.Path,
					"method", r.Method,
				)

				// Return internal server error in envelope format
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := Response{
					Success: false,
					Error: &ErrorInfo{
						Code:    string(apperror.CodeInternal),
						Message: "An unexpected error occurred",
					},
				}

				if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
					rec.logger.ErrorContext(r.Context(), "failed to encode panic response", "error", encodeErr)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
