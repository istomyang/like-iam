package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type JwtScheme struct {
	*jwt.GinJWTMiddleware
}

var _ Scheme = &JwtScheme{}

// NewJwtScheme create a scheme of jwt.
func NewJwtScheme(j *jwt.GinJWTMiddleware) *JwtScheme {
	return &JwtScheme{j}
}

func (b *JwtScheme) AuthFunc() gin.HandlerFunc {
	return b.GinJWTMiddleware.MiddlewareFunc()
}
