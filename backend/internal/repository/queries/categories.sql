-- name: CreateCategory :one
INSERT INTO categories (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: DeleteCategoryByName :one
DELETE FROM categories WHERE name=$1
RETURNING *;