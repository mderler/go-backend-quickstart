-- name: ListTodos :many
SELECT * FROM todo
ORDER BY created_at DESC;

-- name: GetTodosOfUser :many
SELECT todo.id, creator_id, title, description, completed, created_at, updated_at FROM todo
JOIN todo_user ON todo.id = todo_user.todo_id
WHERE todo_user.user_id = $1;

-- name: CreateTodo :one
INSERT INTO todo (title, description, creator_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AssignUserToTodo :execrows
INSERT INTO todo_user (todo_id, user_id)
VALUES ($1, $2);

-- name: UpdateTodo :one
UPDATE todo
SET title = $1, description = $2, completed = $3
WHERE id = $4
RETURNING *;

-- name: DeleteTodo :execrows
DELETE FROM todo
WHERE id = $1;