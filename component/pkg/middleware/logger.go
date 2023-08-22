package middleware

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/log"
)

// UserNameKey defines username key string.
const UserNameKey = "username"

// Logger puts XRequestIDKey and UserNameKey 's value into Context with logger's key.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.XRequestIDKey, c.GetString(XRequestIDKey))
		c.Set(log.UserNameKey, c.GetString(UserNameKey))
		c.Next()
	}
}
