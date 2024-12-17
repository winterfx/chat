package rest

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	UidKey = "uid"
)

func SetUid(c context.Context, uid string) context.Context {
	// Check if the context is of type *gin.Context
	if ginCtx, ok := c.(*gin.Context); ok {
		originalCtx := ginCtx.Request.Context()
		// Set the uid in the original context
		newCtx := context.WithValue(originalCtx, UidKey, uid)
		ginCtx.Request = ginCtx.Request.WithContext(newCtx)
		return newCtx
	}
	return context.WithValue(c, UidKey, uid)
}

func GetUid(c context.Context) string {
	if ginCtx, ok := c.(*gin.Context); ok {
		// 从 *gin.Context 中提取原始的 context.Context
		c = ginCtx.Request.Context()
	}

	v := c.Value(UidKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

func GetUserName(c context.Context) string {
	if ginCtx, ok := c.(*gin.Context); ok {
		c = ginCtx.Request.Context()
	}

	v := c.Value("username")
	if v == nil {
		return ""
	}
	return v.(string)
}
