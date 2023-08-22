package authorization

import (
	"context"
	"github.com/ory/ladon"
	authzV1 "istomyang.github.com/like-iam/api/authzserver/v1"
	"istomyang.github.com/like-iam/log"
	"sync"
)

type Authorizator struct {
	l *ladon.Ladon

	ctx    context.Context
	cancel context.CancelFunc
}

var (
	au   *Authorizator
	once sync.Once
)

func GetAuthorizator() *Authorizator {
	return au
}

func NewAuthorizator(ctx context.Context) (*Authorizator, error) {
	var err error

	once.Do(func() {
		au = &Authorizator{}
		au.ctx, au.cancel = context.WithCancel(ctx)

		l := ladon.Ladon{}
		l.Manager, err = newManager(au.ctx)
		l.AuditLogger, err = newAuditor()
		l.Metric, err = newMetric()
	})

	return au, err
}

func (authz *Authorizator) Authorize(request *ladon.Request) *authzV1.Response {
	log.Debugf("authorize request: %v", request)

	if err := authz.l.IsAllowed(request); err != nil {
		return &authzV1.Response{
			Allowed: false,
			Reason:  err.Error(),
		}
	}

	return &authzV1.Response{
		Allowed: true,
	}
}
