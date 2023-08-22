package user

import (
	"github.com/gin-gonic/gin"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/auth"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/log"
)

// ChangePasswordSchema serves as router: /:username/change-password
type ChangePasswordSchema struct {
	OldPassword string `json:"oldPassword,omitempty" binding:"required"`
	NewPassword string `json:"newPassword,omitempty" binding:"required"`
}

func (c *Controller) ChangePassword(ctx *gin.Context) {
	log.L(ctx).Info("router enters into change-password.")

	var s *ChangePasswordSchema

	var err error

	if err = ctx.ShouldBind(&c); err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrBind, err.Error()), nil)
		return
	}

	var user *v1.User
	user, err = c.svc.Users().Get(ctx, ctx.Param("name"), metav1.GetOperateMeta{})
	if err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	if !auth.Compare(user.Password, s.OldPassword) {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrPasswordIncorrect, "password is incorrect."), nil)
		return
	}

	user.Password, err = auth.Encrypt(s.NewPassword)
	if err != nil {
		web.WriteResponse(ctx, errors.WithCode(errors.ErrEncrypt, "password encrypt fail."), nil)
		return
	}

	if err = c.svc.Users().Update(ctx, user, metav1.UpdateOperateMeta{}); err != nil {
		web.WriteResponse(ctx, err, nil)
		return
	}

	web.WriteResponse(ctx, nil, nil)
}
