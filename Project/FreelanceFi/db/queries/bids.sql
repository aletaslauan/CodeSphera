-- name: PlaceBid :one
INSERT INTO bids (job_id, freelancer_id, amount, cover_letter)
VALUES ($1, $2, $3, $4)
RETURNING id, job_id, freelancer_id, amount, cover_letter, status, created_at;

-- name: ListBidsForJob :many
SELECT *
  FROM bids
 WHERE job_id = $1
 ORDER BY created_at DESC;
