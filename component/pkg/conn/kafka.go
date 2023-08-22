package conn

import (
	"context"
	"github.com/segmentio/kafka-go"
	"istomyang.github.com/like-iam/component/pkg/interfaces"
	"istomyang.github.com/like-iam/component/pkg/options"
	"strings"
	"time"
)

type KafkaClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *options.KafkaOpts

	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaClient(ctx context.Context, opts *options.KafkaOpts) *KafkaClient {
	var c KafkaClient
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.config = opts
	return &c
}

func (k *KafkaClient) connect() error {
	var c = k.config
	k.writer = &kafka.Writer{
		Addr:  kafka.TCP(strings.Split(c.Addrs, ",")...),
		Topic: c.Topic,
	}
	return nil
}

func (k *KafkaClient) Send(batch [][]byte) error {
	var messages = make([]kafka.Message, len(batch))
	for i, data := range batch {
		messages[i] = kafka.Message{
			Topic: k.config.Topic,
			Value: data,
			Time:  time.Now(),
		}
	}

	if err := k.writer.WriteMessages(k.ctx, messages...); err != nil {
		return err
	}
	return nil
}

func (k *KafkaClient) Run() error {
	return k.connect()
}

func (k *KafkaClient) Close() error {
	k.cancel()
	return k.writer.Close()
}

var _ interfaces.ComponentCommon = &KafkaClient{}
