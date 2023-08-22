package apiserver

import (
	"istomyang.github.com/like-iam/component/pkg/app"
	generaloptions "istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/log"
)

type Options struct {
	httpSvrOptions     *generaloptions.ServerOpts
	insecureSvrOptions *generaloptions.InsecureServerOpts
	secureSvrOptions   *generaloptions.SecureServerOpts
	mysqlOptions       *generaloptions.MySQLOpts
	redisOptions       *generaloptions.RedisOpts
	jwtOptions         *generaloptions.JwtOpts
	gRPCOptions        *generaloptions.GRPCOpts
	featureOptions     *generaloptions.FeatureOptions

	Log *log.Options
}

func NewOptions(basename string) *Options {
	return &Options{
		httpSvrOptions:     generaloptions.NewServerOpts(),
		insecureSvrOptions: generaloptions.NewInsecureServerOpts(),
		secureSvrOptions:   generaloptions.NewSecureServerOpts(),
		mysqlOptions:       generaloptions.NewMySQLOpts(),
		redisOptions:       generaloptions.NewRedisOpts(),
		jwtOptions:         generaloptions.NewJwtOpts(),
		gRPCOptions:        generaloptions.NewGRPCOpts(),
		featureOptions:     generaloptions.NewFeatureOptions(),
		Log:                log.NewOptions(basename, nil),
	}
}

func (o *Options) Flags(appFss *app.NamedFlagSets) {
	o.httpSvrOptions.AddFlags(appFss.AddFlagSet("general"))
	o.insecureSvrOptions.AddFlags(appFss.AddFlagSet("insecure server"))
	o.secureSvrOptions.AddFlags(appFss.AddFlagSet("secure server"))
	o.mysqlOptions.AddFlags(appFss.AddFlagSet("mysql"))
	o.redisOptions.AddFlags(appFss.AddFlagSet("redis"))
	o.jwtOptions.AddFlags(appFss.AddFlagSet("jwt"))
	o.gRPCOptions.AddFlags(appFss.AddFlagSet("gRPC"))
	o.featureOptions.AddFlags(appFss.AddFlagSet("feature"))
	o.Log.AddFlags(appFss.AddFlagSet("log"))
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.httpSvrOptions.Validate()...)
	errs = append(errs, o.insecureSvrOptions.Validate()...)
	errs = append(errs, o.secureSvrOptions.Validate()...)
	errs = append(errs, o.mysqlOptions.Validate()...)
	errs = append(errs, o.redisOptions.Validate()...)
	errs = append(errs, o.jwtOptions.Validate()...)
	errs = append(errs, o.gRPCOptions.Validate()...)
	errs = append(errs, o.featureOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	return errs
}
