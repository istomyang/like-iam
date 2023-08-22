package authorization

import (
	"context"
	"github.com/ory/ladon"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service"
)

type manager struct {
	svc service.Service

	ctx    context.Context
	cancel context.CancelFunc
}

var _ ladon.Manager = &manager{}

func newManager(ctx context.Context) (ladon.Manager, error) {
	m := &manager{}

	m.ctx, m.cancel = context.WithCancel(ctx)
	m.svc = service.GetService()

	return m, nil
}

func (m *manager) Create(policy ladon.Policy) error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Update(policy ladon.Policy) error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Get(id string) (ladon.Policy, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) GetAll(limit, offset int64) (ladon.Policies, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	return m.svc.Find(r)
}

func (m *manager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	//TODO implement me
	panic("implement me")
}
