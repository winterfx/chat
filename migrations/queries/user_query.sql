-- name: GetUser :one
SELECT * FROM users WHERE user_id = $1;

-- name: CreateUser :one
INSERT INTO users (user_id, username, email, oauth_token, refresh_token, token_expiry) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateUserTokens :one
UPDATE users
SET
    oauth_token = \$2,
    refresh_token = \$3,
    token_expiry = \$4,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = \$1
    RETURNING user_id, oauth_token, refresh_token, token_expiry, updated_at;