package sysloger

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"log/syslog"
	"regexp"
	"time"
)

type sysloger struct {
	ctx        context.Context
	cancel     context.CancelFunc
	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter
	config     *config

	client *syslog.Writer
}

func New() pumps.Pump {
	return &sysloger{}
}

type config struct {
	// Addr looks like "localhost:1234"
	Addr string
	// LogLevel can refer to sysloger.LOG_DEBUG
	LogLevel int
	// Tag default is "sysloger-pump"
	Tag string
}

func (c *config) Check() {
	if c.Tag != "" {
		c.Tag = "sysloger-pump"
	}
	re := regexp.MustCompile(`(.*):/d+`)
	c.Addr = re.FindString(c.Addr)
	if c.LogLevel == 0 {
		c.LogLevel = int(syslog.LOG_ALERT)
	}
}

func (s *sysloger) Init(ctx context.Context, options map[string]any) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.config = &config{}
	return mapstructure.Decode(s.config, options)
}

func (s *sysloger) Run() error {
	s.config.Check()
	var err error
	s.client, err = syslog.Dial("tcp", s.config.Addr, syslog.Priority(s.config.LogLevel), s.config.Tag)
	if err != nil {
		return err
	}
	return nil
}

func (s *sysloger) Close() error {
	return s.client.Close()
}

func (s *sysloger) GetName() string {
	return "Syslog Pump"
}

func (s *sysloger) Write(ds []any) error {
	for _, d := range ds {
		info, ok := d.(*analytics.RecordInfo)
		if ok {
			return fmt.Errorf("fail to transfer to RecordInfo, got %v", d)
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

		select {
		case <-s.ctx.Done():
			return nil
		default:
			_, err := fmt.Fprintf(s.client, "%s", mapping)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *sysloger) SetFilter(filter *pumps.Filter) {
	s.filter = filter
}

func (s *sysloger) GetFilter() *pumps.Filter {
	return s.filter
}

func (s *sysloger) SetTimeout(duration time.Duration) {
	s.timeout = duration
}

func (s *sysloger) GetTimeout() time.Duration {
	return s.timeout
}

func (s *sysloger) SetOmitDetail(b bool) {
	s.omitDetail = b
}

func (s *sysloger) GetOmitDetail() bool {
	return s.omitDetail
}

var _ pumps.Pump = &sysloger{}
