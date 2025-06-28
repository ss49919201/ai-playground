-- name: create_session
INSERT INTO sessions (session_id, user_id, expires_at)
VALUES (?, ?, ?);

-- name: get_session
SELECT session_id, user_id, expires_at, created_at
FROM sessions
WHERE session_id = ?;

-- name: get_valid_session
SELECT session_id, user_id, expires_at, created_at
FROM sessions
WHERE session_id = ? AND expires_at > CURRENT_TIMESTAMP;

-- name: delete_session
DELETE FROM sessions
WHERE session_id = ?;

-- name: delete_user_sessions
DELETE FROM sessions
WHERE user_id = ?;

-- name: delete_expired_sessions
DELETE FROM sessions
WHERE expires_at <= CURRENT_TIMESTAMP;

-- name: get_user_sessions
SELECT session_id, user_id, expires_at, created_at
FROM sessions
WHERE user_id = ? AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;