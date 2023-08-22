package pumper

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"istomyang.github.com/like-iam/component-base/json"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/shutdown"
	"istomyang.github.com/like-iam/iam/internal/pkg"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	pps "istomyang.github.com/like-iam/iam/internal/pumper/pumps"
	"istomyang.github.com/like-iam/iam/internal/pumper/store"
	"istomyang.github.com/like-iam/iam/internal/pumper/store/redis"
	"istomyang.github.com/like-iam/log"
	"strings"
	"time"
)

type pumper struct {
	mut      *redsync.Mutex
	interval time.Duration

	shutdown *shutdown.Shutdown
	ctx      context.Context
	cancel   context.CancelFunc

	pumpsOptions map[string]any
}

func newPumper(options *Options) *pumper {
	r := &pumper{}

	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.shutdown = shutdown.CreateDefaultShutdown(r.close)

	conn.NewRedisClientOr(options.redisOptions)
	r.mut = redsync.New(goredis.NewPool(conn.GetRedisClient().UniversalClient())).NewMutex(pkg.AnalyticsMutexKey, redsync.WithExpiry(options.mutexExpiry))
	r.pumpsOptions = options.pumps
	r.interval = options.interval

	return r
}

func (p *pumper) run() error {
	p.shutdown.Run()

	store.SetClient(redis.NewRedisStore(p.ctx))
	if err := store.Client().Run(); err != nil {
		return err
	}

	if err := p.errs2err(pps.RegisterRun(p.ctx, p.pumpsOptions), "pumpers run error: "); err != nil {
		return err
	}

	tkr := time.NewTicker(p.interval)
	defer tkr.Stop()

	select {
	case <-tkr.C:
		if err := p.doWork(); err != nil {
			return err
		}
	case <-p.ctx.Done():
		return nil
	}
	return nil
}

func (p *pumper) close() error {
	p.cancel()
	if err := conn.GetRedisClient().Close(); err != nil {
		return err
	}
	var errs = pps.Close()
	return p.errs2err(errs, "pumpers close error: ")
}

func (p *pumper) doWork() error {
	if err := p.mut.Lock(); err != nil {
		return err
	}
	defer func() {
		if _, err := p.mut.Unlock(); err != nil {
			log.Errorf("redis mutex must be unlock, got err: %s", err.Error())
			return
		}
	}()

	raw, err := store.Client().Pop(pkg.AnalyticsKey)
	if err != nil {
		return err
	}

	// TODO: maybe not pointer.
	dataSet := make([]*analytics.RecordInfo, len(raw))
	for i, bytes := range raw {
		var info = &analytics.RecordInfo{}
		if err = json.MPUnmarshal(bytes, info); err != nil {
			return err
		}
		dataSet[i] = info
	}

	pps.Do(func(pump pps.Pump) error {
		var sendInfos []any
		for _, info := range dataSet {
			if pump.GetFilter().ShouldFilter(info) {
				continue
			}
			if pump.GetOmitDetail() {
				info.Deciders = ""
				info.Policies = ""
			}
			sendInfos = append(sendInfos, info)
		}

		if err = pump.Write(sendInfos); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (p *pumper) errs2err(errs []error, prefix string) error {
	if len(errs) == 0 {
		return nil
	}
	var builder strings.Builder
	builder.WriteString(strings.TrimSpace(prefix) + " ")
	for _, err := range errs {
		builder.WriteString(strings.TrimSpace(err.Error()))
		builder.WriteString("; ")
	}
	return fmt.Errorf(builder.String())
}
