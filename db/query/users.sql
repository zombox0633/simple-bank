-- name: CreateUser :one
INSERT INTO users (username, password, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT username, full_name, email, created_at, updated_at FROM users
ORDER BY username
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2, updated_at = now()
WHERE username = $1
RETURNING *;

-- name: ChangePassword :exec
UPDATE users
SET password = $2, updated_at = now()
WHERE username = $1;
