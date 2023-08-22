package options

import (
	"github.com/spf13/pflag"
)

type InfluxOpts struct {
	Addr     string `json:"addr,omitempty" mapstructure:"addr"`
	Username string `json:"username,omitempty" mapstructure:"username"`
	Password string `json:"password,omitempty" mapstructure:"password"`

	DatabaseName string   `json:"database-name,omitempty" mapstructure:"database-name"`
	Fields       []string `json:"fields,omitempty" mapstructure:"fields"`
	Tags         []string `json:"tags,omitempty" mapstructure:"tags"`
}

func (i *InfluxOpts) AddFlags(set *pflag.FlagSet) {
	//TODO implement me
	panic("implement me")
}

func (i *InfluxOpts) Validate() []error {
	//TODO implement me
	panic("implement me")
}

var _ ValidatableOptions = &InfluxOpts{}

var _ FlagOptions = &InfluxOpts{}
