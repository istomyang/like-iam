package service

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
)

type SecretSvc interface {
	Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOperateMeta) error
	Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOperateMeta) error
	Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOperateMeta) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOperateMeta) error
	Get(ctx context.Context, username, secretID string, opts metav1.GetOperateMeta) (*v1.Secret, error)
	List(ctx context.Context, username string, opts metav1.ListOperateMeta) (*v1.SecretList, error)
}

type secretSvc struct {
	svc *service
}

func newSecretSvc(svc *service) SecretSvc {
	return &secretSvc{svc: svc}
}

func (s *secretSvc) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOperateMeta) error {
	return s.svc.store.Secret().Create(ctx, secret, opts)
}

func (s *secretSvc) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOperateMeta) error {
	return s.svc.store.Secret().Update(ctx, secret, opts)
}

func (s *secretSvc) Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOperateMeta) error {
	return s.svc.store.Secret().Delete(ctx, username, secretID, opts)
}

func (s *secretSvc) DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOperateMeta) error {
	return s.svc.store.Secret().DeleteCollection(ctx, username, secretIDs, opts)
}

func (s *secretSvc) Get(ctx context.Context, username, secretID string, opts metav1.GetOperateMeta) (*v1.Secret, error) {
	return s.svc.store.Secret().Get(ctx, username, secretID, opts)
}

func (s *secretSvc) List(ctx context.Context, username string, opts metav1.ListOperateMeta) (*v1.SecretList, error) {
	return s.svc.store.Secret().List(ctx, username, opts)
}
