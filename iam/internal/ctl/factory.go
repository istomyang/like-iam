package ctl

import (
	"context"
	"github.com/spf13/pflag"
	"istomyang.github.com/like-iam/iam-sdk-go/pkg/util/coder"
	"istomyang.github.com/like-iam/iam-sdk-go/service"
	"istomyang.github.com/like-iam/iam-sdk-go/service/iam"
	"strings"
)

type factory struct {
	ctx             context.Context
	username        string
	password        string
	secretId        string
	secretKey       string
	bearerToken     string
	bearerTokenFile string
	apiAddr         string
	authAddr        string
	apiVersion      int
	userAgent       string

	svr service.Services
}

var defaultFactory = &factory{}

func (f *factory) AddFlagTo(fs *pflag.FlagSet) {
	fs.StringVar(&f.username, "username", "", "Optional.")
	fs.StringVar(&f.password, "password", "", "Optional.")
	fs.StringVar(&f.secretId, "secretId", "", "Optional.")
	fs.StringVar(&f.secretKey, "secret-key", "", "Optional.")
	fs.StringVar(&f.bearerToken, "bearer-token", "", "Optional.")
	fs.StringVar(&f.bearerTokenFile, "bearer-token-file", "", "Optional.")
	fs.StringVar(&f.apiAddr, "api-addr", "", "Example: host:port")
	fs.StringVar(&f.authAddr, "auth-addr", "", "Example: host:port")
	fs.IntVar(&f.apiVersion, "api-version", 1, "default is 1")
	fs.StringVar(&f.userAgent, "user-agent", "", "Optional.")
}

func (f *factory) Service() service.Services {
	if f.svr != nil {
		return f.svr
	}

	f.apiAddr = strings.TrimSpace(f.apiAddr)
	f.authAddr = strings.TrimSpace(f.authAddr)
	f.userAgent = strings.TrimSpace(f.userAgent)

	var authConfig iam.Authentication
	if f.bearerToken != "" {
		authConfig = iam.Authentication{BearerToken: f.bearerToken}
	} else if f.bearerTokenFile != "" {
		authConfig = iam.Authentication{BearerTokenFile: f.bearerTokenFile}
	} else if f.secretId != "" && f.secretKey != "" {
		authConfig = iam.Authentication{SecretID: f.secretId, SecretKey: f.secretKey}
	} else if f.username != "" && f.password != "" {
		authConfig = iam.Authentication{Username: f.username, Password: f.password}
	}

	var svr = service.NewService(f.ctx, &service.Config{IAM: &iam.Config{
		Authentication: authConfig,
		ApiAddr:        f.apiAddr,
		AuthAddr:       f.authAddr,
		CodeOption:     string(coder.CodeJson),
		ApiVersion:     f.apiVersion,
		UserAgent:      f.userAgent,
	}})
	f.svr = svr
	return svr
}
