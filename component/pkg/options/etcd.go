package options

import (
	"github.com/spf13/pflag"
)

type EtcdOpts struct {
}

func NewEtcdOpts() *EtcdOpts {
	return &EtcdOpts{}
}

func (o *EtcdOpts) Validate() []error {
	var err []error

	return err
}

func (o *EtcdOpts) AddFlags(fs *pflag.FlagSet) {
}
