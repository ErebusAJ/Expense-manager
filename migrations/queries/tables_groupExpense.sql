-- name: AddGroupExpense :one
INSERT INTO group_expense(title, description, amount, group_id, created_by)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetGroupExpenseByID :one
SELECT * FROM group_expense
WHERE id=$1;

-- name: GetAllGroupExpenses :many
SELECT * FROM group_expense
WHERE group_id=$1
ORDER BY created_at DESC;

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


-- name: GetTotalGroupExpense :one
SELECT SUM(amount)::FLOAT AS total_expense FROM group_expense
WHERE group_id=$1;


-- name: GetMembersTotalExpense :many
SELECT u.name, u.id, SUM(group_expense_participants.amount)::FLOAT AS total_expense FROM group_expense_participants
INNER JOIN users u ON u.id = group_expense_participants.user_id
INNER JOIN group_expense ON group_expense.id = group_expense_participants.group_expense_id
WHERE group_id=$1
GROUP BY u.name, u.id
ORDER BY total_expense DESC;

