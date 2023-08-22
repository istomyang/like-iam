package policy

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) List(ctx *gin.Context) {
	log.L(ctx).Info("list policy.")

	web.WriteResponse(ctx, nil, nil)
}
