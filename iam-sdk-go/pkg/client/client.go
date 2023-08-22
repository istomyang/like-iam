package client

import (
	"context"
	"net/url"
	"time"
)

type client struct {
	timeout   time.Duration
	coderName string

	c *Config

	ctx    context.Context
	cancel context.CancelFunc

	u *url.URL
}

// NewClient create a reset client.
func NewClient(ctx context.Context, config *Config) Client {
	var c client
	c.timeout = config.Timeout * time.Microsecond
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.c = config
	c.coderName = config.CoderName
	return &c
}

var _ Client = &client{}

func (c *client) createBaseRequest() Request {
	var req = newRequest(c.u, c.coderName, c.c.CertPEMData, c.c.KeyPEMData, c.c.CAData, c.c.Insecure)
	req.Header("Authorization", c.c.Token)
	return req
}

func (c *client) Verb(method Verb) Request {
	return c.createBaseRequest().Verb(method)
}

func (c *client) Get() Request {
	return c.createBaseRequest().Verb(VerbGET)
}

func (c *client) Post() Request {
	return c.createBaseRequest().Verb(VerbPost)
}

func (c *client) Put() Request {
	return c.createBaseRequest().Verb(VerbPUT)
}

func (c *client) Delete() Request {
	return c.createBaseRequest().Verb(VerbDelete)
}
