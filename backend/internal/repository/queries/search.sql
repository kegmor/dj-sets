-- name: GetSetsByDjName :many
SELECT * FROM sets WHERE dj_name ILIKE '%' || $1::text || '%';

-- name: GetSetsByChannelName :many
SELECT * FROM sets WHERE channel_name ILIKE '%' || $1::text || '%';

-- name: GetSetsByTitle :many
SELECT * FROM sets WHERE title ILIKE '%' || $1::text || '%';

-- name: GetSetsByCategory :many
SELECT sets.* FROM sets
JOIN set_categories ON sets.id = set_categories.set_id
JOIN categories ON set_categories.category_id = categories.id 
WHERE categories.name ILIKE '%' || $1::text || '%';