package user

import (
	"github.com/gin-gonic/gin"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Delete(ctx *gin.Context) {

	log.L(ctx).Info("delete a user.")

	err := c.svc.Users().Delete(ctx, ctx.Param("name"), metav1.DeleteOperateMeta{})
	if err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, nil)
}
