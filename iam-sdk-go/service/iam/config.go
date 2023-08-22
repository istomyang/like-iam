package iam

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"istomyang.github.com/like-iam/component-base/auth"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/coder"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/url"
	"os"
	"time"
)

type Config struct {
	Authentication `yaml:",inline"`

	// ApiAddr and AuthAddr 's format is scheme://host:port.
	ApiAddr  string `yaml:"api-addr,omitempty"`
	AuthAddr string `yaml:"auth-addr,omitempty"`

	restClientConfig *RestClientConfig `yaml:"rest-client-config,omitempty"`

	// See coder.CodeJson
	CodeOption string `yaml:"code-option,omitempty"`

	// ApiVersion affects IAM.Api and IAM.Authz, which auto-select by ApiVersion.
	ApiVersion int `yaml:"api-version,omitempty"`

	// UserAgent will append default user-agent string.
	// optional
	UserAgent string `yaml:"user-agent,omitempty"`
}

func (c *Config) Correct() []error {
	var errs []error
	errs = append(errs, c.restClientConfig.Correct()...)
	{
		a := string(coder.CodeJson)
		b := string(coder.CodeMsgPack)
		if c.CodeOption != a && c.CodeOption != b {
			errs = append(errs, fmt.Errorf("code option must be %s or %s", a, b))
		}
	}
	{
		// Check host.
		var err error
		if c.ApiAddr, err = url.PickAddr(c.ApiAddr); err != nil {
			errs = append(errs, err)
		}
		if c.AuthAddr, err = url.PickAddr(c.AuthAddr); err != nil {
			errs = append(errs, err)
		}
	}
	{
		if c.ApiVersion == 0 {
			c.ApiVersion = 1
		}
		if c.ApiVersion < 0 {
			errs = append(errs, fmt.Errorf("ApiVersion must be greater than 0 and impl corresponding api, got: %d", c.ApiVersion))
		}
	}
	return errs
}

type RestClientConfig struct {
	Timeout int64 `yaml:"timeout-ms,omitempty"`

	// EnableClientAuthn enable server-client-two-way authn, client must assign client.key、client.crt and ca.crt.
	EnableClientAuthn bool   `yaml:"enable-client-authn"`
	CertPEMFile       string `yaml:"cert-pem-file,omitempty"`
	CertPEMData       []byte `yaml:"cert-pem-data,omitempty"`
	KeyPEMData        []byte `yaml:"key-pem-data,omitempty"`
	KeyPEMFile        string `yaml:"key-pem-file,omitempty"`
	CAFile            string `yaml:"ca-file,omitempty"`
	CAData            []byte `yaml:"ca-data,omitempty"`
}

func (c *RestClientConfig) Correct() []error {
	var errs []error
	{
		if c.Timeout < 1000 {
			c.Timeout = 1000
		}
	}
	{
		if c.EnableClientAuthn {
			if c.CertPEMFile == "" || c.CertPEMData == nil {
				errs = append(errs, fmt.Errorf("must assign cert pem config"))
			}
			if c.KeyPEMFile == "" || c.KeyPEMData == nil {
				errs = append(errs, fmt.Errorf("must assign cert key config"))
			}
			// optional
			//if c.CAFile == "" || c.CAData == nil {
			//	errs = append(errs, fmt.Errorf("must assign cert ca config"))
			//}
		}
	}

	return errs
}

type Authentication struct {
	DevelopmentEnv  bool   `yaml:"development-env,omitempty"`
	Username        string `yaml:"username,omitempty"`
	Password        string `yaml:"password,omitempty"`
	SecretID        string `yaml:"secret-id,omitempty"`
	SecretKey       string `yaml:"secret-key,omitempty"`
	BearerTokenFile string `yaml:"bearer-token-file,omitempty"`
	BearerToken     string `yaml:"bearer-token,omitempty"`
}

// Header returns corresponding Token Header by different config.
// Priority: BearerToken > SecretID/KEY > Username&Password
func (au *Authentication) Header() (k string, v string, err error) {
	k = "Authorization"

	if au.BearerToken != "" {
		v = fmt.Sprintf("Bearer %s", au.BearerToken)
		return
	}
	if au.BearerTokenFile != "" {
		var b []byte
		if b, err = os.ReadFile(au.BearerTokenFile); err != nil {
			return
		}
		v = fmt.Sprintf("Bearer %s", string(b))
		return
	}
	if au.SecretID != "" && au.SecretKey != "" {
		var expire time.Duration
		if au.DevelopmentEnv {
			expire = auth.JwtExpireDevelopment
		} else {
			expire = auth.JwtExpireProduction
		}
		var token string
		token, err = auth.Sign(au.SecretID, au.SecretKey, "iam-sdk-go", ".like-iam.github.com", expire)
		v = fmt.Sprintf("Bearer %s", token)
	}
	if au.Username != "" && au.Password != "" {
		token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", au.Username, au.Password)))
		v = fmt.Sprintf("Basic %s", token)
		return
	}

	s, _ := json.Marshal(au)
	err = fmt.Errorf("you should complete auth config, got: %s", string(s))
	return
}
