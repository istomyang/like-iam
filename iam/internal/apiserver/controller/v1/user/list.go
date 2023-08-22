package user

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) List(ctx *gin.Context) {
	log.L(ctx).Info("list users.")

	var meta metav1.ListOperateMeta

	if err := ctx.ShouldBindQuery(&meta); err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrBind, err.Error()), nil)
		return
	}
	userList, err := c.svc.Users().List(ctx, meta)
	if err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, userList)
}
