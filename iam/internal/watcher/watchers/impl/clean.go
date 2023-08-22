package impl

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store/mysql"
	"istomyang.github.com/like-iam/iam/internal/watcher"
	"istomyang.github.com/like-iam/iam/internal/watcher/watchers"
	"istomyang.github.com/like-iam/log"
)

type cleaner struct {
	mut    *redsync.Mutex
	ctx    context.Context
	config *watcher.WatchOpts
}

func (c *cleaner) Init(ctx context.Context, mut *redsync.Mutex, config interface{}) {
	c.ctx = ctx
	c.mut = mut
	cfg, ok := config.(*watcher.WatchOpts)
	if !ok {
		panic(watchers.ErrConfig)
	}
	c.config = cfg
}

func (c *cleaner) Schedules() string {
	return "@every 1d"
}

func (c *cleaner) Run() {
	if err := c.mut.Lock(); err != nil {
		log.Warnf("clean watcher already run, got err: %s", err.Error())
		return
	}
	if _, err := c.mut.Unlock(); err != nil {
		log.Errorf("clean watcher could not release lock, got err; %s", err.Error())
		return
	}

	factory, err := mysql.GetMySQLFactoryOr(nil)
	if err != nil {
		log.Errorf("clean watcher get mysql factory got err: %s", err.Error())
		return
	}

	effectCounts, err := factory.Policy().ClearOutdated(c.ctx, c.config.Clean.MaxUserActiveDays)
	if err != nil {
		log.Errorf("clean watcher send clean-outdated to factory got err: %s", err.Error())
		return
	}

	log.Infof("clean watcher clean outdated policy for %d numbers.", effectCounts)
}

var _ watchers.Watcher = &cleaner{}

func init() {
	watchers.Register("clean watcher", &cleaner{})
}
