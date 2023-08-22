package apiserver

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go/v4"
	"github.com/ory/ladon"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/iam/internal/authzserver/store"
	"istomyang.github.com/like-iam/log"
)

type policy struct {
	pb  pb.CacheClient
	ctx context.Context
}

func newPolicy(ctx context.Context, pb pb.CacheClient) store.PolicyStore {
	return &policy{pb: pb, ctx: ctx}
}

func (p *policy) List() (map[string][]*ladon.DefaultPolicy, error) {
	log.Info("loading list policies.")

	req := pb.ListRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1), // cancel offset condition with -1
	}

	var policies *pb.ListPoliciesReply
	var err error

	err = retry.Do(func() error {
		policies, err = p.pb.ListPolicies(p.ctx, &req)
		return err
	}, retry.Attempts(3))
	if err != nil {
		return nil, errors.Wrap(err, "list policies coming from apiserver failed after 3 times.")
	}

	log.Infof("policies loaded count: &d", policies.Count)

	r := make(map[string][]*ladon.DefaultPolicy)

	for _, item := range policies.Items {
		log.Infof("get policies: %s:%s", item.Username, item.Name)

		l := ladon.DefaultPolicy{}

		if err := l.UnmarshalJSON([]byte(item.PolicyShadow)); err != nil {
			return nil, err
		}

		r[item.Username] = append(r[item.Username], &l)
	}

	return r, nil
}
