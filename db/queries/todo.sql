-- name: CreateTodo :one
INSERT INTO todos (
    id,
    user_id,
    title,
    description,
    completed
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTodoByID :one
SELECT * FROM todos
WHERE id = $1 LIMIT 1;

-- name: ListTodosByUserID :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListTodosByUserIDAndStatus :many
SELECT * FROM todos
WHERE user_id = $1 AND completed = $2
ORDER BY created_at DESC;

-- name: UpdateTodo :one
UPDATE todos
SET
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    completed = COALESCE(sqlc.narg('completed'), completed),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;

-- name: CountTodosByUserID :one
SELECT COUNT(*) FROM todos
WHERE user_id = $1;

-- name: CountCompletedTodosByUserID :one
SELECT COUNT(*) FROM todos
WHERE user_id = $1 AND completed = true;
