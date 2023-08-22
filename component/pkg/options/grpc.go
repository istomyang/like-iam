package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

type GRPCOpts struct {
	Address    string `json:"address,omitempty" mapstructure:"address"`
	Port       int    `json:"port,omitempty" mapstructure:"port"`
	MaxMsgSize int    `json:"max_msg_size,omitempty" mapstructure:"max_msg_size"`
}

func NewGRPCOpts() *GRPCOpts {
	return &GRPCOpts{
		Address:    "0.0.0.0",
		Port:       8081,
		MaxMsgSize: 1024 * 1024 * 4,
	}
}

func (o *GRPCOpts) Validate() []error {
	var err []error
	if !(o.Port > 1023 && o.Port <= 65535) {
		err = append(err, fmt.Errorf("secure-server port %d must between 1024 to 65535", o.Port))
	}
	return err
}

func (o *GRPCOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "apiserver.bind-address", o.Address, ""+
		"The IP address on which to serve the --apiserver.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&o.Port, "apiserver.bind-port", o.Port, ""+
		"The port on which to serve unsecured, unauthenticated apiserver access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")

	fs.IntVar(&o.MaxMsgSize, "apiserver.max-msg-size", o.MaxMsgSize, "gRPC max message size.")
}
