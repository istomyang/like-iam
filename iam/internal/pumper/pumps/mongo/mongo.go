package mongo

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"time"
)

type mongo struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *options.MongoOpts
	client *conn.MongoClient

	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter
}

func New() pumps.Pump {
	return &mongo{}
}

func (m *mongo) Init(ctx context.Context, config map[string]any) error {
	m.ctx, m.cancel = context.WithCancel(ctx)
	m.config = &options.MongoOpts{}
	if err := mapstructure.Decode(m.config, config); err != nil {
		return err
	}
	return nil
}

func (m *mongo) Run() error {
	m.client = conn.NewMongoClient(m.ctx, m.config)
	return m.client.Run()
}

func (m *mongo) Close() error {
	return m.client.Close()
}

func (m *mongo) GetName() string {
	return "MongoDB Pump"
}

func (m *mongo) Write(ds []any) (err error) {
	var documents = make([]map[string]any, len(ds))
	for i, d := range ds {
		info, ok := d.(*analytics.RecordInfo)
		if ok {
			err = fmt.Errorf("fail to transfer to RecordInfo, got %v", d)
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
		documents[i] = mapping
	}
	return m.client.Send(documents)
}

func (m *mongo) SetFilter(filter *pumps.Filter) {
	m.filter = filter
}

func (m *mongo) GetFilter() *pumps.Filter {
	return m.filter
}

func (m *mongo) SetTimeout(duration time.Duration) {
	m.timeout = duration
}

func (m *mongo) GetTimeout() time.Duration {
	return m.timeout
}

func (m *mongo) SetOmitDetail(b bool) {
	m.omitDetail = b
}

func (m *mongo) GetOmitDetail() bool {
	return m.omitDetail
}

var _ pumps.Pump = &mongo{}
