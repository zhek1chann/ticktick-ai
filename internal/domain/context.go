package domain

import (
	"github.com/gin-gonic/gin"
)

// Context keys for storing user information in Gin context
const (
	CtxKeyUserID = "userID"
	CtxKeyRole   = "role"
)

func GetUserIDFromContext(c *gin.Context) int64 {
	v, exists := c.Get(CtxKeyUserID)
	if !exists {
		return 0
	}
	userID, ok := v.(int64)
	if !ok {
		return 0
	}
	return userID
}
