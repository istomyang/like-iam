package auth

import "github.com/gin-gonic/gin"

// Scheme defines interface to use auth scheme.
// https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml
type Scheme interface {
	AuthFunc() gin.HandlerFunc
}
