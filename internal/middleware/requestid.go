package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestID is a middleware that generates and propagates request IDs
type RequestID struct{}

// NewRequestID creates a new RequestID middleware
func NewRequestID() *RequestID {
	return &RequestID{}
}

// Handle adds a request ID to the context and response headers
func (rid *RequestID) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			// Generate a new request ID
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Add request ID to response header
		w.Header().Set(RequestIDHeader, requestID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID extracts the request ID from the context
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}
