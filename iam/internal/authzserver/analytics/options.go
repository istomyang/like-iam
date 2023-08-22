package analytics

import (
	"github.com/spf13/pflag"
	"runtime"
	"time"
)

type Options struct {
	Enabled bool `json:"enabled,omitempty" mapstructure:"enabled"`

	Workers       int           `json:"workers,omitempty" mapstructure:"workers"`
	BatchSize     uint64        `json:"batch-size,omitempty" mapstructure:"batchSize"`
	FlushInterval time.Duration `json:"flush-interval,omitempty" mapstructure:"flushInterval"`
	Expire        time.Duration `json:"expire" mapstructure:"expire"`
}

func NewAnalyticsOptions() *Options {
	return &Options{
		Enabled:       true,
		Workers:       runtime.NumCPU(),
		BatchSize:     10,
		FlushInterval: time.Second * 3,
		Expire:        time.Hour * 24 * 365,
	}
}

func (o *Options) Validate() []error {
	return nil
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Enabled, "enabled", o.Enabled, "Enable store.")
	fs.IntVar(&o.Workers, "workers", o.Workers, "The number of workers to handle recording info concurrently.")
	fs.Uint64Var(&o.BatchSize, "batch-size", o.BatchSize, "Size of recording infos to be sent to storage once.")
	fs.DurationVar(&o.FlushInterval, "flush-interval", o.FlushInterval, "Time duration to send to storage.")
	fs.DurationVar(&o.Expire, "expire", o.Expire, "Expiration of every recording info.")
}
