package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChatService struct {
	pubsub PubSubService
}
type ChatMessage struct {
	MessageId  string `json:"messageId"`
	ReceiverId string `json:"receiverId"`
	SenderId   string `json:"senderId"`
	Data       string `json:"data"`
	TimeStamp  int64  `json:"timeStamp"`
}

type Channel struct {
	ChannelId string `json:"channelId"`
	Type      string `json:"type"`
}

type Deliverer struct {
	DeliverId string `json:"deliverId"`
	Type      string `json:"type"`
}

type ChannelMessage struct {
	MessageId string    `json:"messageId"`
	Channel   Channel   `json:"channel"`
	Deliverer Deliverer `json:"deliverer"`
	Data      string    `json:"data"`
}

func NewChatService(pubsub PubSubService) *ChatService {
	return &ChatService{pubsub: pubsub}
}

func (c *ChatService) PubMessage(ctx context.Context, message *ChatMessage) error {
	if message.TimeStamp == 0 {
		message.TimeStamp = time.Now().Unix()
	}
	if message.MessageId == "" {
		message.MessageId = uuid.New().String()
	}
	return c.pubsub.PublishUserChanelWrite(ctx, message)
}

func (c *ChatService) SubscribeUserChanel(ctx context.Context, userId string, handler MessageHandler, subReady chan struct{}) {
	err := c.pubsub.SubscribeUserChanelWrite(ctx, userId, handler, subReady)
	if err != nil {
		fmt.Printf("subscribe user %s error:%s\n", userId, err)
	}
}

func (c *ChatService) UnsubscribeUserChannel(ctx context.Context, userId string) error {
	return c.pubsub.UnsubscribeUserChannel(ctx, userId)
}
