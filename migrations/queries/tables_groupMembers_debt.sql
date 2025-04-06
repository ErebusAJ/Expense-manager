-- name: UpdateUserDebts :one
INSERT INTO group_members_debt(from_user, to_user, group_id, expense_id, amount)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (from_user, to_user, group_id) 
DO UPDATE SET amount = group_members_debt.amount + $5
RETURNING *;


-- name: FetchNetBalance :many
WITH group_members_debt AS(
    SELECT
        e.created_by AS payer_id,
        ep.user_id AS member_id,
        e.amount AS total_amount,
        ep.amount AS share
    FROM
        group_expense e
    INNER JOIN
        group_expense_participants ep ON e.id = ep.group_expense_id
    WHERE e.group_id=$1
)
SELECT
    u.id AS user_id,
    u.name,
    u.image_url,
    COALESCE(SUM(
        CASE
            WHEN ud.member_id = ud.payer_id THEN ud.total_amount - ud.share
            ELSE -ud.share
        END
    ), 0)::NUMERIC AS netBalance
FROM
    users u
INNER JOIN
    group_members_debt ud ON u.id = ud.member_id
GROUP BY u.id;


-- name: AddSimplifiedTransaction :one
INSERT INTO simplified_transactions(group_id, from_user, to_user, amount)
VALUES($1, $2, $3, $4)
ON CONFLICT (group_id, from_user, to_user)
DO UPDATE SET amount = EXCLUDED.amount
RETURNING *;


-- name: GetSimplifiedTransactions :many
SELECT
    st.id AS transaction_id,
    st.from_user AS from_user_id,
    u_from.name AS from_user_name,
    st.to_user AS to_user_id,
    u_to.name AS to_user_name,
    st.amount
FROM simplified_transactions st
INNER JOIN users u_from ON st.from_user = u_from.id
INNER JOIN users u_to ON st.to_user = u_to.id
WHERE group_id=$1;

-- name: GetUserSimplifiedTransaction :many
SELECT
    st.id AS transaction_id,
    st.from_user AS from_user_id,
    u_from.name AS from_user_name,
    st.to_user AS to_user_id,
    u_to.name AS to_user_name,
    st.amount
FROM simplified_transactions st
INNER JOIN users u_from ON st.from_user = u_from.id
INNER JOIN users u_to ON st.to_user = u_to.id
WHERE group_id=$1 AND (u_from.id=$2 OR u_to.id=$2);


-- name: UpdateTransaction :exec
UPDATE simplified_transactions 
SET amount=$1
WHERE id=$2;

-- name: DeleteTransaction :exec
DELETE FROM simplified_transactions 
WHERE id=$1;

-- name: GetUserSettleTransaction :one
SELECT
    st.from_user AS from_user_id,
    u_from.name AS from_user_name,
    st.to_user AS to_user_id,
    u_to.name AS to_user_name,
    st.amount
FROM simplified_transactions st
INNER JOIN users u_from ON st.from_user = u_from.id
INNER JOIN users u_to ON st.to_user = u_to.id
WHERE st.id=$1;




