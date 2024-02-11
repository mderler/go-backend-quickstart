-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY username;

-- name: CreateUser :one
INSERT INTO "user" (
  username, email, password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :one
UPDATE "user"
  set username = $2,
  email = $3,
  password = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM "user"
WHERE id = $1;