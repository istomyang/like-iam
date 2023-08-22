package conn

import (
	"context"
	"fmt"
	influx "github.com/influxdata/influxdb/client/v2"
	"istomyang.github.com/like-iam/component/pkg/interfaces"
	"istomyang.github.com/like-iam/component/pkg/options"
	"sync"
	"time"
)

type InfluxClient struct {
	ctx          context.Context
	cancel       context.CancelFunc
	config       *options.InfluxOpts
	client       influx.HTTPClient
	mut          sync.Mutex
	retry        int // reconnect counts.
	defaultRetry int
}

func NewInfluxClient(ctx context.Context, opts *options.InfluxOpts) *InfluxClient {
	var c InfluxClient
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.config = opts
	c.defaultRetry = 3
	return &c
}

var _ interfaces.ComponentCommon = &InfluxClient{}

func (c *InfluxClient) Run() error {
	if err := c.connect(); err != nil {
		return err
	}
	return nil
}

func (c *InfluxClient) Close() error {
	c.cancel()
	return c.client.Close()
}

func (c *InfluxClient) Send(name string, tags map[string]string, fields map[string]any) error {
	points, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Precision: "ms",
		Database:  c.config.DatabaseName,
	})
	if err != nil {
		return err
	}
	point, err := influx.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return err
	}
	points.AddPoint(point)
	return c.retrySending(points)
}

func (c *InfluxClient) retrySending(points influx.BatchPoints) (err error) {
	if err = c.client.WriteCtx(c.ctx, points); err != nil {
		err = nil
		if err = c.Reconnect(); err != nil {
			return fmt.Errorf("send fail: %s", err.Error())
		} else {
			return c.retrySending(points)
		}
	}
	return nil
}

func (c *InfluxClient) Reconnect() error {
	if c.retry > c.defaultRetry {
		return fmt.Errorf("influx client doesn't connect to server, reconnect for %d counts/", c.retry)
	}

	if c.client == nil {
		if err := c.connect(); err != nil {
			c.retry++
			time.Sleep(time.Second * 5)
			return c.Reconnect()
		}
	}

	if _, _, err := c.client.Ping(time.Second * 5); err != nil {
		err = nil
		if err = c.connect(); err != nil {
			c.retry++
			time.Sleep(time.Second * 5)
			return c.Reconnect()
		}
	}

	c.retry = 0
	return nil
}

func (c *InfluxClient) connect() (err error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     c.config.Addr,
		Username: c.config.Username,
		Password: c.config.Password,
	})
	return
}
