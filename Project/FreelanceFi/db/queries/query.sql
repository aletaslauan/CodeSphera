-- name: GetUserByUsername :one
SELECT id, username, password_hash, role FROM users WHERE username = $1;

-- name: CreateSession :exec
INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT id, username, role FROM users WHERE id = $1;