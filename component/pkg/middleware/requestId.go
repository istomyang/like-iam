package middleware

import "github.com/gin-gonic/gin"
import uu "github.com/gofrs/uuid"

const (
	// XRequestIDKey defines X-Request-ID key string.
	XRequestIDKey = "X-Request-ID"
)

// RequestId puts "X-Request-ID" into Headers which new an uuid string if not existed in Request's Headers.
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(XRequestIDKey)

		if id == "" {
			u := uuid()
			c.Request.Header.Set(XRequestIDKey, u)
			c.Set(XRequestIDKey, u)
		}

		c.Writer.Header().Set(XRequestIDKey, id)
		c.Next()
	}
}

// uuid provide generating uuid string.
func uuid() string {
	// use github.com/gofrs/uuid/v3 for:
	// 1, it's predecessor github.com/satori/go.uuid provide different algorithms.
	// 2, github.com/satori/go.uuid is no longer maintained.
	return uu.Must(uu.NewV4()).String()
}
