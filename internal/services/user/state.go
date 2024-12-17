package user

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserState struct {
	cache *redis.Client
}

type UserStatus string

const expiredTime = 60 * time.Duration(time.Minute)
const (
	Offline UserStatus = "offline"
	Online  UserStatus = "online"
)

func userStatusKey(userId string) string {
	return "user_status:" + userId
}
func (u UserStatus) IsOnline() bool {
	return u == Online
}

func (u *UserState) GetUserState(ctx context.Context, userId string) UserStatus {
	k := userStatusKey(userId)
	status, err := u.cache.Get(ctx, k).Result()
	if err != nil {
		return Offline
	}
	return UserStatus(status)
}

func (u *UserState) RefreshUserState(ctx context.Context, userId string, status UserStatus) error {
	k := userStatusKey(userId)
	return u.cache.Set(ctx, k, string(status), expiredTime).Err()
}
func NewUserState(r *redis.Client) *UserState {
	return &UserState{
		cache: r,
	}
}
