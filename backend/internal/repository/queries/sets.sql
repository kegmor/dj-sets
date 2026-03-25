-- name: CreateSet :one
INSERT INTO sets (id, video_id, title, dj_name, channel_name, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllSets :many
SELECT * FROM sets;

-- name: GetSetById :one
SELECT * FROM sets WHERE id=$1;

-- name: DeleteSetById :one
Delete FROM sets WHERE id=$1
RETURNING *;