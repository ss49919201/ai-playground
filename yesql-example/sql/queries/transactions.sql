-- name: create_transaction
INSERT INTO transactions (transaction_id, from_account, to_account, transaction_type, amount, description)
VALUES (?, ?, ?, ?, ?, ?);

-- name: get_transaction_by_id
SELECT transaction_id, from_account, to_account, transaction_type, amount, description, created_at
FROM transactions
WHERE transaction_id = ?;

-- name: get_account_transactions
SELECT transaction_id, from_account, to_account, transaction_type, amount, description, created_at
FROM transactions
WHERE from_account = ? OR to_account = ?
ORDER BY created_at DESC;

-- name: list_transactions
SELECT transaction_id, from_account, to_account, transaction_type, amount, description, created_at
FROM transactions
ORDER BY created_at DESC;

-- name: get_account_balance_for_update
SELECT balance
FROM accounts
WHERE account_id = ?
FOR UPDATE;