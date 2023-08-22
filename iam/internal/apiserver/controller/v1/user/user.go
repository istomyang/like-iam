package user

import (
	"istomyang.github.com/like-iam/iam/internal/apiserver/service"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
)

type Controller struct {
	svc service.Service
}

func NewUserController(store store.Factory) *Controller {
	return &Controller{svc: service.NewService(store)}
}
