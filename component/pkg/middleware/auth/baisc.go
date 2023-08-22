package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"strings"
)

type CheckFunc = func(username, password string) bool

// Check whether BasicScheme impl Scheme or not.
// If not, must throw error.
var _ Scheme = &BasicScheme{}

type BasicScheme struct {
	check CheckFunc
}

// NewBasicScheme create a basic-authn scheme.
// checkFunc is a func to check whether username and passwd is correct.
func NewBasicScheme(checkFunc CheckFunc) *BasicScheme {
	return &BasicScheme{check: checkFunc}
}

func (b *BasicScheme) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: Basic 'admin:Admin@2021'|base64
		auths := strings.SplitN(c.GetHeader("Authorization"), " ", 2)

		if len(auths) != 2 || auths[0] != "Basic" {
			web.WriteResponse(c,
				errors.WithCode(errors.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil)
			c.Abort()
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auths[1])
		up := strings.SplitN(string(payload), ":", 2)

		if len(up) != 2 && b.check(up[0], up[1]) {
			web.WriteResponse(c,
				errors.WithCode(errors.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil)
			c.Abort()
			return
		}

		// set username represents user get authn.
		c.Set(middleware.UserNameKey, up[0])

		c.Next()
	}
}
