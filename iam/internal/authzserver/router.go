package authzserver

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/component/pkg/middleware/auth"
	"istomyang.github.com/like-iam/iam/internal/authzserver/controller/v1/authorize"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service"
)

func installRouter(engine *gin.Engine) {

	// Parse token in the gin context, put username into context.
	m := auth.NewCacheScheme(func(kid string) (*auth.Secret, error) {
		secret, err := service.GetService().FindSecret(kid)
		if err != nil {
			return nil, err
		}
		return &auth.Secret{
			Username: secret.Username,
			ID:       secret.SecretId,
			Key:      secret.SecretKey,
			Expires:  secret.Expires,
		}, nil
	}).AuthFunc()

	engine.NoRoute(m, func(ctx *gin.Context) {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrPageNotFound, "route not found"), nil)
	})

	authzCtrl := authorize.NewAuthorizeController()
	v1 := engine.Group("/v1", m)
	v1.POST("/authz", authzCtrl.Authorize)
}
