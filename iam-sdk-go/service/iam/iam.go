package iam

import (
	"context"
	"fmt"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/client"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/coder"
	v1 "istomyang.github.com/like-iam/iam-sdk-go/service/iam/v1"
	"os"
	"strings"
	"time"
)

type IAM interface {
	// Api and Authz returns depends on Config.ApiVersion.
	Api() v1.Api
	Authz() v1.Authz

	ApiV1() v1.Api
	AuthzV1() v1.Authz
}

const CoderRegisterName = "iam"

type iam struct {
	ctx    context.Context
	config *Config
}

func NewIam(ctx context.Context, config *Config) IAM {
	_ = coder.Register(CoderRegisterName, coder.CodeJson)
	return &iam{
		ctx:    ctx,
		config: config,
	}
}

func (i *iam) Api() v1.Api {
	switch i.config.ApiVersion {
	case 1:
		return i.ApiV1()
	default:
		panic(fmt.Sprintf("ApiVerson %d must be implmented", i.config.ApiVersion))
	}
}

func (i *iam) Authz() v1.Authz {
	switch i.config.ApiVersion {
	case 1:
		return i.AuthzV1()
	default:
		panic(fmt.Sprintf("ApiVerson %d must be implmented", i.config.ApiVersion))
	}
}

func (i *iam) ApiV1() v1.Api {
	c, errs := i.buildCommonConfig()
	if errs != nil {
		panic(fmt.Sprintf("client config has error: %s", i.errs2String(errs)))
	}
	c.BaseUrl = i.config.ApiAddr
	if errs = c.Correct(); errs != nil {
		panic(fmt.Sprintf("client config has error: %s", i.errs2String(errs)))
	}
	return v1.NewApiV1(client.NewClient(i.ctx, c))
}

func (i *iam) AuthzV1() v1.Authz {
	c, errs := i.buildCommonConfig()
	if errs != nil {
		panic(fmt.Sprintf("client config has error: %s", i.errs2String(errs)))
	}
	c.BaseUrl = i.config.AuthAddr
	if errs := c.Correct(); errs != nil {
		panic(fmt.Sprintf("client config has error: %s", i.errs2String(errs)))
	}
	return v1.NewAuthz(client.NewClient(i.ctx, c))
}

func (i *iam) buildCommonConfig() (*client.Config, []error) {
	var errs []error

	errs = append(errs, i.config.Correct()...)

	// In this step, source has error, so next step is meaningless.
	if len(errs) > 0 {
		return nil, errs
	}

	_, token, err := i.config.Header()
	if err != nil {
		errs = append(errs, err)
		// you should prepare username and password.
		return nil, errs
	}

	c := client.Config{
		BaseUrl:           "",
		Token:             token,
		CoderName:         CoderRegisterName,
		Timeout:           time.Duration(i.config.restClientConfig.Timeout) * time.Microsecond,
		EnableClientAuthn: i.config.restClientConfig.EnableClientAuthn,
		CertPEMData:       i.config.restClientConfig.CertPEMData,
		KeyPEMData:        i.config.restClientConfig.KeyPEMData,
		CAData:            i.config.restClientConfig.CAData,
	}

	if c.EnableClientAuthn {
		if c.CertPEMData == nil {
			var err error
			c.CertPEMData, err = os.ReadFile(i.config.restClientConfig.CertPEMFile)
			errs = append(errs, err)
		}
		if c.KeyPEMData == nil {
			var err error
			c.KeyPEMData, err = os.ReadFile(i.config.restClientConfig.KeyPEMFile)
			errs = append(errs, err)
		}
		if c.CAData == nil && i.config.restClientConfig.CAFile != "" {
			var err error
			c.CAData, err = os.ReadFile(i.config.restClientConfig.CAFile)
			errs = append(errs, err)
		}
	}

	return &c, errs
}

func (i *iam) errs2String(errs []error) string {
	var builder strings.Builder
	builder.WriteString("errors: ")
	for _, err := range errs {
		builder.WriteString(err.Error())
		builder.WriteString(";")
	}
	return builder.String()
}

var _ IAM = &iam{}
