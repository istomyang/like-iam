package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"time"
)

var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache")
)

// Secret defines user secret key.
type Secret struct {
	Username string
	ID       string
	// Key is a secret key to encrypt plain jwt string.
	Key     string
	Expires int64
}

// CacheScheme defines a authn Scheme using cache-solution in redis and memory.
type CacheScheme struct {
	// get query Secret by kid (key id, https://www.rfc-editor.org/rfc/rfc7515#section-4.1.4)
	get func(kid string) (*Secret, error)
}

func NewCacheScheme(get func(kid string) (*Secret, error)) *CacheScheme {
	return &CacheScheme{get: get}
}

func (s *CacheScheme) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			web.WriteResponse(c, errors.WithCode(errors.ErrMissingHeader, "Authorization header cannot be empty."), nil)
			c.Abort()

			return
		}

		var token string
		_, _ = fmt.Sscanf(h, "Bearer %s", token)

		var secret *Secret

		tokenT, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {

			// Inside jwt.ParseWithClaims, use this function to check token is valid.
			// If invalid, err is not nil, will return err inside function.

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}

			var err error
			secret, err = s.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}

			// jwt.ParseWithClaims use secret.Key to valid HMAC decrypt key is correct.
			return []byte(secret.Key), nil
		})

		if err != nil || !tokenT.Valid {
			web.WriteResponse(c, errors.WithCode(errors.ErrSignatureInvalid, err.Error()), nil)
			c.Abort()

			return
		}

		if keyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			web.WriteResponse(c, errors.WithCode(errors.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		c.Set(middleware.UserNameKey, secret.Username)

		c.Next()
	}
}

func keyExpired(ts int64) bool {
	// ts will not set, default is zero.
	if ts > 0 {
		return time.Now().After(time.Unix(ts, 0))
	}
	return false
}
