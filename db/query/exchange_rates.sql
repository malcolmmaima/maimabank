-- name: CreateExchangeRate :one
INSERT INTO exchange_rates (
  base_currency,
  target_currency,
  exchange_rate
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetExchangeRate :one
SELECT * FROM exchange_rates
WHERE base_currency = $1
AND target_currency = $2
LIMIT 1;

-- name: ListExchangeRates :many
SELECT * FROM exchange_rates
ORDER BY id DESC;

-- name: UpdateExchangeRate :one
UPDATE exchange_rates SET
  exchange_rate = $2
WHERE id = $1
RETURNING *;

-- name: DeleteExchangeRate :exec
DELETE FROM exchange_rates
WHERE id = $1;