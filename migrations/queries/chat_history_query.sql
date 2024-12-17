-- name: CreateChatHistory :exec
INSERT INTO chat_history (
    message_id, message_timestamp, sender_id, receiver_id, message, retry_count, status
) VALUES ( ?, ?, ?, ?, ?, ?, ?);

-- name: GetChatHistoryByMessageId :one
SELECT
    id, message_id, message_timestamp, sender_id, receiver_id, message, retry_count, status, created_at, updated_at
FROM
    chat_history
WHERE
    message_id = ?;

-- name: UpdateChatHistoryStatusByMessageId :exec
UPDATE
    chat_history
SET
    message = ?, retry_count = ?, status = ?, updated_at = CURRENT_TIMESTAMP
WHERE
    message_id = ?;

-- name: GetChatHistoryByReceiverIdAndStatus :many
SELECT
    id, message_id, message_timestamp, sender_id, receiver_id, message, retry_count, status, created_at, updated_at
FROM
    chat_history
WHERE
    receiver_id = ? AND status = ?
ORDER BY
    message_timestamp ASC;