package service

import (
	"context"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
)

type PolicySvc interface {
	Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOperateMeta) error
	Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOperateMeta) error
	Delete(ctx context.Context, username string, name string, opts metav1.DeleteOperateMeta) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOperateMeta) error
	Get(ctx context.Context, username string, name string, opts metav1.GetOperateMeta) (*v1.Policy, error)
	List(ctx context.Context, username string, opts metav1.ListOperateMeta) (*v1.PolicyList, error)
}

type policySvc struct {
	svc *service
}

func newPolicySvc(svc *service) PolicySvc {
	return &policySvc{svc: svc}
}

func (p *policySvc) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOperateMeta) error {
	return p.svc.store.Policy().Create(ctx, policy, opts)
}

func (p *policySvc) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOperateMeta) error {
	return p.svc.store.Policy().Update(ctx, policy, opts)
}

func (p *policySvc) Delete(ctx context.Context, username string, name string, opts metav1.DeleteOperateMeta) error {
	return p.svc.store.Policy().Delete(ctx, username, name, opts)
}

func (p *policySvc) DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOperateMeta) error {
	return p.svc.store.Policy().DeleteCollection(ctx, username, names, opts)
}

func (p *policySvc) Get(ctx context.Context, username string, name string, opts metav1.GetOperateMeta) (*v1.Policy, error) {
	return p.svc.store.Policy().Get(ctx, username, name, opts)
}

func (p *policySvc) List(ctx context.Context, username string, opts metav1.ListOperateMeta) (*v1.PolicyList, error) {
	return p.svc.store.Policy().List(ctx, username, opts)
}
