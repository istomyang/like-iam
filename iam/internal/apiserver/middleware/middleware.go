package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/middleware"
)

func NewPublishSecretMiddleFunc() gin.HandlerFunc {
	return middleware.Publish(NewSecretPublishInfo(), func() redis.UniversalClient {
		return conn.NewRedisClientOr(nil).UniversalClient()
	})
}

func NewPublishPolicyMiddleFunc() gin.HandlerFunc {
	return middleware.Publish(NewPolicyPublishInfo(), func() redis.UniversalClient {
		return conn.NewRedisClientOr(nil).UniversalClient()
	})
}
