package apperror

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	CodeValidation         ErrorCode = "VALIDATION_ERROR"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeUserExists         ErrorCode = "USER_EXISTS"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeForbidden          ErrorCode = "FORBIDDEN"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
	CodeBadRequest         ErrorCode = "BAD_REQUEST"
)

// AppError represents an application error
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Status  int       `json:"-"`
	Details []string  `json:"-"`
	Err     error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code ErrorCode, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// WithDetails returns a copy of the error with details added
func (e *AppError) WithDetails(details ...string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Status:  e.Status,
		Details: details,
		Err:     e.Err,
	}
}

// Predefined errors
var (
	ErrInvalidCredentials = &AppError{
		Code:    CodeInvalidCredentials,
		Message: "Invalid email or password",
		Status:  http.StatusUnauthorized,
	}

	ErrUserExists = &AppError{
		Code:    CodeUserExists,
		Message: "User with this email already exists",
		Status:  http.StatusConflict,
	}

	ErrNotFound = &AppError{
		Code:    CodeNotFound,
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}

	ErrForbidden = &AppError{
		Code:    CodeForbidden,
		Message: "You don't have permission to access this resource",
		Status:  http.StatusForbidden,
	}

	ErrUnauthorized = &AppError{
		Code:    CodeUnauthorized,
		Message: "Authentication required",
		Status:  http.StatusUnauthorized,
	}

	ErrInternal = &AppError{
		Code:    CodeInternal,
		Message: "An unexpected error occurred",
		Status:  http.StatusInternalServerError,
	}

	ErrValidation = &AppError{
		Code:    CodeValidation,
		Message: "Validation failed",
		Status:  http.StatusBadRequest,
	}

	ErrBadRequest = &AppError{
		Code:    CodeBadRequest,
		Message: "Bad request",
		Status:  http.StatusBadRequest,
	}
)

// ErrorResponse represents the JSON error response structure
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents the error detail in the response
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details []string  `json:"details,omitempty"`
}

// ToErrorResponse converts an AppError to an ErrorResponse
func (e *AppError) ToErrorResponse() ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		},
	}
}
