-- name: InsertToken :exec
INSERT INTO password_tokens(id, user_id, token, expires_at)
VALUES($1, $2, $3, $4);

-- name: GetUserToken :one
SELECT * FROM password_tokens
WHERE token=$1;

-- name: DeleteToken :exec
DELETE FROM password_tokens
WHERE token=$1;


-- name: SetPassword :exec
UPDATE users
SET password_hash=$1
WHERE id=$2;

