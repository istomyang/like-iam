package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

type InsecureServerOpts struct {
	Address string `json:"address" mapstructure:"address"`
	Port    int    `json:"port"    mapstructure:"port"`
}

func NewInsecureServerOpts() *InsecureServerOpts {
	return &InsecureServerOpts{
		Address: "127.0.0.1",
		Port:    8080,
	}
}

func (o *InsecureServerOpts) Validate() []error {
	var err []error

	if !(o.Port > 1023 && o.Port <= 65535) {
		err = append(err, fmt.Errorf("insecure-server port %d must between 1024 to 65535", o.Port))
	}

	if net.ParseIP(o.Address) == nil {
		err = append(err, fmt.Errorf("insecure-server address %s is not valid", o.Address))
	}

	return err
}

func (o *InsecureServerOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "insecure.address", o.Address, ""+
		"The IP address on which to serve the --insecure.port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	fs.IntVar(&o.Port, "insecure.port", o.Port, ""+
		"The port on which to serve unsecured, unauthenticated access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxy to this "+
		"port. This is performed by nginx in the default setup.")
}
