package middleware

import "github.com/gin-gonic/gin"

// Secure ensures http security.
func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CSP
		c.Header("Content-Security-Policy", "default-src 'self'")

		// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")

		// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/X-Frame-Options
		c.Header("X-Frame-Options", "DENY")

		// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/X-XSS-Protection
		c.Header("X-XSS-Protection", "1")

		if c.Request.TLS != nil {
			// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Strict-Transport-Security
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}
