package options

import (
	"github.com/spf13/pflag"
)

// RedisOpts provides options refer to redis.Options.
type RedisOpts struct {
	Host     string   `json:"host" mapstructure:"host"`
	Port     int      `json:"port" mapstructure:"port"`
	Addrs    []string `json:"addrs" mapstructure:"addrs"`
	Username string   `json:"username" mapstructure:"username"`
	Password string   `json:"password" mapstructure:"password"`
	DB       int      `json:"db" mapstructure:"db"`
	// TODO:
}

func NewRedisOpts() *RedisOpts {
	return &RedisOpts{
		Host:     "127.0.0.1",
		Port:     6379,
		Addrs:    nil,
		Username: "",
		Password: "",
		DB:       0,
	}
}

func (o *RedisOpts) Validate() []error {
	var err []error
	return err
}

func (o *RedisOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "redis.host", o.Host, "Hostname of your Redis server.")
	fs.IntVar(&o.Port, "redis.port", o.Port, "The port the Redis server is listening on.")
	fs.StringSliceVar(&o.Addrs, "redis.addrs", o.Addrs, "A set of redis address(format: 127.0.0.1:6379).")
	fs.StringVar(&o.Username, "redis.username", o.Username, "Username for access to redis service.")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Optional auth password for Redis db.")
	fs.IntVar(&o.DB, "redis.database", o.DB, ""+
		"By default, the database is 0. Setting the database is not supported with redis cluster. "+
		"As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.")
}
