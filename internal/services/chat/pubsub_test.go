package chat

import (
	"chat/internal/cache"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func Test1(t *testing.T) {
	pubsub := NewRedisPubSubService(cache.NewRedisCache())
	stopChan := make(chan struct{})
	ctx := context.Background()
	subReady := make(chan struct{})
	go func() {
		err := pubsub.SubscribeUserChanelWrite(ctx, "22", func(message *redis.Message) error {
			fmt.Printf("Received message: %v\n", message.Payload)
			return nil
		}, subReady)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Wait for the subscription to be ready
	<-subReady

	err := pubsub.PublishUserChanelWrite(ctx, &PubMessage{
		UserId: "22",
		Data:   "test222",
	})
	if err != nil {
		t.Fatal(err)
	}

	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			close(stopChan)
			return
		default:
		}
	}
}
