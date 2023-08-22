package authorize

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/iam/internal/authzserver/authorization"
)

type Controller struct {
}

func NewAuthorizeController() *Controller {
	return &Controller{}
}

func (c *Controller) Authorize(ctx *gin.Context) {
	var req ladon.Request

	if err := ctx.ShouldBind(&req); err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrBind, err.Error()), nil)
		return
	}

	req.Context["username"] = ctx.GetString("username")

	web.WriteResponse(ctx, nil, authorization.GetAuthorizator().Authorize(&req))
}
