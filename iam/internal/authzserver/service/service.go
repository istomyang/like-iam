package service

import (
	"context"
	"fmt"
	"github.com/ory/ladon"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service/cache"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service/subscribe"
	"sync"
)

type Service interface {
	Find(request *ladon.Request) (ladon.Policies, error)

	FindSecret(kid string) (*pb.SecretInfo, error)

	Run() error
	Close() error
}

type service struct {
	cache Cache
	sub   Subscribe

	sync chan bool

	ctx    context.Context
	cancel context.CancelFunc
}

func GetService() Service {
	return svr
}

var (
	svr  Service
	once sync.Once
)

func NewService(ctx context.Context) (Service, error) {
	var err error

	once.Do(func() {
		s := &service{}
		s.ctx, s.cancel = context.WithCancel(ctx)

		if s.cache, err = cache.NewMemory(ctx); err != nil {
			return
		}

		if s.sub, err = subscribe.NewRedisSubClient(ctx); err != nil {
			return
		}

		s.sync = make(chan bool, 1)

		if err := s.sub.OnReceive(s.sync); err != nil {
			return
		}

		svr = s
	})

	return svr, err
}

func (s *service) Find(request *ladon.Request) (ladon.Policies, error) {
	username, ok := request.Context["username"]
	if !ok {
		return nil, fmt.Errorf("username not in request %v", request.Context)
	}

	return s.cache.GetPolicy(username.(string))
}

func (s *service) FindSecret(kid string) (*pb.SecretInfo, error) {
	return s.cache.GetSecret(kid)
}

func (s *service) Run() error {

	go func() {
		_ = s.sub.Run()
	}()

	go func() {
		_ = s.cache.Run()
	}()

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-s.sync:
				_ = s.cache.Run()
			}
		}
	}()

	return nil
}

func (s *service) Close() error {
	s.cancel()

	if err := s.cache.Close(); err != nil {
		return err
	}

	if err := s.sub.Close(); err != nil {
		return err
	}

	return nil
}

var _ Service = &service{}
