package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/validator"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"istomyang.github.com/like-iam/component/pkg/middleware/auth"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"istomyang.github.com/like-iam/log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func GetAutoScheme() auth.Scheme {
	return auth.NewAutoScheme(GetBasicScheme().(*auth.BasicScheme), GetJwtSchemeOr(nil).(*auth.JwtScheme))
}

func GetBasicScheme() auth.Scheme {
	return auth.NewBasicScheme(func(username, password string) bool {
		user, err := store.Client().User().Get(context.TODO(), username, metav1.GetOperateMeta{})
		if err != nil {
			log.Errorf("basic error: %s", err.Error())
			return false
		}

		if user.Compare(password) {
			user.LoginAt = time.Now()
			if store.Client().User().Update(context.TODO(), user, metav1.UpdateOperateMeta{}) != nil {
				log.Errorf("basic error: %s", err.Error())
				return false
			}
			return true
		}

		return false
	})
}

var (
	jwtAuth *auth.JwtScheme
	jwtOnce sync.Once
)

// GetJwtSchemeOr should run in `create stage`, ensures can be used in GetAutoScheme.
func GetJwtSchemeOr(opts *options.JwtOpts) auth.Scheme {
	if jwtAuth == nil && opts == nil {
		_ = fmt.Errorf("jwt scheme options must not be nil")
		return nil
	}

	var jwtScheme *auth.JwtScheme

	jwtOnce.Do(func() {
		jwtScheme = auth.NewJwtScheme(&jwt.GinJWTMiddleware{
			Realm:            opts.Realm,
			SigningAlgorithm: "HS256",
			Key:              []byte(opts.Key),
			Timeout:          opts.Timeout,
			MaxRefresh:       opts.MaxRefresh,
			Authenticator:    loginAuthenticator(),
			Authorizator: func(data interface{}, c *gin.Context) bool {
				return true
			},
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				claims := jwt.MapClaims{
					"iss": "apiserver",
					"aud": "apiserver.iam.com",
				}
				if user, ok := data.(*v1.User); ok {
					claims["sub"] = user.Username
					claims[middleware.UserNameKey] = user.Username
				}
				return claims
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			},
			LoginResponse: func(c *gin.Context, i int, s string, t time.Time) {
				c.JSON(http.StatusOK, gin.H{
					"code":   http.StatusOK,
					"token":  s,
					"expire": t.Format(time.RFC3339),
				})
			},
			LogoutResponse: logout(),
			RefreshResponse: func(c *gin.Context, i int, s string, t time.Time) {
				c.JSON(http.StatusOK, gin.H{
					"code":   http.StatusOK,
					"token":  s,
					"expire": t.Format(time.RFC3339),
				})
			},
			IdentityHandler:       nil,
			TokenLookup:           "",
			TokenHeadName:         "",
			TimeFunc:              nil,
			HTTPStatusMessageFunc: nil,
			PrivKeyFile:           "",
			PubKeyFile:            "",
			SendCookie:            false,
			SecureCookie:          false,
			SendAuthorization:     false,
		})

		jwtAuth = jwtScheme
	})

	return jwtScheme
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func loginAuthenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var ln *login
		var err error
		if c.GetHeader("Authorization") != "" {
			ln, err = parseWithHeader(c)
		} else {
			ln, err = parseWithBody(c)
		}
		if err != nil {
			return nil, err
		}

		var user *v1.User
		user, err = store.Client().User().Get(c, ln.Username, metav1.GetOperateMeta{})
		if err != nil {
			return nil, err
		}

		if !validator.CheckPassword(user.Password) {
			return nil, jwt.ErrFailedAuthentication
		}

		user.LoginAt = time.Now()
		if err = store.Client().User().Update(c, user, metav1.UpdateOperateMeta{}); err != nil {
			return nil, err
		}

		return user, nil
	}
}

func logout() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
	}
}

func parseWithHeader(c *gin.Context) (*login, error) {
	// "Authorization Basic username:password"
	authArr := strings.SplitN(c.GetHeader("Authorization"), "", 2)
	if len(authArr) != 2 && authArr[0] != "Basic" {
		log.Errorf("parse token fail: %v", authArr)
		return nil, jwt.ErrFailedAuthentication
	}
	decodeString, err := base64.StdEncoding.DecodeString(authArr[2])
	if err != nil {
		log.Errorf("decode base64 token fail: %v", authArr[2])
		return nil, jwt.ErrFailedAuthentication
	}
	up := strings.SplitN(string(decodeString), "", 2)
	if len(up) != 2 {
		log.Errorf("format wrong, must username:password, got: %v", up)
		return nil, jwt.ErrFailedAuthentication
	}
	return &login{
		Username: up[0],
		Password: up[1],
	}, nil
}

func parseWithBody(c *gin.Context) (*login, error) {
	var ln *login

	if err := c.ShouldBind(&ln); err != nil {
		log.Error(jwt.ErrFailedAuthentication.Error())
		return nil, jwt.ErrFailedAuthentication
	}

	return ln, nil

}
