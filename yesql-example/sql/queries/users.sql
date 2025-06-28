-- name: create_user
INSERT INTO users (user_id, username, email, password_hash)
VALUES (?, ?, ?, ?);

-- name: get_user_by_id
SELECT user_id, username, email, password_hash, created_at, updated_at
FROM users
WHERE user_id = ?;

-- name: get_user_by_username
SELECT user_id, username, email, password_hash, created_at, updated_at
FROM users
WHERE username = ?;

-- name: get_user_by_email
SELECT user_id, username, email, password_hash, created_at, updated_at
FROM users
WHERE email = ?;

-- name: update_user_password
UPDATE users
SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?;

-- name: delete_user
DELETE FROM users
WHERE user_id = ?;

-- name: list_users
SELECT user_id, username, email, created_at, updated_at
FROM users
ORDER BY created_at;