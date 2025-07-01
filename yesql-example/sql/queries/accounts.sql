-- name: create_account
INSERT INTO accounts (account_id, account_name, balance)
VALUES (?, ?, ?);

-- name: get_account_by_id
SELECT account_id, account_name, balance, created_at, updated_at
FROM accounts
WHERE account_id = ?;

-- name: get_account_balance
SELECT balance
FROM accounts
WHERE account_id = ?;

-- name: update_account_balance
UPDATE accounts
SET balance = ?, updated_at = CURRENT_TIMESTAMP
WHERE account_id = ?;

-- name: list_accounts
SELECT account_id, account_name, balance, created_at, updated_at
FROM accounts
ORDER BY created_at;

-- name: deposit_update_balance
UPDATE accounts 
SET balance = balance + ?, updated_at = CURRENT_TIMESTAMP
WHERE account_id = ?;

-- name: deposit_create_transaction
INSERT INTO transactions (transaction_id, from_account, to_account, transaction_type, amount, description)
VALUES (?, NULL, ?, 'deposit', ?, ?)
RETURNING transaction_id, from_account, to_account, transaction_type, amount, description, created_at;

-- name: withdraw_update_balance
UPDATE accounts 
SET balance = balance - ?, updated_at = CURRENT_TIMESTAMP
WHERE account_id = ? AND balance >= ?;

-- name: withdraw_create_transaction
INSERT INTO transactions (transaction_id, from_account, to_account, transaction_type, amount, description)
VALUES (?, ?, NULL, 'withdrawal', ?, ?)
RETURNING transaction_id, from_account, to_account, transaction_type, amount, description, created_at;