package store

import (
	"context"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/iam/internal/authzserver/analytics"
	"time"
)

type redisStore struct {
	client *conn.RedisClient

	expire time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisStore(ctx context.Context) analytics.Store {
	s := redisStore{}

	s.ctx, s.cancel = context.WithCancel(ctx)

	s.client = conn.GetRedisClient()

	return &s
}

func (r *redisStore) Run() error {
	return nil
}

func (r *redisStore) Close() error {
	r.cancel()
	return nil
}

func (r *redisStore) WithExpire(duration time.Duration) {
	r.expire = duration
}

func (r *redisStore) Send(k string, batch [][]byte) error {
	pl := r.client.UniversalClient().Pipeline()
	pl.Expire(r.ctx, k, r.expire)
	for _, bytes := range batch {
		pl.RPush(r.ctx, k, bytes)
	}
	if _, err := pl.Exec(r.ctx); err != nil {
		return err
	}
	return nil
}

var _ analytics.Store = &redisStore{}
