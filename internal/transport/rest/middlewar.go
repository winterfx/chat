package rest

import (
	"bytes"
	"chat/internal/services/user"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := GetTokenCookie(c)
		if err != nil || token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		uid, err := user.ValidateToken(c, token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		SetUid(c, uid)
		c.Next()
	}
}

// allow cors
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的Origin
		origin := c.Request.Header.Get("Origin")

		// 设置允许的Origin
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 设置为你的前端URL
		}

		// 设置允许的方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 设置允许发送Cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

type WriteResponseBody struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *WriteResponseBody) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LogRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取请求主体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 恢复请求主体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 记录请求信息
		log.Printf("Request: %s %s\nHeaders: %v\nBody: %s\n", c.Request.Method, c.Request.URL, c.Request.Header, string(bodyBytes))
		//记录cookie
		log.Printf("Cookies: %v\n", c.Request.Cookies())

		// 包装 ResponseWriter 以记录响应主体
		responseBody := &WriteResponseBody{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBody

		// 处理请求
		c.Next()

		// 记录响应信息
		log.Printf("Response: %d %s\nHeaders: %v\nBody: %s\n", c.Writer.Status(), http.StatusText(c.Writer.Status()), c.Writer.Header(), responseBody.body.String())
	}
}
