package cache

import (
	"context"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/ory/ladon"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service"
	"istomyang.github.com/like-iam/iam/internal/authzserver/store"
	"sync"
)

const ErrNotFoundTpl = "k %s not found in %s cache"

type memory struct {
	policy *ristretto.Cache
	secret *ristretto.Cache

	ctx context.Context
	l   *sync.RWMutex
}

func NewMemory(ctx context.Context) (service.Cache, error) {

	config := &ristretto.Config{
		NumCounters:        1e7,     // number of keys to track frequency of (10M).
		MaxCost:            1 << 30, // maximum cost of cache (1GB).
		BufferItems:        64,      // number of keys per Get buffer.
		Metrics:            false,
		OnEvict:            nil,
		OnReject:           nil,
		OnExit:             nil,
		KeyToHash:          nil,
		Cost:               nil,
		IgnoreInternalCost: false,
	}

	m := &memory{ctx: ctx}

	var err error

	if m.secret, err = ristretto.NewCache(config); err != nil {
		return nil, err
	}
	if m.policy, err = ristretto.NewCache(config); err != nil {
		return nil, err
	}

	return m, err
}

func (m *memory) GetPolicy(k string) ([]ladon.Policy, error) {
	m.l.Lock()
	defer m.l.Unlock()

	v, ok := m.policy.Get(k)
	if !ok {
		return nil, fmt.Errorf(ErrNotFoundTpl, k, "policy")
	}
	return v.([]ladon.Policy), nil
}

func (m *memory) GetSecret(k string) (*pb.SecretInfo, error) {
	m.l.Lock()
	defer m.l.Unlock()

	v, ok := m.secret.Get(k)
	if !ok {
		return nil, fmt.Errorf(ErrNotFoundTpl, k, "secret")
	}
	return v.(*pb.SecretInfo), nil
}

func (m *memory) Sync() error {
	m.l.Lock()
	defer m.l.Unlock()

	if err := m.Clear(); err != nil {
		return err
	}

	policies, err := store.Client().Policies().List()
	if err != nil {
		return errors.Wrapf(err, "memory cache sync policy fail")
	}
	for k, policy := range policies {
		m.policy.Set(k, policy, 1)
	}

	secrets, err := store.Client().Secrets().List()
	if err != nil {
		return errors.Wrapf(err, "memory cache sync secret fail")
	}
	for k, secret := range secrets {
		m.secret.Set(k, secret, 1)
	}

	return nil
}

func (m *memory) Clear() error {
	m.l.Lock()
	defer m.l.Unlock()

	m.secret.Clear()
	m.policy.Clear()
	return nil
}

func (m *memory) Run() error {
	go func() {
		err := m.Sync()
		if err != nil {
			// TODO: research it, upload error to top and close to bottom.
			panic(err.Error())
		}
	}()
	return nil
}

func (m *memory) Close() error {
	m.secret.Close()
	m.policy.Close()
	return nil
}

var _ service.Cache = &memory{}
