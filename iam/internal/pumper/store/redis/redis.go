package redis

import (
	"context"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/iam/internal/pkg"
	"istomyang.github.com/like-iam/iam/internal/pumper/store"
	"istomyang.github.com/like-iam/log"
	"strings"
)

type redisStore struct {
	client *conn.RedisClient

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisStore(ctx context.Context) store.Factory {
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

func (r *redisStore) Pop(key string) ([][]byte, error) {
	log.Debug("start pop analytics data from pipeline.")

	if strings.TrimSpace(key) == "" {
		key = pkg.AnalyticsKey
	}

	var rg []string
	var err error
	pipeline := r.client.UniversalClient().TxPipeline()
	rg, err = pipeline.LRange(r.ctx, pkg.AnalyticsKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	pipeline.Del(r.ctx, key)
	if _, err = pipeline.Exec(r.ctx); err != nil {
		return nil, err
	} else {
		_ = pipeline.Close()
	}

	result := make([][]byte, len(rg))
	for i, s := range rg {
		result[i] = []byte(s)
	}

	log.Debugf("analytics data's result: %v", result)

	return result, nil
}

var _ store.Factory = &redisStore{}
