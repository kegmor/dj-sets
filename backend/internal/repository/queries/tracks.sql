-- name: CreateTrack :one
INSERT INTO tracks (id, set_id, song_name, artist, timestamp_in_set, position)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTracksForSet :many
SELECT * FROM tracks WHERE set_id = $1;

-- name: DeleteTrackFromSet :one
DELETE FROM tracks WHERE id = $1
RETURNING *;