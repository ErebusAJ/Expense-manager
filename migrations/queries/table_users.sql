-- name: CreateUser :exec
INSERT INTO users(id, name, email, password_hash, image_url)
VALUES($1, $2, $3, $4, $5);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT id, name, password_hash, access_level FROM users
WHERE email=$1;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id=$1;

-- name: UpdateUserDetails :exec
UPDATE users 
SET name=$1, password_hash=$2, email=$3
WHERE id=$4;