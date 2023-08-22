package conn

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"istomyang.github.com/like-iam/component/pkg/interfaces"
	"istomyang.github.com/like-iam/component/pkg/options"
	"net/http"
)

type PrometheusClient struct {
	ctx          context.Context
	cancel       context.CancelFunc
	config       *options.PrometheusOpts
	retry        int
	defaultRetry int

	authz *prometheus.CounterVec
}

func NewPrometheusClient(ctx context.Context, opts *options.PrometheusOpts) *PrometheusClient {
	var p PrometheusClient
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.config = opts
	p.defaultRetry = 3
	return &p
}

func (p *PrometheusClient) Connect() error {
	reg := prometheus.NewRegistry()
	p.withAuthz(reg)
	http.Handle(p.config.Path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go func() {
		panic(http.ListenAndServe(p.config.Addr, nil))
	}()
	return nil
}

func (p *PrometheusClient) withAuthz(registry *prometheus.Registry) {
	authz := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "iam_user_authorization_status_total",
		Help: "authorization effect per user",
	}, []string{"isAllowed", "username"})
	p.authz = authz
	registry.MustRegister(authz)
}

func (p *PrometheusClient) Send(allow, username string) {
	p.authz.WithLabelValues(allow, username).Inc()
}

func (p *PrometheusClient) Run() error {
	return p.Connect()
}

func (p *PrometheusClient) Close() error {
	p.cancel()
	return nil
}

var _ interfaces.ComponentCommon = &PrometheusClient{}
