package client

import (
	"time"
)

type Config struct {
	// BaseUrl 's format: scheme://host:port
	BaseUrl string `yaml:"base-url,omitempty" mapstructure:"base-url"`

	// Token is Http Header `Authorization` 's value
	Token string

	// CoderName is service name to get corresponding base.Coder.
	CoderName string

	// Timeout represents microsecond for http client.
	Timeout time.Duration `yaml:"timeout,omitempty" mapstructure:"timeout"`

	// EnableClientAuthn enable server-client-two-way authn, client must assign client.key、client.crt and ca.crt.
	EnableClientAuthn bool   `json:"enable-client-authn,omitempty" mapstructure:"enable-client-authn"`
	Insecure          bool   `yaml:"insecure,omitempty" mapstructure:"insecure"`
	CertPEMData       []byte `json:"cert-pem-data,omitempty" mapstructure:"cert-pem-data"`
	KeyPEMData        []byte `json:"key-pem-data,omitempty" mapstructure:"key-pem-data"`
	CAData            []byte `json:"ca-data,omitempty" mapstructure:"ca-data"`
}

// Correct will check and correct config.
func (c *Config) Correct() []error {
	var errs []error
	// Do work in upper layer.
	return errs
}
