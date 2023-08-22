package service

import "istomyang.github.com/like-iam/iam/internal/apiserver/store"

type Service interface {
	Users() UserSvc
	Secrets() SecretSvc
	Policies() PolicySvc
}

type service struct {
	store store.Factory
}

func NewService(factory store.Factory) Service {
	return &service{store: factory}
}

func (s *service) Users() UserSvc {
	return newUserSvc(s)
}

func (s *service) Secrets() SecretSvc {
	return newSecretSvc(s)
}

func (s *service) Policies() PolicySvc {
	return newPolicySvc(s)
}
