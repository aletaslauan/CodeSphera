-- name: SelectUsers :many
SELECT id, username, password_hash, role
FROM users;
-- name: AddUser :one
INSERT INTO users (username, password_hash, role)
VALUES ($1, $2, $3)
RETURNING id, username, password_hash, role;
