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
	"net/url"
	"strings"
)

type User interface {
	Create(ctx context.Context, user *v1.User, opts metaV1.CreateOperateMeta) error
	Update(ctx context.Context, user *v1.User, opts metaV1.UpdateOperateMeta) error
	Delete(ctx context.Context, name string, opts metaV1.DeleteOperateMeta) error
	DeleteCollection(ctx context.Context, names []string, opts metaV1.DeleteOperateMeta) error
	Get(ctx context.Context, name string, opts metaV1.GetOperateMeta) (*v1.User, error)
	List(ctx context.Context, opts metaV1.ListOperateMeta) (*v1.UserList, error)
}

type user struct {
	client client.Client
}

func newUser(client client.Client) User {
	return &user{client: client}
}

// prepare is a template for this page.
func (u *user) prepare() client.Request {
	return u.client.Get().Resource(client.ResUser).Version(client.V1)
}

// handleResErr is a template to return error.
func (u *user) handleResErr(res client.Response) error {
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

func (u *user) Create(ctx context.Context, user *v1.User, opts metaV1.CreateOperateMeta) error {
	res := u.prepare().
		Verb(client.VerbPost).
		Body(user).
		Meta(opts).
		Send(ctx)
	return u.handleResErr(res)
}

func (u *user) Update(ctx context.Context, user *v1.User, opts metaV1.UpdateOperateMeta) (err error) {
	if opts.DryRun {
		return nil
	}
	res := u.prepare().
		Verb(client.VerbPUT).
		Name(user.Name).
		Body(user).
		Meta(opts).
		Send(ctx)
	return u.handleResErr(res)
}

func (u *user) Delete(ctx context.Context, name string, opts metaV1.DeleteOperateMeta) error {
	res := u.prepare().Verb(client.VerbDelete).Meta(opts).Name(name).Send(ctx)
	return u.handleResErr(res)
}

func (u *user) DeleteCollection(ctx context.Context, names []string, opts metaV1.DeleteOperateMeta) error {
	var ps = url.Values{}
	ps.Set("names", strings.Join(names, ","))
	res := u.prepare().Verb(client.VerbDelete).Meta(opts).Params(ps).Send(ctx)
	return u.handleResErr(res)
}

func (u *user) Get(ctx context.Context, name string, opts metaV1.GetOperateMeta) (user *v1.User, err error) {
	res := u.prepare().Verb(client.VerbGET).Meta(opts).Name(name).Send(ctx)
	if err = u.handleResErr(res); err == nil {
		err = res.Into(user)
	}
	return
}

func (u *user) List(ctx context.Context, opts metaV1.ListOperateMeta) (list *v1.UserList, err error) {
	res := u.prepare().Verb(client.VerbGET).Meta(opts).Send(ctx)
	if err = u.handleResErr(res); err == nil {
		err = res.Into(list)
	}
	return
}

var _ User = &user{}
