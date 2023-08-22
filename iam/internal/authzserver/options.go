package authzserver

import (
	"istomyang.github.com/like-iam/component/pkg/app"
	generaloptions "istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/authzserver/analytics"
	"istomyang.github.com/like-iam/log"
)

type Options struct {
	httpSvrOptions     *generaloptions.ServerOpts
	insecureSvrOptions *generaloptions.InsecureServerOpts
	secureSvrOptions   *generaloptions.SecureServerOpts
	featureOptions     *generaloptions.FeatureOptions
	redisOptions       *generaloptions.RedisOpts
	gRPCOptions        *generaloptions.GRPCOpts
	analyticsOptions   *analytics.Options
	Log                *log.Options
	clientCA           string
}

func NewOptions(basename string) *Options {
	return &Options{
		httpSvrOptions:     generaloptions.NewServerOpts(),
		insecureSvrOptions: generaloptions.NewInsecureServerOpts(),
		secureSvrOptions:   generaloptions.NewSecureServerOpts(),
		featureOptions:     generaloptions.NewFeatureOptions(),
		redisOptions:       generaloptions.NewRedisOpts(),
		gRPCOptions:        generaloptions.NewGRPCOpts(),
		Log:                log.NewOptions(basename, nil),
		analyticsOptions:   analytics.NewAnalyticsOptions(),
	}
}

func (o *Options) Flags(appFss *app.NamedFlagSets) {
	o.httpSvrOptions.AddFlags(appFss.AddFlagSet("general"))
	o.insecureSvrOptions.AddFlags(appFss.AddFlagSet("insecure server"))
	o.secureSvrOptions.AddFlags(appFss.AddFlagSet("secure server"))
	o.featureOptions.AddFlags(appFss.AddFlagSet("feature"))
	o.redisOptions.AddFlags(appFss.AddFlagSet("redis"))
	o.gRPCOptions.AddFlags(appFss.AddFlagSet("gRPC"))
	o.Log.AddFlags(appFss.AddFlagSet("log"))

	o.analyticsOptions.AddFlags(appFss.AddFlagSet("store"))
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.httpSvrOptions.Validate()...)
	errs = append(errs, o.insecureSvrOptions.Validate()...)
	errs = append(errs, o.secureSvrOptions.Validate()...)
	errs = append(errs, o.redisOptions.Validate()...)
	errs = append(errs, o.gRPCOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.featureOptions.Validate()...)
	errs = append(errs, o.analyticsOptions.Validate()...)
	return errs
}
