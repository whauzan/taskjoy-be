package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/whauzan/todo-api/internal/domain"
	"github.com/whauzan/todo-api/internal/middleware"
	"github.com/whauzan/todo-api/internal/pkg/apperror"
	"github.com/whauzan/todo-api/internal/service"
)

// TodoHandler handles todo requests
type TodoHandler struct {
	todoService *service.TodoService
	logger      *slog.Logger
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(todoService *service.TodoService, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Create handles todo creation
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	var req domain.CreateTodoRequest

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

	// Create todo
	todo, err := h.todoService.Create(r.Context(), userID, &req)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return created todo with envelope
	JSON(w, http.StatusCreated, todo)
}

// List handles listing all todos for a user
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// List todos
	todos, err := h.todoService.List(r.Context(), userID)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return todos with envelope
	JSON(w, http.StatusOK, todos)
}

// GetByID handles getting a single todo
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		JSONError(w, h.logger, r, apperror.NewAppError(
			apperror.CodeBadRequest,
			"Invalid todo ID",
			http.StatusBadRequest,
			err,
		))
		return
	}

	// Get todo
	todo, err := h.todoService.GetByID(r.Context(), userID, todoID)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return todo with envelope
	JSON(w, http.StatusOK, todo)
}

// Update handles updating a todo
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		JSONError(w, h.logger, r, apperror.NewAppError(
			apperror.CodeBadRequest,
			"Invalid todo ID",
			http.StatusBadRequest,
			err,
		))
		return
	}

	var req domain.UpdateTodoRequest

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

	// Update todo
	todo, err := h.todoService.Update(r.Context(), userID, todoID, &req)
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return updated todo with envelope
	JSON(w, http.StatusOK, todo)
}

// Delete handles deleting a todo
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		JSONError(w, h.logger, r, apperror.NewAppError(
			apperror.CodeBadRequest,
			"Invalid todo ID",
			http.StatusBadRequest,
			err,
		))
		return
	}

	// Delete todo
	if err := h.todoService.Delete(r.Context(), userID, todoID); err != nil {
		JSONError(w, h.logger, r, err)
		return
	}

	// Return success message with envelope
	JSON(w, http.StatusOK, map[string]string{
		"message": "Todo deleted successfully",
	})
}
