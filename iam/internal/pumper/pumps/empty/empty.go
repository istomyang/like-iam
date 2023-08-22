package empty

import (
	"context"
	"istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"istomyang.github.com/like-iam/log"
	"time"
)

type emptyDump struct {
	timeout    time.Duration
	omitDetail bool
	filter     *pumps.Filter
}

func New() pumps.Pump {
	return &emptyDump{}
}

func (e *emptyDump) Init(ctx context.Context, config map[string]any) error {
	return nil
}

func (e *emptyDump) SetFilter(filter *pumps.Filter) {
	e.filter = filter
}

func (e *emptyDump) GetFilter() *pumps.Filter {
	return e.filter
}

func (e *emptyDump) SetTimeout(duration time.Duration) {
	e.timeout = duration
}

func (e *emptyDump) GetTimeout() time.Duration {
	return e.timeout
}

func (e *emptyDump) SetOmitDetail(b bool) {
	e.omitDetail = b
}

func (e *emptyDump) GetOmitDetail() bool {
	return e.omitDetail
}

func (e *emptyDump) Run() error {
	return nil
}

func (e *emptyDump) Close() error {
	return nil
}

func (e *emptyDump) GetName() string {
	return "Empty Pump"
}

func (e *emptyDump) Write(i []interface{}) error {
	log.Debugf("%s: write %s data.", e.GetName(), len(i))
	return nil
}

var _ pumps.Pump = &emptyDump{}
