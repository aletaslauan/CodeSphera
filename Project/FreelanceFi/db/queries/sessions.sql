-- name: CreateUserSession :exec
INSERT INTO sessions (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetSessionByTokenHash :one
SELECT user_id, token_hash, expires_at FROM sessions
WHERE token_hash = $1;
