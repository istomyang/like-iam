package policy

import (
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) DeleteCollection(ctx *gin.Context) {
	log.L(ctx).Info("delete policy collection.")

	web.WriteResponse(ctx, nil, nil)
}
