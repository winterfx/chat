package rest

import (
	"chat/internal/cache"
	"chat/internal/repository"
	"chat/internal/services/chat"
	"chat/internal/services/user"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	chatService *chat.ChatService
	userState   *user.UserState
	userService *user.UsersService
)

func LaunchApiServer() {
	InitService()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(CorsMiddleware())

	registerAuthApis(router)
	registerChatApis(router)

	log.Println("Starting server on :8081")
	err := router.Run("0.0.0.0:8081")
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func InitService() {
	pubsub := chat.NewRedisPubSubService(cache.NewRedisCache(), repository.GetDefaultChatHistoryRepo())
	chatService = chat.NewChatService(pubsub)
	userState = user.NewUserState(cache.NewRedisCache())
	userService = user.NewUsersService(repository.GetDefaultUserRepo())
}
