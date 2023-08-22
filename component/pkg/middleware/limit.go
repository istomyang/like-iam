package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

const ErrLimitExceeded = "limit exceeded"

// Limit prevents server from too many requests.
func Limit(maxEventPerSec float64, maxBurstToken int) gin.HandlerFunc {
	l := rate.NewLimiter(rate.Limit(maxEventPerSec), maxBurstToken)
	return func(c *gin.Context) {
		if l.Allow() {
			c.Next()
		} else {
			_ = c.AbortWithError(http.StatusTooManyRequests, errors.New(ErrLimitExceeded))
		}
	}
}
