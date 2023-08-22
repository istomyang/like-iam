package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Cors() gin.HandlerFunc {
	return libCors()
}

// self use MDN 's reference.
// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
func self() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "authorize, origin, content-interface, accept")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", strconv.FormatInt(int64((time.Hour*24)/time.Second), 10))
			c.Header("Content-Type", "application/json")
			c.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// libCors use lib implement.
func libCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors.New(cors.DefaultConfig())
	}
}
