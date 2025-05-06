-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
RETURNING id, name, created_at;

-- name: ListCategories :many
SELECT id, name, created_at
  FROM categories
 ORDER BY name;
