-- name: RecordTransaction :one
INSERT INTO transactions (user_id, job_id, bid_id, type, amount)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, job_id, bid_id, type, amount, recorded_at;

-- name: ListTransactionsForUser :many
SELECT *
  FROM transactions
 WHERE user_id = $1
 ORDER BY recorded_at DESC
 LIMIT $2 OFFSET $3;
