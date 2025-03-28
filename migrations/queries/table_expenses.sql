-- name: AddExpense :exec
INSERT INTO expenses(user_id, amount, title, description)
VALUES($1, $2, $3, $4);

-- name: GetAllExpenses :many
SELECT id, amount, title, description, created_at, updated_at FROM expenses
WHERE user_id=$1
ORDER BY created_at DESC;

-- name: GetExpenseByID :one
SELECT id, amount, title, description, created_at, updated_at FROM expenses
WHERE id=$1;

-- name: UpdateExpense :one
UPDATE expenses
SET amount=$1, title=$2, description=$3, updated_at=CURRENT_TIMESTAMP
WHERE id=$4
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id=$1;

-- name: TotalExpense :one
SELECT sum(amount)::FLOAT AS total_expense FROM expenses
WHERE user_id=$1;   




