package pumper

import (
	"istomyang.github.com/like-iam/component/pkg/app"
	generaloptions "istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"istomyang.github.com/like-iam/log"
	"time"
)

type Options struct {
	interval     time.Duration
	mutexExpiry  time.Duration
	redisOptions *generaloptions.RedisOpts
	Log          *log.Options
	pumps        map[string]any
}

func NewOptions(basename string) *Options {
	return &Options{
		interval:     time.Minute,
		Log:          log.NewOptions(basename, nil),
		mutexExpiry:  time.Minute * 5,
		redisOptions: generaloptions.NewRedisOpts(),
		pumps: map[string]any{"csv": pumps.PumpConfig{
			Name:       "CSVs",
			Filter:     nil,
			OmitDetail: false,
			Timeout:    time.Second * 5,
			Extend: map[string]any{
				"file-dir": "",
			},
		}},
	}
}

func (o *Options) Flags(appFss *app.NamedFlagSets) {
	o.Log.AddFlags(appFss.AddFlagSet("log"))
	o.redisOptions.AddFlags(appFss.AddFlagSet("redis"))
	p := appFss.AddFlagSet("pumper")
	p.DurationVar(&o.mutexExpiry, "mutex-expiry", o.mutexExpiry, "Expiry duration of distributed mutex.")
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.Log.Validate()...)
	return errs
}
