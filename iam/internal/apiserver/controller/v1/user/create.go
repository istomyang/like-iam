package user

import (
	"github.com/gin-gonic/gin"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/auth"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/validator"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
	"time"
)

func (c *Controller) Create(ctx *gin.Context) {
	log.L(ctx).Info("create a user.")

	var r *v1.User

	if err := ctx.ShouldBind(&r); err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrBind, err.Error()), nil)
		return
	}

	if err := validator.CheckPasswordErr(r.Password); err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	r.Password, _ = auth.Encrypt(r.Password)
	r.LoginAt = time.Now()

	if err := c.svc.Users().Create(ctx, r, metav1.CreateOperateMeta{}); err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, nil)
}
