-- name: CreateUser :one
INSERT INTO users (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT  FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

