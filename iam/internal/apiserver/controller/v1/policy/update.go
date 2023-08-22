package policy

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Update(ctx *gin.Context) {
	log.L(ctx).Info("update policy.")

	web.WriteResponse(ctx, nil, nil)
}
