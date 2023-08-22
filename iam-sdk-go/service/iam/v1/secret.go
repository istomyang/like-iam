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

type Secret interface {
	Create(ctx context.Context, secret *v1.Secret, opts metaV1.CreateOperateMeta) error
	Update(ctx context.Context, secret *v1.Secret, opts metaV1.UpdateOperateMeta) error
	Delete(ctx context.Context, name string, opts metaV1.DeleteOperateMeta) error
	DeleteCollection(ctx context.Context, opts metaV1.DeleteOperateMeta, listOpts metaV1.ListOperateMeta) error
	Get(ctx context.Context, name string, opts metaV1.GetOperateMeta) (*v1.Secret, error)
	List(ctx context.Context, opts metaV1.ListOperateMeta) (*v1.SecretList, error)
}

type secret struct {
	client client.Client
}

func newSecret(client client.Client) Secret {
	return &secret{client: client}
}

// prepare is a template for this page.
func (s *secret) prepare() client.Request {
	return s.client.Get().Resource(client.ResSecret).Version(client.V1)
}

// handleResErr is a template to return error.
func (s *secret) handleResErr(res client.Response) error {
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

func (s *secret) Create(ctx context.Context, secret *v1.Secret, opts metaV1.CreateOperateMeta) error {
	res := s.prepare().Verb(client.VerbPost).Meta(opts).Body(secret).Send(ctx)
	return s.handleResErr(res)
}

func (s *secret) Update(ctx context.Context, secret *v1.Secret, opts metaV1.UpdateOperateMeta) error {
	res := s.prepare().Verb(client.VerbPUT).Meta(opts).Body(secret).Send(ctx)
	return s.handleResErr(res)
}

func (s *secret) Delete(ctx context.Context, name string, opts metaV1.DeleteOperateMeta) error {
	res := s.prepare().Verb(client.VerbDelete).Meta(opts).Name(name).Send(ctx)
	return s.handleResErr(res)
}

func (s *secret) DeleteCollection(ctx context.Context, opts metaV1.DeleteOperateMeta, listOpts metaV1.ListOperateMeta) error {
	res := s.prepare().Verb(client.VerbDelete).Meta(opts).Meta(listOpts).Send(ctx)
	return s.handleResErr(res)
}

func (s *secret) Get(ctx context.Context, name string, opts metaV1.GetOperateMeta) (sec *v1.Secret, err error) {
	res := s.prepare().Verb(client.VerbGET).Meta(opts).Name(name).Send(ctx)
	if err = s.handleResErr(res); err == nil {
		err = res.Into(sec)
	}
	return
}

func (s *secret) List(ctx context.Context, opts metaV1.ListOperateMeta) (list *v1.SecretList, err error) {
	res := s.prepare().Verb(client.VerbGET).Meta(opts).Send(ctx)
	if err = s.handleResErr(res); err == nil {
		err = res.Into(list)
	}
	return
}

var _ Secret = &secret{}
