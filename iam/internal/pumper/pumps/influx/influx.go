package influx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"time"
)

type influx struct {
	ctx    context.Context
	cancel context.CancelFunc

	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter

	client *conn.InfluxClient
	config *options.InfluxOpts
}

func New() pumps.Pump {
	return &influx{}
}

func (i *influx) Init(ctx context.Context, config map[string]any) error {
	i.ctx, i.cancel = context.WithCancel(ctx)
	i.config = &options.InfluxOpts{}
	if err := mapstructure.Decode(config, i.config); err != nil {
		return err
	}
	return nil
}

func (i *influx) Run() error {
	i.client = conn.NewInfluxClient(i.ctx, i.config)
	return i.client.Run()
}

func (i *influx) Close() error {
	i.cancel()
	return i.client.Close()
}

func (i *influx) GetName() string {
	return "Influx Pump"
}

func (i *influx) Write(ds []interface{}) error {
	for _, d := range ds {
		var info, ok = d.(*analytics.RecordInfo)
		if !ok {
			return fmt.Errorf("data as RecordInfo fail: %v", d)
		}

		var mapping = map[string]any{
			"username":   info.UserName,
			"deciders":   info.Deciders,
			"conclusion": info.Conclusion,
			"expire-at":  info.ExpireAt,
			"request":    info.Request,
			"effect":     info.Effect,
			"timestamp":  info.Timestamp,
			"policies":   info.Policies,
		}

		var fields = make(map[string]any)
		var tags = make(map[string]string)

		for _, k := range i.config.Fields {
			var has bool
			if fields[k], has = mapping[k]; !has {
				return fmt.Errorf("this field key does not exist: %s, you should assign correct fileds", k)
			}
		}

		for _, k := range i.config.Tags {

			v, has := mapping[k]
			if !has {
				return fmt.Errorf("this tag key does not exist: %s, you should assign correct tags", k)
			}

			var sv string
			svb, err := json.Marshal(v)
			if err != nil {
				sv = "CAN'T JSON"
			} else {
				sv = string(svb)
			}

			tags[k] = sv
		}

		return i.client.Send("analytics", tags, fields)
	}

	return nil
}

func (i *influx) SetFilter(filter *pumps.Filter) {
	i.filter = filter
}

func (i *influx) GetFilter() *pumps.Filter {
	return i.filter
}

func (i *influx) SetTimeout(duration time.Duration) {
	i.timeout = duration
}

func (i *influx) GetTimeout() time.Duration {
	return i.timeout
}

func (i *influx) SetOmitDetail(b bool) {
	i.omitDetail = b
}

func (i *influx) GetOmitDetail() bool {
	return i.omitDetail
}

var _ pumps.Pump = &influx{}
