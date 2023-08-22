package prometheus

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/ladon"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"time"
)

type prometheus struct {
	ctx    context.Context
	cancel context.CancelFunc

	config *options.PrometheusOpts
	client *conn.PrometheusClient

	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter
}

func New() pumps.Pump {
	return &prometheus{}
}

func (p *prometheus) Init(ctx context.Context, config map[string]any) error {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.config = &options.PrometheusOpts{}
	return mapstructure.Decode(p.config, config)
}

func (p *prometheus) Run() error {
	p.client = conn.NewPrometheusClient(p.ctx, p.config)
	return nil
}

func (p *prometheus) Close() error {
	return p.client.Close()
}

func (p *prometheus) GetName() string {
	return "Prometheus Pump"
}

func (p *prometheus) Write(ds []any) error {
	for _, d := range ds {
		info, ok := d.(*analytics.RecordInfo)
		if !ok {
			return fmt.Errorf("fail to transfer to RecordInfo, got %v", d)
		}
		var allowed string
		if info.Effect == ladon.AllowAccess {
			allowed = "allow"
		} else {
			allowed = "deny"
		}
		p.client.Send(allowed, info.UserName)
	}
	return nil
}

func (p *prometheus) SetFilter(filter *pumps.Filter) {
	p.filter = filter
}

func (p *prometheus) GetFilter() *pumps.Filter {
	return p.filter
}

func (p *prometheus) SetTimeout(duration time.Duration) {
	p.timeout = duration
}

func (p *prometheus) GetTimeout() time.Duration {
	return p.timeout
}

func (p *prometheus) SetOmitDetail(b bool) {
	p.omitDetail = b
}

func (p *prometheus) GetOmitDetail() bool {
	return p.omitDetail
}

var _ pumps.Pump = &prometheus{}
