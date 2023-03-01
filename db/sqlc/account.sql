-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id=$1 LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts
ORDER BY owner OFFSET $1 LIMIT $2;

-- name: UpdateAccount :one
UPDATE accounts
  set owner=$2, balance=$3
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
  set balance=balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id=$1;