package rest

import (
	"os"

	"github.com/gin-gonic/gin"
)

const ChatTokenCookieName = "chat_token"

func SetTokenCookie(c *gin.Context, token string) {

	c.SetCookie(ChatTokenCookieName, token, 3600, "/", os.Getenv("DOMAIN"), false, true)
}

func GetTokenCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie(ChatTokenCookieName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ClearTokenCookie(c *gin.Context) {
	c.SetCookie(ChatTokenCookieName, "", -1, "/", os.Getenv("DOMAIN"), false, true)
}
