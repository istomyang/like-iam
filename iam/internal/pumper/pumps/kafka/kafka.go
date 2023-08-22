package kafka

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

type kafka struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *options.KafkaOpts

	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter

	client *conn.KafkaClient
}

func New() pumps.Pump {
	return &kafka{}
}

func (k *kafka) Init(ctx context.Context, config map[string]any) error {
	k.config = &options.KafkaOpts{}
	k.ctx, k.cancel = context.WithCancel(ctx)
	if err := mapstructure.Decode(k.config, config); err != nil {
		return err
	}
	return nil
}

func (k *kafka) Run() error {
	k.client = conn.NewKafkaClient(k.ctx, k.config)
	return nil
}

func (k *kafka) Close() error {
	k.cancel()
	return k.client.Close()
}

func (k *kafka) GetName() string {
	return "Kafka Pump"
}

func (k *kafka) Write(ds []any) error {
	var messages = make([][]byte, len(ds))

	for i, data := range ds {

		info, ok := data.(*analytics.RecordInfo)
		if !ok {
			return fmt.Errorf("fail to transfer to RecordInfo, got %v", data)
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

		for k, v := range k.config.MessageExtend {
			mapping[k] = v
		}

		message, err := json.Marshal(mapping)
		if err != nil {
			return err
		}

		messages[i] = message
	}

	if err := k.client.Send(messages); err != nil {
		return err
	}
	return nil
}

func (k *kafka) SetFilter(filter *pumps.Filter) {
	k.filter = filter
}

func (k *kafka) GetFilter() *pumps.Filter {
	return k.filter
}

func (k *kafka) SetTimeout(duration time.Duration) {
	k.timeout = duration
}

func (k *kafka) GetTimeout() time.Duration {
	return k.timeout
}

func (k *kafka) SetOmitDetail(b bool) {
	k.omitDetail = b
}

func (k *kafka) GetOmitDetail() bool {
	return k.omitDetail
}

var _ pumps.Pump = &kafka{}
