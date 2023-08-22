package options

import (
	"github.com/spf13/pflag"
	"regexp"
)

type PrometheusOpts struct {
	Addr string
	Path string
}

func (p *PrometheusOpts) AddFlags(set *pflag.FlagSet) {
	//TODO implement me
	panic("implement me")
}

func (p *PrometheusOpts) Validate() []error {
	var errs []error
	re := regexp.MustCompile(`(.*):(.*)`)
	p.Addr = re.FindString(p.Addr)
	if p.Path == "" {
		p.Path = "/metrics"
	}
	return errs
}

var _ ValidatableOptions = &PrometheusOpts{}

var _ FlagOptions = &PrometheusOpts{}
