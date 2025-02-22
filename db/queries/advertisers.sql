-- name: CreateAdvertiser :one
INSERT INTO advertisers (id, name)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET name = $2
RETURNING *;

-- name: GetAdvertiserById :one
SELECT *
FROM advertisers
WHERE id = $1;