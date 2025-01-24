-- name: AddGroupExpense :one
INSERT INTO group_expense(title, description, amount, group_id, created_by)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetGroupExpenseByID :one
SELECT * FROM group_expense
WHERE id=$1;

-- name: GetAllGroupExpenses :many
SELECT * FROM group_expense
WHERE group_id=$1;

-- name: UpdateGroupExpense :exec
UPDATE group_expense
SET title=$1, description=$2, amount=$3, updated_at=CURRENT_TIMESTAMP
WHERE id=$4;

-- name: DeleteGroupExpense :exec
DELETE FROM group_expense
WHERE id=$1;

-- name: AddGroupExpenseMembers :one
INSERT INTO group_expense_participants(group_expense_id, user_id, amount)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetGroupExpenseMembersByID :many
SELECT * FROM group_expense_participants
WHERE group_expense_id=$1;

-- name: UpdateGroupExpenseMembers :exec
UPDATE group_expense_participants
SET amount=$1, updated_at=CURRENT_TIMESTAMP
WHERE id=$2;





