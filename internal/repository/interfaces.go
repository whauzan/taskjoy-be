package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/whauzan/todo-api/internal/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates a user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
}

// TodoRepository defines the interface for todo data operations
type TodoRepository interface {
	// Create creates a new todo
	Create(ctx context.Context, todo *domain.Todo) error

	// GetByID retrieves a todo by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)

	// ListByUserID retrieves all todos for a user
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Todo, error)

	// ListByUserIDAndStatus retrieves todos for a user filtered by completion status
	ListByUserIDAndStatus(ctx context.Context, userID uuid.UUID, completed bool) ([]*domain.Todo, error)

	// Update updates a todo
	Update(ctx context.Context, todo *domain.Todo) error

	// Delete deletes a todo
	Delete(ctx context.Context, id uuid.UUID) error
}
