package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/whauzan/todo-api/internal/pkg/apperror"
)

var validate = validator.New()

// Response is the standard envelope for all API responses
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo contains structured error information
type ErrorInfo struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// Meta contains optional metadata like pagination and request tracking
type Meta struct {
	RequestID  string      `json:"request_id,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination information for list responses
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// JSON sends a success response with data
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	}); err != nil {
		// If encoding fails, there's not much we can do at this point
		slog.Error("failed to encode response", "error", err)
	}
}

// JSONWithMeta sends a success response with data and metadata
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}); err != nil {
		slog.Error("failed to encode response with meta", "error", err)
	}
}

// JSONError sends an error response from AppError
func JSONError(w http.ResponseWriter, logger *slog.Logger, r *http.Request, err error) {
	appErr, ok := err.(*apperror.AppError)
	if !ok {
		// If it's not an AppError, treat it as internal server error
		logger.ErrorContext(r.Context(), "unexpected error", "error", err)
		appErr = apperror.ErrInternal
	}

	// Log errors that are not client errors
	if appErr.Status >= 500 {
		logger.ErrorContext(r.Context(), "server error",
			"error", appErr.Error(),
			"code", appErr.Code,
			"status", appErr.Status,
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}); err != nil {
		logger.ErrorContext(r.Context(), "failed to encode error response", "error", err)
	}
}

// JSONErrorWithStatus sends an error response with custom status
func JSONErrorWithStatus(w http.ResponseWriter, status int, code, message string, details []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}

// decodeJSON decodes a JSON request body
func decodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return apperror.NewAppError(
			apperror.CodeBadRequest,
			"Invalid JSON request body",
			http.StatusBadRequest,
			err,
		)
	}
	return nil
}

// validateStruct validates a struct using validator
func validateStruct(v interface{}) error {
	if err := validate.Struct(v); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return apperror.ErrValidation
		}
		details := formatValidationErrors(validationErrors)
		return apperror.ErrValidation.WithDetails(details...)
	}
	return nil
}

// formatValidationErrors formats validation errors into detailed messages
func formatValidationErrors(errs validator.ValidationErrors) []string {
	var details []string
	for _, e := range errs {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			details = append(details, fmt.Sprintf("%s: is required", field))
		case "email":
			details = append(details, fmt.Sprintf("%s: must be a valid email", field))
		case "min":
			details = append(details, fmt.Sprintf("%s: must be at least %s characters", field, e.Param()))
		case "max":
			details = append(details, fmt.Sprintf("%s: must be at most %s characters", field, e.Param()))
		default:
			details = append(details, fmt.Sprintf("%s: failed %s validation", field, e.Tag()))
		}
	}
	return details
}
