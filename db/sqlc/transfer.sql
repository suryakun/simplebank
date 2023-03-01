-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id=$1 LIMIT 1;

-- name: GetTransfers :many
SELECT * FROM transfers
ORDER BY from_account_id OFFSET $1 LIMIT $2;

-- name: UpdateTransfer :one
UPDATE transfers
  set from_account_id=$2, to_account_id=$3, amount=$4
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id=$1;

