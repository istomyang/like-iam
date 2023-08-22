package policy

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Delete(ctx *gin.Context) {
	log.L(ctx).Info("delete policy.")

	web.WriteResponse(ctx, nil, nil)
}
