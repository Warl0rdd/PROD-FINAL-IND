-- name: CreateClient :one
INSERT INTO clients (id, login, age, location, gender) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetClientById :one
SELECT * FROM clients WHERE id = $1;