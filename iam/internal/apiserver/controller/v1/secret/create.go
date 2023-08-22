package secret

import (
	"github.com/gin-gonic/gin"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/util/idutil"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"istomyang.github.com/like-iam/log"
)

func (c *Controller) Create(ctx *gin.Context) {

	log.L(ctx).Info("create secret.")

	var secret *v1.Secret

	if err := ctx.ShouldBind(&secret); err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrBind, err.Error()), nil)
		return
	}

	secret.Username = ctx.GetString(middleware.UserNameKey)

	secret.SecretID, _ = idutil.GetRandString(idutil.AlphabetL+idutil.AlphabetU+idutil.Number, 36)
	secret.SecretKey, _ = idutil.GetRandString(idutil.AlphabetL+idutil.AlphabetU+idutil.Number, 32)

	if err := c.svc.Secrets().Create(ctx, secret, metav1.CreateOperateMeta{}); err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, nil)
}
