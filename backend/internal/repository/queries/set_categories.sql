-- name: AddCategoryToSet :exec
INSERT INTO set_categories (set_id, category_id)
VALUES ($1, $2);

-- name: GetAllCategoriesForSet :many
SELECT categories.id, categories.name FROM set_categories JOIN categories ON
set_categories.category_id = categories.id WHERE set_categories.set_id = $1;

-- name: RemoveCategoryFromSet :exec
DELETE FROM set_categories WHERE set_id = $1 AND category_id = $2;