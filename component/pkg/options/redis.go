package options

import (
	"github.com/spf13/pflag"
	"time"
)

// RedisOpts is used for cobra and viper.
// Refer to redis.Options, redis.UniversalOptions, redis.ClusterOptions and redis.FailoverOptions.
type RedisOpts struct {
	Addrs    []string `json:"addrs,omitempty" mapstructure:"addrs"`
	Username string   `json:"username,omitempty" mapstructure:"username"`
	Password string   `json:"password,omitempty" mapstructure:"password"`

	Timeout time.Duration `json:"timeout,omitempty" mapstructure:"timeout"`

	// The sentinel master name.
	// Only failover clients.
	MasterName string `json:"masterName,omitempty" mapstructure:"masterName"`

	// TODO:
	// see: https://pkg.go.dev/crypto/tls#LoadX509KeyPair
	UseTLS bool `json:"useTLS,omitempty" mapstructure:"useTLS"`
	UseSSH bool `json:"useSSH,omitempty" mapstructure:"useSSH"`
}

func NewRedisOpts() *RedisOpts {
	return &RedisOpts{
		Addrs: []string{"127.0.0.1:6379"},
	}
}

func (o *RedisOpts) Validate() []error {
	return nil
}

func (o *RedisOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Addrs, "redis.addrs", o.Addrs, "If length of it greater than 1, use cluster, then use simple.")
	fs.StringVar(&o.Username, "redis.username", o.Username, "Username for access to redis service.")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Optional auth password for Redis db.")
	fs.DurationVar(&o.Timeout, "redis.timeout", o.Timeout, "Timeout for all waiting scenes.")
	fs.StringVar(&o.MasterName, "redis.master-name", o.MasterName, "master name only for failover client.")
	fs.BoolVar(&o.UseTLS, "redis.use-tls", o.UseTLS, "Enable TLS/SSL.")
	fs.BoolVar(&o.UseSSH, "redis.use-ssh", o.UseSSH, "Enable connect over SSH channel.")
}
