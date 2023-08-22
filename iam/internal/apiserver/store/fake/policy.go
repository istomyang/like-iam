package fake

import (
	"context"
	"gorm.io/gorm"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	metav1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/util/idutil"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"istomyang.github.com/like-iam/iam/internal/pkg/codes"
	"strings"
	"time"
)

type policy struct {
	db *datastore
}

func newPolicy(ds *datastore) store.PolicyStore {
	return &policy{db: ds}
}

func (p *policy) Create(c context.Context, policy *v1.Policy, opts metav1.CreateOperateMeta) error {
	p.db.Lock()
	defer p.db.Unlock()

	for _, v := range p.db.policies {
		if v.Username == policy.Username && v.InstanceID == policy.InstanceID && v.Name == policy.Name {
			return errors.WithCode(codes.ErrPolicyAlreadyExit, "policy has already exits.")
		}
	}

	policy.ID = uint64(len(p.db.policies))
	policy.InstanceID, _ = idutil.GetInstanceId(policy.ID, "policy", 6)

	p.db.policies = append(p.db.policies, policy)

	return nil
}

func (p *policy) Update(c context.Context, policy *v1.Policy, opts metav1.UpdateOperateMeta) error {
	p.db.Lock()
	defer p.db.Unlock()

	for i, v := range p.db.policies {
		if v.Username == policy.Username && v.InstanceID == policy.InstanceID {
			p.db.policies[i] = policy
			return nil
		}
	}

	return errors.WithCode(codes.ErrPolicyNotFound, "secret not found.")
}

func (p *policy) Delete(c context.Context, username string, name string, opts metav1.DeleteOperateMeta) error {
	p.db.Lock()
	defer p.db.Unlock()

	for _, v := range p.db.policies {
		if v.Username == username && v.Name == name {
			v.DeletedAt = gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			}
			return nil
		}
	}

	return errors.WithCode(codes.ErrSecretNotFound, "secret not found.")
}

func (p *policy) DeleteCollection(c context.Context, username string, names []string, opts metav1.DeleteOperateMeta) error {
	p.db.Lock()
	defer p.db.Unlock()

	ss := p.db.policies
	for _, v := range ss {
		if v.Username == username {
			for _, d := range names {
				if v.Name == d {
					v.DeletedAt = gorm.DeletedAt{
						Time:  time.Now(),
						Valid: true,
					}
				}
			}
		}
	}

	return nil
}

func (p *policy) Get(c context.Context, username string, name string, opts metav1.GetOperateMeta) (*v1.Policy, error) {
	p.db.Lock()
	defer p.db.Unlock()

	for _, v := range p.db.policies {
		if v.Username == username && v.Name == name {
			return v, nil
		}
	}

	return nil, errors.WithCode(codes.ErrSecretNotFound, "policy name `%s` in user `%s` not found", name, username)
}

func (p *policy) List(c context.Context, username string, opts metav1.ListOperateMeta) (*v1.PolicyList, error) {
	p.db.Lock()
	defer p.db.Unlock()

	var r []*v1.Policy
	for _, v := range p.db.policies {
		if strings.Contains(v.Username, username) {
			r = append(r, v)
		}
	}

	return &v1.PolicyList{
		ListMeta: metav1.ListMeta{TotalCount: int64(len(r))},
		Items:    r,
	}, nil
}

func (p *policy) ClearOutdated(c context.Context, maxReserveDays int) (int64, error) {
	ps := make([]*v1.Policy, 0)
	for _, v := range p.db.policies {
		if v.DeletedAt.Valid && time.Now().After(v.DeletedAt.Time) {
			ps = append(ps, v)
		}
	}
	p.db.policies = ps
	return int64(len(ps)), nil
}
