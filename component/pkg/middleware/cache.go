package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// NoCache ensures http remove cache.
// See https://httpwg.org/specs/rfc9111.html
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// used to try to work around "old and not updated proxy cache" implementations
		// that do not understand current HTTP Caching spec directives like no-store.
		c.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
		c.Header("Expires", "Tue, 28 Feb 1970 22:22:22 GMT")
		c.Header("Last-Modified", time.Now().Format(http.TimeFormat))
		c.Next()
	}
}
