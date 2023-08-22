package apiserver

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/component-base/web"
	auth2 "istomyang.github.com/like-iam/component/pkg/middleware/auth"
	"istomyang.github.com/like-iam/iam/internal/apiserver/auth"
	"istomyang.github.com/like-iam/iam/internal/apiserver/controller/v1/policy"
	"istomyang.github.com/like-iam/iam/internal/apiserver/controller/v1/secret"
	"istomyang.github.com/like-iam/iam/internal/apiserver/controller/v1/user"
	"istomyang.github.com/like-iam/iam/internal/apiserver/middleware"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
)

func installRouter(g *gin.Engine) {

	jwtCtrl := auth.GetJwtSchemeOr(nil).(*auth2.JwtScheme)
	g.GET("/login", jwtCtrl.LoginHandler)
	g.GET("/logout", jwtCtrl.LogoutHandler)
	g.GET("/refresh", jwtCtrl.RefreshHandler)

	g.NoRoute(auth.GetAutoScheme().AuthFunc(), func(c *gin.Context) {
		web.WriteResponse(c, errors.WithCode(errors.ErrPageNotFound, "page not found"), nil)
	})

	v1 := g.Group("/v1")
	{
		userCtrl := user.NewUserController(store.Client())

		users := v1.Group("/users")
		users.POST("", userCtrl.Create)

		users.Use(auth.GetAutoScheme().AuthFunc())
		users.GET("", userCtrl.List)
		users.GET(":name", userCtrl.Get)
		users.PUT(":name", userCtrl.Update)
		users.PUT(":name/change-password", userCtrl.ChangePassword)
		users.DELETE("", userCtrl.DeleteCollection)
		users.DELETE(":name", userCtrl.Delete)
	}

	v1.Use(auth.GetAutoScheme().AuthFunc())

	{
		policyCtrl := policy.NewPolicyController(store.Client())

		policies := v1.Group("/policies", middleware.NewPublishPolicyMiddleFunc())
		policies.POST("", policyCtrl.Create)
		policies.GET("", policyCtrl.List)
		policies.GET(":name", policyCtrl.Get)
		policies.PUT("", policyCtrl.Update)
		policies.DELETE("", policyCtrl.DeleteCollection)
		policies.DELETE(":name", policyCtrl.Delete)
	}

	{
		secretCtrl := secret.NewSecretController(store.Client())

		secrets := v1.Group("/secrets", middleware.NewPublishSecretMiddleFunc())
		secrets.POST("", secretCtrl.Create)
		secrets.GET("", secretCtrl.List)
		secrets.GET(":name", secretCtrl.Get)
		secrets.PUT("", secretCtrl.Update)
		secrets.DELETE("", secretCtrl.DeleteCollection)
		secrets.DELETE(":name", secretCtrl.Delete)
	}
}
