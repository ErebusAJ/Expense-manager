-- name: CreateUser :exec
INSERT INTO users(id, name, email, password_hash)
VALUES($1, $2, $3, $4);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT id, password_hash FROM users
WHERE email=$1;