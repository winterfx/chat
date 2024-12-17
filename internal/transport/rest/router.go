package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAuthApis(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/:provider", Auth)
		authGroup.GET("/:provider/callback", AuthCallback)
	}
	userGroup := router.Group("/user", LogRequestMiddleware())
	{
		userGroup.POST("/login", Login)
		userGroup.POST("/logout", func(ctx *gin.Context) {
			ClearTokenCookie(ctx)
			ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
		})
		userGroup.POST("/register", Register)
		userGroup.GET("/login-status", CheckIsLogin)
		userGroup.GET("/friends", AuthMiddleware(), RetrieveFriends)
	}

}

// registerChatApis registers the chat related routes.
func registerChatApis(router *gin.Engine) {
	chatGroup := router.Group("/chat")
	{
		chatGroup.GET("/ws", AuthMiddleware(), handleWs)
	}
}
