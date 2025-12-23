package domain

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item
type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest represents the request to create a new todo
type CreateTodoRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
}

// UpdateTodoRequest represents the request to update a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	Completed   *bool   `json:"completed"`
}
