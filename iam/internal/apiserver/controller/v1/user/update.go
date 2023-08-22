package user

import (
	"github.com/gin-gonic/gin"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Update(ctx *gin.Context) {
	log.L(ctx).Info("get a user.")

	var user v1.User

	err := ctx.ShouldBind(&user)
	if err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrBind, err.Error()), nil)
		return
	}

	if err := c.svc.Users().Update(ctx, &user, metav1.UpdateOperateMeta{}); err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, nil)
}
