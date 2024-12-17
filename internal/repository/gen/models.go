// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package gen

import (
	"database/sql"
	"time"
)

type ChatHistory struct {
	ID               int32          `json:"id"`
	MessageID        string         `json:"message_id"`
	MessageTimestamp time.Time      `json:"message_timestamp"`
	SenderID         string         `json:"sender_id"`
	ReceiverID       string         `json:"receiver_id"`
	Message          string         `json:"message"`
	RetryCount       sql.NullInt32  `json:"retry_count"`
	Status           sql.NullString `json:"status"`
	CreatedAt        sql.NullTime   `json:"created_at"`
	UpdatedAt        sql.NullTime   `json:"updated_at"`
}

type User struct {
	UserID       string       `json:"user_id"`
	Email        string       `json:"email"`
	Name         string       `json:"name"`
	AuthProvider string       `json:"auth_provider"`
	PasswordHash string       `json:"password_hash"`
	CreatedAt    sql.NullTime `json:"created_at"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
}
