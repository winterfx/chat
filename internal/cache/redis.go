package cache

import (
	"chat/internal/config"
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisCache() *redis.Client {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.RedisUrl,
		//Password: config.Config.RedisPassword,
	})
	s, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	println(s)
	return rdb
}
