package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/whauzan/todo-api/internal/domain"
	"github.com/whauzan/todo-api/internal/pkg/apperror"
	"github.com/whauzan/todo-api/internal/repository"
)

// TodoService handles todo business logic
type TodoService struct {
	todoRepo repository.TodoRepository
	logger   *slog.Logger
}

// NewTodoService creates a new TodoService
func NewTodoService(
	todoRepo repository.TodoRepository,
	logger *slog.Logger,
) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

// Create creates a new todo
func (s *TodoService) Create(ctx context.Context, userID uuid.UUID, req *domain.CreateTodoRequest) (*domain.Todo, error) {
	todo := &domain.Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		s.logger.ErrorContext(ctx, "failed to create todo", "error", err, "user_id", userID)
		return nil, apperror.ErrInternal
	}

	s.logger.InfoContext(ctx, "todo created successfully", "todo_id", todo.ID, "user_id", userID)

	return todo, nil
}

// GetByID retrieves a todo by ID and verifies ownership
func (s *TodoService) GetByID(ctx context.Context, userID, todoID uuid.UUID) (*domain.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, todoID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get todo by ID", "error", err, "todo_id", todoID)
		return nil, apperror.ErrInternal
	}

	if todo == nil {
		return nil, apperror.NewAppError(
			apperror.CodeNotFound,
			"Todo not found",
			404,
			fmt.Errorf("todo with ID %s not found", todoID),
		)
	}

	// Verify ownership
	if todo.UserID != userID {
		s.logger.WarnContext(ctx, "user attempted to access todo they don't own",
			"user_id", userID, "todo_id", todoID, "owner_id", todo.UserID)
		return nil, apperror.ErrForbidden
	}

	return todo, nil
}

// List retrieves all todos for a user
func (s *TodoService) List(ctx context.Context, userID uuid.UUID) ([]*domain.Todo, error) {
	todos, err := s.todoRepo.ListByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list todos", "error", err, "user_id", userID)
		return nil, apperror.ErrInternal
	}

	// Return empty slice instead of nil if no todos found
	if todos == nil {
		todos = []*domain.Todo{}
	}

	return todos, nil
}

// Update updates a todo
func (s *TodoService) Update(ctx context.Context, userID, todoID uuid.UUID, req *domain.UpdateTodoRequest) (*domain.Todo, error) {
	// First, get the todo and verify ownership
	todo, err := s.GetByID(ctx, userID, todoID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	// Save the updated todo
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		s.logger.ErrorContext(ctx, "failed to update todo", "error", err, "todo_id", todoID)
		return nil, apperror.ErrInternal
	}

	s.logger.InfoContext(ctx, "todo updated successfully", "todo_id", todoID, "user_id", userID)

	return todo, nil
}

// Delete deletes a todo
func (s *TodoService) Delete(ctx context.Context, userID, todoID uuid.UUID) error {
	// First, verify the todo exists and the user owns it
	_, err := s.GetByID(ctx, userID, todoID)
	if err != nil {
		return err
	}

	// Delete the todo
	if err := s.todoRepo.Delete(ctx, todoID); err != nil {
		s.logger.ErrorContext(ctx, "failed to delete todo", "error", err, "todo_id", todoID)
		return apperror.ErrInternal
	}

	s.logger.InfoContext(ctx, "todo deleted successfully", "todo_id", todoID, "user_id", userID)

	return nil
}
