-- name: create_user_account_association
INSERT INTO user_accounts (user_id, account_id)
VALUES (?, ?);

-- name: get_user_accounts
SELECT ua.user_id, ua.account_id, a.account_name, a.balance, a.created_at, a.updated_at
FROM user_accounts ua
JOIN accounts a ON ua.account_id = a.account_id
WHERE ua.user_id = ?;

-- name: get_account_user
SELECT ua.user_id, u.username, u.email
FROM user_accounts ua
JOIN users u ON ua.user_id = u.user_id
WHERE ua.account_id = ?;

-- name: check_user_account_access
SELECT COUNT(*) as count
FROM user_accounts
WHERE user_id = ? AND account_id = ?;

-- name: delete_user_account_association
DELETE FROM user_accounts
WHERE user_id = ? AND account_id = ?;

-- name: delete_user_all_accounts
DELETE FROM user_accounts
WHERE user_id = ?;