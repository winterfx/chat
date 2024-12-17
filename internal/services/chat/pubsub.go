package chat

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type PubMessage struct {
	UserId string `json:"userId"`
	Data   string `json:"data"`
}

func channleName(userId string) string {
	return "user:" + userId
}

type PubSubService interface {
	PublishUserChanelWrite(ctx context.Context, message *PubMessage) error
	SubscribeUserChanelWrite(ctx context.Context, userChannel string, conn *websocket.Conn)
}

type RedisPubSubService struct {
	Rc *redis.Client
}

func (r *RedisPubSubService) PublishUserChanelWrite(ctx context.Context, message *PubMessage) error {
	_, err := r.Rc.Publish(ctx, channleName(message.UserId), message.Data).Result()
	return err
}

func (r *RedisPubSubService) SubscribeUserChanelWrite(ctx context.Context, userId string, conn *websocket.Conn) {
	pubsub := r.Rc.Subscribe(ctx, channleName(userId))
	defer pubsub.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		}
	}
}

func NewRedisPubSubService(redis *redis.Client) PubSubService {
	return &RedisPubSubService{Rc: redis}
}
