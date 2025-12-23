package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whauzan/todo-api/internal/domain"
	"github.com/whauzan/todo-api/internal/repository/postgres/db"
)

// TodoRepository implements the repository.TodoRepository interface
type TodoRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(pool *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// Create creates a new todo
func (r *TodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	var description sql.NullString
	if todo.Description != nil {
		description = sql.NullString{String: *todo.Description, Valid: true}
	}

	params := db.CreateTodoParams{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: description,
		Completed:   todo.Completed,
	}

	dbTodo, err := r.queries.CreateTodo(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	// Update the todo with generated values
	todo.CreatedAt = dbTodo.CreatedAt
	todo.UpdatedAt = dbTodo.UpdatedAt

	return nil
}

// GetByID retrieves a todo by ID
func (r *TodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	dbTodo, err := r.queries.GetTodoByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get todo by ID: %w", err)
	}

	return r.toDomainTodo(dbTodo), nil
}

// ListByUserID retrieves all todos for a user
func (r *TodoRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Todo, error) {
	dbTodos, err := r.queries.ListTodosByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos by user ID: %w", err)
	}

	todos := make([]*domain.Todo, 0, len(dbTodos))
	for _, dbTodo := range dbTodos {
		todos = append(todos, r.toDomainTodo(dbTodo))
	}

	return todos, nil
}

// ListByUserIDAndStatus retrieves todos for a user filtered by completion status
func (r *TodoRepository) ListByUserIDAndStatus(ctx context.Context, userID uuid.UUID, completed bool) ([]*domain.Todo, error) {
	params := db.ListTodosByUserIDAndStatusParams{
		UserID:    userID,
		Completed: completed,
	}

	dbTodos, err := r.queries.ListTodosByUserIDAndStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos by user ID and status: %w", err)
	}

	todos := make([]*domain.Todo, 0, len(dbTodos))
	for _, dbTodo := range dbTodos {
		todos = append(todos, r.toDomainTodo(dbTodo))
	}

	return todos, nil
}

// Update updates a todo
func (r *TodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	var description sql.NullString
	if todo.Description != nil {
		description = sql.NullString{String: *todo.Description, Valid: true}
	}

	params := db.UpdateTodoParams{
		ID:          todo.ID,
		Title:       sql.NullString{String: todo.Title, Valid: true},
		Description: description,
		Completed:   sql.NullBool{Bool: todo.Completed, Valid: true},
	}

	dbTodo, err := r.queries.UpdateTodo(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("failed to update todo: %w", err)
	}

	// Update the todo with new values
	todo.UpdatedAt = dbTodo.UpdatedAt

	return nil
}

// Delete deletes a todo
func (r *TodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteTodo(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

// toDomainTodo converts a db.Todo to domain.Todo
func (r *TodoRepository) toDomainTodo(dbTodo db.Todo) *domain.Todo {
	var description *string
	if dbTodo.Description.Valid {
		description = &dbTodo.Description.String
	}

	return &domain.Todo{
		ID:          dbTodo.ID,
		UserID:      dbTodo.UserID,
		Title:       dbTodo.Title,
		Description: description,
		Completed:   dbTodo.Completed,
		CreatedAt:   dbTodo.CreatedAt,
		UpdatedAt:   dbTodo.UpdatedAt,
	}
}
