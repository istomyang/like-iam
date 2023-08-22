package v1

import (
	"context"
	"fmt"
	"github.com/ory/ladon"
	authzV1 "istomyang.github.com/like-iam/api/authzserver/v1"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/client"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/coder"
	"istomyang.github.com/like-iam/iam-sdk-go/service/iam"
)

type Authz interface {
	Authorize(ctx context.Context, request *ladon.Request) *authzV1.Response
}

type authz struct {
	client client.Client
}

func NewAuthz(client client.Client) Authz {
	return &authz{client: client}
}

func (a *authz) Authorize(ctx context.Context, request *ladon.Request) *authzV1.Response {
	res := a.client.Post().Action("authz").Version(client.V1).Body(request).Send(ctx)
	if err := res.Error(); err != nil {
		return a.errRes(err)
	}
	raw, err := res.Raw()
	if err != nil {
		return a.errRes(err)
	}

	if es := web.IsErrResponse(raw, coder.Get(iam.CoderRegisterName)); es != nil {
		return a.errRes(fmt.Errorf(es.String()))
	}

	var r *authzV1.Response
	if err := res.Into(r); err != nil {
		return a.errRes(err)
	}
	return r
}

func (a *authz) errRes(err error) *authzV1.Response {
	return &authzV1.Response{
		Allowed: false,
		Reason:  err.Error(),
		Error:   err.Error(),
	}
}

var _ Authz = &authz{}
