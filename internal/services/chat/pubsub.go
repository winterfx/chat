package chat

import (
	"chat/internal/repository"
	"chat/internal/repository/gen"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func channleName(userId string) string {
	return "user:" + userId
}

type ChatMessageStatus string

const (
	Pending  ChatMessageStatus = "pending"
	Sent     ChatMessageStatus = "sent"
	Received ChatMessageStatus = "received"
	Failed   ChatMessageStatus = "failed"
)

type MessageHandler func(message *redis.Message) error

func GetWebsocketMessageHandler(conn *websocket.Conn) MessageHandler {
	return func(message *redis.Message) error {
		log.Printf("message will be sent to websocket: %s", message.Payload)
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // 设置写入超时时间
		err := conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("write message error by ws closed: %v", err)
			}
			return err
		}
		return nil
		//w, err := conn.NextWriter(websocket.TextMessage)
		//if err != nil {
		//	return err
		//}
		//_, err = w.Write([]byte(message.Payload))
		//if err != nil {
		//	return err
		//}
		//if err = w.Close(); err != nil {
		//	return err
		//}
		return nil
	}
}

type PubSubService interface {
	PublishUserChanelWrite(ctx context.Context, message *ChatMessage) error
	SubscribeUserChanelWrite(ctx context.Context, userId string, f MessageHandler, subReady chan struct{}) error
	UnsubscribeUserChannel(ctx context.Context, userId string) error
}

type RedisPubSubService struct {
	Rc          *redis.Client
	HistoryRepo repository.ChatHistoryRepo
	subs        sync.Map
}

func (r *RedisPubSubService) UnsubscribeUserChannel(ctx context.Context, userId string) error {
	log.Printf("unsubscribing user %s\n", userId)
	sub, ok := r.subs.Load(userId)
	if !ok {
		log.Printf("user %s not subscribed\n", userId)
		return nil
	}
	err := sub.(*redis.PubSub).Unsubscribe(ctx, channleName(userId))
	if err != nil {
		log.Printf("unsubscribe user %s error:%s\n", userId, err)
		return err
	}
	return nil
}

func (r *RedisPubSubService) PublishUserChanelWrite(ctx context.Context, message *ChatMessage) error {
	msg, _ := json.Marshal(message)
	_, err := r.Rc.Publish(ctx, channleName(message.ReceiverId), msg).Result()
	if err != nil {
		return ErrPubMessage.WithError(err)
	}
	log.Printf("message sent to %s,from user %s,message: %s\n", message.ReceiverId, message.SenderId, message.Data)
	err = r.HistoryRepo.CreateChatHistory(ctx, gen.CreateChatHistoryParams{
		MessageID:        message.MessageId,
		MessageTimestamp: time.UnixMilli(message.TimeStamp),
		SenderID:         message.SenderId,
		ReceiverID:       message.ReceiverId,
		Message:          message.Data,
		Status: sql.NullString{
			String: string(Pending),
			Valid:  true,
		},
	})
	if err != nil {
		return ErrDBOperation.WithError(err)
	}
	return err
}

func (r *RedisPubSubService) SubscribeUserChanelWrite(ctx context.Context, userId string, f MessageHandler, subReady chan struct{}) error {
	sub := r.Rc.Subscribe(ctx, channleName(userId))
	r.subs.Store(userId, sub)
	defer r.subs.Delete(userId)
	_, err := sub.Receive(ctx)
	if err != nil {
		return ErrSubMessage.WithError(err)
	}
	close(subReady)
	for {
		select {
		case <-ctx.Done():
			log.Println("context done")
			return nil
		default:
			msg, err := sub.ReceiveMessage(ctx)
			if err != nil {
				return ErrSubMessage.WithError(err)
			}
			log.Printf("message received from %s,message: %s\n", userId, msg.Payload)
			err = f(msg)
			if err != nil {
				return ErrSubMessage.WithError(err)
			}
		}
	}
}

func NewRedisPubSubService(redis *redis.Client, historyRepo repository.ChatHistoryRepo) PubSubService {
	return &RedisPubSubService{
		Rc:          redis,
		HistoryRepo: historyRepo,
		subs:        sync.Map{},
	}
}
