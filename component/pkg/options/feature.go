package options

import "github.com/spf13/pflag"

type FeatureOptions struct {
	// https://pkg.go.dev/github.com/zsais/go-gin-prometheus
	Metrics bool `json:"metrics" mapstructure:"metrics"`

	// https://pkg.go.dev/github.com/gin-contrib/pprof
	Profile bool `json:"profile" mapstructure:"profile"`
}

func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{
		Metrics: true,
		Profile: true,
	}
}

func (o *FeatureOptions) Validate() []error {
	var err []error
	return err
}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Metrics, "server.metrics", o.Metrics, ""+
		"Enable Metrics, see more: https://pkg.go.dev/github.com/zsais/go-gin-prometheus .")

	fs.BoolVar(&o.Profile, "server.profile", o.Profile, ""+
		"Enable Profile, see more: https://pkg.go.dev/github.com/gin-contrib/pprof .")
}
