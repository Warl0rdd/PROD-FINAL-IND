-- name: CreateClient :one
INSERT INTO clients (id, login, age, location, gender)
VALUES ($1, $2, $3, $4, $5::gender)
ON CONFLICT (id)
    DO UPDATE SET login    = $2,
                  age      = $3,
                  location = $4,
                  gender   = $5
RETURNING *;

-- name: GetClientById :one
SELECT *
FROM clients
WHERE id = $1;