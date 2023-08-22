package watcher

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/robfig/cron/v3"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/iam/internal/pkg"
	"istomyang.github.com/like-iam/iam/internal/watcher/watchers"
	_ "istomyang.github.com/like-iam/iam/internal/watcher/watchers/impl"
	"istomyang.github.com/like-iam/log"
	"time"
)

type watcher struct {
	ctx     context.Context
	cancel  context.CancelFunc
	cron    *cron.Cron
	log     *log.Logger
	options *WatchOpts
	rs      *redsync.Redsync
}

func newWatcher(ctx context.Context, options *WatchOpts) *watcher {
	var w = &watcher{}
	w.ctx, w.cancel = context.WithCancel(ctx)

	client := conn.GetRedisClient()

	w.log = log.Default()
	w.cron = cron.New(cron.WithLogger(w), cron.WithChain(cron.Recover(w)), cron.WithSeconds())
	w.rs = redsync.New(goredis.NewPool(client.UniversalClient()))
	w.options = options

	return w
}

func (w *watcher) run() error {
	for name, wat := range watchers.ListMap() {
		ctx := context.WithValue(w.ctx, pkg.WatcherContextKey, name)
		mut := w.rs.NewMutex(name, redsync.WithExpiry(time.Hour))
		wat.Init(ctx, mut, w.options.Clean)
		if _, err := w.cron.AddJob(wat.Schedules(), wat); err != nil {
			return err
		}
	}

	go func() {
		w.cron.Start()

		<-w.ctx.Done()

		w.cron.Stop()
	}()

	return nil
}

func (w *watcher) Info(msg string, keysAndValues ...interface{}) {
	w.log.Infof(msg, keysAndValues)
}

func (w *watcher) Error(err error, msg string, keysAndValues ...interface{}) {
	w.log.Errorf(msg, keysAndValues)
}
