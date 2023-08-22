package service

import (
	"context"
	"istomyang.github.com/like-iam/iam-sdk-go/service/iam"
)

type Services interface {
	Iam() iam.IAM
	// Other Services
}

type services struct {
	ctx    context.Context
	config *Config
}

func NewService(ctx context.Context, config *Config) Services {
	return &services{
		ctx:    ctx,
		config: config,
	}
}

func (s *services) Iam() iam.IAM {
	return iam.NewIam(s.ctx, s.config.IAM)
}

var _ Services = &services{}
