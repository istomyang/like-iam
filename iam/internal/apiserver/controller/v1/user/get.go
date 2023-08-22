package user

import (
	"github.com/gin-gonic/gin"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Get(ctx *gin.Context) {
	log.L(ctx).Info("get a user.")

	user, err := c.svc.Users().Get(ctx, ctx.Param("name"), metav1.GetOperateMeta{})
	if err != nil {
		web.WriteResponse(ctx, err, nil)

		return
	}

	web.WriteResponse(ctx, nil, user)
}
