-- name: CreateGroup :one
INSERT INTO groups(name, description, created_by)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id=$1;

-- name: UpdateGroup :exec
UPDATE groups
SET name=$1, description=$2, updated_at=CURRENT_TIMESTAMP
WHERE id=$3;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id=$1;


-- name: AddMember :exec
INSERT INTO group_members(user_id, group_id)
VALUES($1, $2);

-- name: GetGroupMembers :many
SELECT users.id, users.name, users.email, users.image_url, group_members.added_at,
CASE 
    WHEN g.created_by = users.id THEN TRUE
    ELSE FALSE
END AS is_admin
FROM users
INNER JOIN group_members ON users.id = group_members.user_id
INNER JOIN groups g ON group_members.group_id = g.id
WHERE group_id=$1;

-- name: DeleteGroupMember :exec
DELETE FROM group_members
WHERE user_id=$1 AND group_id=$2;


-- name: GetUserAllGroups :many
SELECT groups.id, name, description, created_at FROM groups
INNER JOIN group_members ON groups.id = group_members.group_id
WHERE user_id=$1;    

-- name: GetUserGroups :many
SELECT * FROM groups
WHERE created_by=$1;

-- name: CheckMemeber :one
SELECT * FROM group_members
WHERE group_id=$1 AND user_id=$2;