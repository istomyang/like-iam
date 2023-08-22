package secret

import (
	"github.com/gin-gonic/gin"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Delete(ctx *gin.Context) {
	log.L(ctx).Info("delete secret.")

	if err := c.svc.Secrets().Delete(ctx, ctx.GetString(middleware.UserNameKey), ctx.Param("secret-id"), metav1.DeleteOperateMeta{Unscoped: true}); err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, nil)
}
