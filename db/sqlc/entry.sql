-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id=$1 LIMIT 1;

-- name: GetEntries :many
SELECT * FROM entries
ORDER BY account_id OFFSET $1 LIMIT $2;

-- name: UpdateEntry :one
UPDATE entries
  set amount=$2
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id=$1;
