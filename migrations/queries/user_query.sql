-- name: GetUserById :one
SELECT user_id, email, name, auth_provider, password_hash
FROM users
WHERE user_id = ?;

-- name: GetUserByEmail :one
SELECT user_id, email, name, auth_provider, password_hash
FROM users
WHERE email = ?;

-- name: CreateUser :exec
INSERT INTO users (user_id, email, name, auth_provider, password_hash)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateUserByEmail :exec
UPDATE users
SET
    name = ?,
    auth_provider = ?,
    password_hash=?,
    updated_at = CURRENT_TIMESTAMP
WHERE email = ?;

-- name: UpdateUserById :exec
UPDATE users
SET
    name = ?,
    auth_provider = ?,
    password_hash = ?,
    email = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?;

-- name: GetAllUsers :many
SELECT user_id, email, name, auth_provider, password_hash
FROM users;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE user_id = ?;