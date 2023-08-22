package v1

import (
	"context"
	"fmt"
	v1 "istomyang.github.com/like-iam/api/apiserver/v1"
	metaV1 "istomyang.github.com/like-iam/component-base/meta/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/client"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/coder"
	"istomyang.github.com/like-iam/iam-sdk-go/service/iam"
)

type Policy interface {
	Create(ctx context.Context, policy *v1.Policy, opts metaV1.CreateOperateMeta) error
	Update(ctx context.Context, policy *v1.Policy, opts metaV1.UpdateOperateMeta) error
	Delete(ctx context.Context, name string, opts metaV1.DeleteOperateMeta) error
	DeleteCollection(ctx context.Context, opts metaV1.DeleteOperateMeta, listOpts metaV1.ListOperateMeta) error
	Get(ctx context.Context, name string, opts metaV1.GetOperateMeta) (*v1.Policy, error)
	List(ctx context.Context, opts metaV1.ListOperateMeta) (*v1.PolicyList, error)
}

type policy struct {
	client client.Client
}

func newPolicy(client client.Client) Policy {
	return &policy{client: client}
}

// prepare is a template for this page.
func (p *policy) prepare() client.Request {
	return p.client.Get().Resource(client.ResPolicy).Version(client.V1)
}

// handleResErr is a template to return error.
func (p *policy) handleResErr(res client.Response) error {
	if err := res.Error(); err != nil {
		return err
	}
	raw, err := res.Raw()
	if err != nil {
		return err
	}
	if es := web.IsErrResponse(raw, coder.Get(iam.CoderRegisterName)); es != nil {
		return fmt.Errorf(es.String())
	}
	return nil
}

func (p *policy) Create(ctx context.Context, policy *v1.Policy, opts metaV1.CreateOperateMeta) error {
	res := p.prepare().Verb(client.VerbPost).Meta(opts).Body(policy).Send(ctx)
	return p.handleResErr(res)
}

func (p *policy) Update(ctx context.Context, policy *v1.Policy, opts metaV1.UpdateOperateMeta) error {
	res := p.prepare().Verb(client.VerbPUT).Meta(opts).Body(policy).Send(ctx)
	return p.handleResErr(res)
}

func (p *policy) Delete(ctx context.Context, name string, opts metaV1.DeleteOperateMeta) error {
	res := p.prepare().Verb(client.VerbDelete).Meta(opts).Name(name).Send(ctx)
	return p.handleResErr(res)
}

func (p *policy) DeleteCollection(ctx context.Context, opts metaV1.DeleteOperateMeta, listOpts metaV1.ListOperateMeta) error {
	res := p.prepare().Verb(client.VerbDelete).Meta(opts).Meta(listOpts).Send(ctx)
	return p.handleResErr(res)
}

func (p *policy) Get(ctx context.Context, name string, opts metaV1.GetOperateMeta) (po *v1.Policy, err error) {
	res := p.prepare().Verb(client.VerbGET).Meta(opts).Name(name).Send(ctx)
	if err = p.handleResErr(res); err == nil {
		err = res.Into(po)
	}
	return
}

func (p *policy) List(ctx context.Context, opts metaV1.ListOperateMeta) (lst *v1.PolicyList, err error) {
	res := p.prepare().Verb(client.VerbGET).Meta(opts).Send(ctx)
	if err = p.handleResErr(res); err == nil {
		err = res.Into(lst)
	}
	return
}

var _ Policy = &policy{}
