package repository

import (
	"chat/internal/repository/gen"
	"context"
	"database/sql"
)

type ChatHistoryRepo interface {
	CreateChatHistory(ctx context.Context, arg gen.CreateChatHistoryParams) error
	GetChatHistoryByMessageId(ctx context.Context, messageID string) (gen.ChatHistory, error)
	GetChatHistoryByReceiverIdAndStatus(ctx context.Context, arg gen.GetChatHistoryByReceiverIdAndStatusParams) ([]gen.ChatHistory, error)
	UpdateChatHistoryStatusByMessageId(ctx context.Context, arg gen.UpdateChatHistoryStatusByMessageIdParams) error
}

type chatHistoryRepo struct {
	query *gen.Queries
}

func (c *chatHistoryRepo) CreateChatHistory(ctx context.Context, arg gen.CreateChatHistoryParams) error {
	return c.query.CreateChatHistory(ctx, arg)

}

func (c *chatHistoryRepo) GetChatHistoryByMessageId(ctx context.Context, messageID string) (gen.ChatHistory, error) {
	return c.query.GetChatHistoryByMessageId(ctx, messageID)

}

func (c *chatHistoryRepo) GetChatHistoryByReceiverIdAndStatus(ctx context.Context, arg gen.GetChatHistoryByReceiverIdAndStatusParams) ([]gen.ChatHistory, error) {
	return c.query.GetChatHistoryByReceiverIdAndStatus(ctx, arg)

}

func (c *chatHistoryRepo) UpdateChatHistoryStatusByMessageId(ctx context.Context, arg gen.UpdateChatHistoryStatusByMessageIdParams) error {
	return c.query.UpdateChatHistoryStatusByMessageId(ctx, arg)

}
func NewChatHistoryRepo(db *sql.DB) ChatHistoryRepo {
	return &chatHistoryRepo{
		query: gen.New(db),
	}
}
