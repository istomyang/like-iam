package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"istomyang.github.com/like-iam/log"
	"net/http"
)

// PublishInfoInterface provides essential info to publish a message.
type PublishInfoInterface interface {
	// WithContext put current gin context in struct.
	WithContext(*gin.Context)

	// Channel return redis channel string.
	Channel() string
	// Message return redis message string.
	Message() string

	// Success and Fail is result work of Publish.
	Success()
	Fail(err error)
}

// Publish 's publisher must be func because redis may not be created.
func Publish(info PublishInfoInterface, publisher func() redis.UniversalClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		client := publisher()

		if c.Writer.Status() != http.StatusOK {
			log.L(c).Debugf("response not 200 with code %d, cancel sending publish message.", c.Writer.Status())
			return
		}

		info.WithContext(c)

		// TODO: this context.
		if err := client.Publish(context.TODO(), info.Channel(), info.Message()).Err(); err != nil {
			info.Fail(err)
			return
		}
		info.Success()
	}
}
