package conn

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/log"
	"sync"
)

type ClientMode string

const (
	ClientCluster  ClientMode = "Cluster"
	ClientSentinel ClientMode = "Sentinel"
	ClientSingle   ClientMode = "Single"
)

var (
	singletonRedis *RedisClient
	onceRedis      sync.Once
)

// RedisClient holds redis's different clients.
type RedisClient struct {
	mode      ClientMode
	universal redis.UniversalClient
}

// GetRedisClient is a correct-semantic func to get a redis client.
// Generally, app manages lifecycle of client and subsystem just gets client, not run or close it.
func GetRedisClient() *RedisClient {
	return NewRedisClientOr(nil)
}

// NewRedisClientOr return a universal redis client provided by go-redis.
// Call this more times will return singleton, and nil opts return singleton.
// Refer to redis.Options, redis.UniversalOptions, redis.ClusterOptions and redis.FailoverOptions.
func NewRedisClientOr(opts *options.RedisOpts) *RedisClient {

	onceRedis.Do(func() {
		c := &RedisClient{}

		// Use redis.NewUniversalClient() logic.
		if opts.MasterName != "" {
			log.Info("--> [REDIS] Creating sentinel-backed failover client")
			c.mode = ClientSentinel
		} else if len(opts.Addrs) > 1 {
			log.Info("--> [REDIS] Creating cluster client")
			c.mode = ClientCluster
		} else {
			log.Info("--> [REDIS] Creating single-node client")
			c.mode = ClientSingle
		}

		c.universal = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:              opts.Addrs,
			DB:                 0,
			Dialer:             nil,
			OnConnect:          nil,
			Username:           opts.Username,
			Password:           opts.Password,
			SentinelUsername:   "",
			SentinelPassword:   "",
			MaxRetries:         0,
			MinRetryBackoff:    0,
			MaxRetryBackoff:    0,
			DialTimeout:        opts.Timeout,
			ReadTimeout:        opts.Timeout,
			WriteTimeout:       opts.Timeout,
			PoolFIFO:           false,
			PoolSize:           0,
			MinIdleConns:       0,
			MaxConnAge:         0,
			PoolTimeout:        0,
			IdleTimeout:        0,
			IdleCheckFrequency: 0,
			TLSConfig:          nil,
			MaxRedirects:       0,
			ReadOnly:           false,
			RouteByLatency:     false,
			RouteRandomly:      false,
			MasterName:         opts.MasterName,
		})

		singletonRedis = c
	})

	return singletonRedis
}

func (c *RedisClient) Run() error {
	return nil
}

func (c *RedisClient) Mode() ClientMode {
	return c.mode
}

func (c *RedisClient) UniversalClient() redis.UniversalClient {
	return c.universal
}

func (c *RedisClient) Close() error {
	return c.UniversalClient().Close()
}

func (c *RedisClient) ClusterPing() error {
	if c.mode != ClientCluster {
		return fmt.Errorf("error mode: %s, check client's mode before run", c.mode)
	}
	return c.UniversalClient().(*redis.ClusterClient).ForEachShard(context.TODO(), func(ctx context.Context, client *redis.Client) error {
		return client.Ping(ctx).Err()
	})
}
