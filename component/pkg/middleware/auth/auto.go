package auth

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/component-base/web"
	"strings"
)

// Manager switches different schemes with Scheme interface.
type Manager struct {
	s Scheme
}

func (m *Manager) SetScheme(s Scheme) {
	m.s = s
}

func (m *Manager) AuthFunc() gin.HandlerFunc {
	return m.s.AuthFunc()
}

type AutoScheme struct {
	basic *BasicScheme
	jwt   *JwtScheme
}

func NewAutoScheme(basic *BasicScheme, jwt *JwtScheme) *AutoScheme {
	return &AutoScheme{
		basic: basic,
		jwt:   jwt,
	}
}

func (a *AutoScheme) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			web.WriteResponse(c,
				errors.WithCode(errors.ErrPermissionDenied, "Don't has `Authorization` header."),
				nil)
			c.Abort()
			return
		}

		au := strings.SplitN(h, " ", 2)
		if len(au) != 2 {
			web.WriteResponse(
				c,
				errors.WithCode(errors.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()
			return
		}

		m := Manager{}

		switch au[0] {
		case "Basic":
			m.SetScheme(a.basic)
		case "Bearer":
			m.SetScheme(a.jwt)
		default:
			web.WriteResponse(
				c,
				errors.WithCode(errors.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()
			return
		}

		// Set username-key
		m.AuthFunc()(c)

		c.Next()
	}
}
