-- name: CreateJob :one
INSERT INTO jobs (
  client_id, category_id, title, description,
  budget_min, budget_max, deadline
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING
  id, client_id, category_id, title, description,
  budget_min, budget_max, deadline, status,
  created_at, updated_at;

-- name: ListOpenJobs :many
SELECT
  j.id, j.client_id, j.category_id,
  j.title, j.description,
  j.budget_min, j.budget_max, j.deadline,
  j.status, j.created_at, j.updated_at
  FROM jobs j
 WHERE j.status = 'open'
 ORDER BY j.created_at DESC
 LIMIT $1 OFFSET $2;

-- name: CountOpenJobs :one
SELECT COUNT(*) FROM jobs WHERE status = 'open';


-- name: GetJobByID :one
SELECT *
  FROM jobs
 WHERE id = $1;

 -- name: UpdateJob :exec
UPDATE jobs
SET title = $2, description = $3, category_id = $4,
    budget_min = $5, budget_max = $6
WHERE id = $1;

-- name: DeleteJob :exec
DELETE FROM jobs WHERE id = $1;

