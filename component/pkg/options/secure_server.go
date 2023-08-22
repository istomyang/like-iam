package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

type SecureServerOpts struct {
	Address string `json:"address" mapstructure:"address"`
	Port    int    `json:"port"    mapstructure:"port"`

	Enable bool       `json:"enable"    mapstructure:"enable"`
	Tls    *TlsConfig `json:"tls"    mapstructure:"tls"`
}

// TlsConfig provides tls key and cert file.
// https://colobu.com/2016/06/07/simple-golang-tls-examples/
type TlsConfig struct {
	KeyFile  string `json:"key-file" mapstructure:"key-file"`
	CertFile string `json:"cert-file" mapstructure:"cert-file"`

	// CertDirectory is a directory which cert files put in, E.g. /var/run/iam
	// If KeyFile and CertFile do not set.
	CertDirectory string `json:"cert-directory"    mapstructure:"cert-directory"`
	// PairName is filename of cert files. E.g. pairName.key and pairName.crt
	PairName string `json:"pair-name"         mapstructure:"pair-name"`
}

func NewSecureServerOpts() *SecureServerOpts {
	return &SecureServerOpts{
		Address: "0.0.0.0",
		Port:    8443,
		Enable:  true,
		Tls: &TlsConfig{
			KeyFile:       "",
			CertFile:      "",
			CertDirectory: "/var/run/iam",
			PairName:      "iam",
		},
	}
}

func (o *SecureServerOpts) Validate() []error {
	var err []error

	if !(o.Port > 1023 && o.Port <= 65535) {
		err = append(err, fmt.Errorf("secure-server port %d must between 1024 to 65535", o.Port))
	}

	if net.ParseIP(o.Address) == nil {
		err = append(err, fmt.Errorf("secure-server address %s is not valid", o.Address))
	}

	return err
}

func (o *SecureServerOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "secure.address", o.Address, ""+
		"The IP address on which to listen for the --secure.port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.BoolVar(&o.Enable, "secure.enable", o.Enable, "The switcher to serve HTTPS with authentication and authorize.")

	fs.IntVar(&o.Port, "secure.port", o.Port, "The port on which to serve HTTPS with authentication and authorize.")

	fs.StringVar(&o.Tls.CertDirectory, "secure.tls.cert-dir", o.Tls.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&o.Tls.PairName, "secure.tls.pair-name", o.Tls.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&o.Tls.CertFile, "secure.tls.cert-file", o.Tls.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&o.Tls.KeyFile, "secure.tls.key-file", o.Tls.KeyFile, "File containing the default x509 "+
		"private key matching --secure.tls.cert-key.cert-file.")
}
