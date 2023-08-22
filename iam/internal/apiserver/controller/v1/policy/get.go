package policy

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Get(ctx *gin.Context) {
	log.L(ctx).Info("get policy.")

	web.WriteResponse(ctx, nil, nil)
}
