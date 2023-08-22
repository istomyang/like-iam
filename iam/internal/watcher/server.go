package watcher

import (
	"context"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/shutdown"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store/mysql"
)

type watchServer struct {
	ctx      context.Context
	shutdown *shutdown.Shutdown
	options  *Options

	watch *watcher
}

func newWatchServer(ctx context.Context, options *Options) *watchServer {
	var s watchServer
	s.ctx = ctx
	s.shutdown = shutdown.CreateDefaultShutdown(s.close)
	s.watch = newWatcher(ctx, options.Watcher)
	return &s
}

func (w *watchServer) run() {
	w.shutdown.Run()
	conn.NewRedisClientOr(w.options.RedisOptions)
	if err := conn.GetRedisClient().Run(); err != nil {
		panic(err)
	}
	if _, err := mysql.GetMySQLFactoryOr(w.options.MysqlOptions); err != nil {
		panic(err)
	}
	if err := w.watch.run(); err != nil {
		panic(err)
	}
}

func (w *watchServer) close() error {
	if err := conn.GetRedisClient().Close(); err != nil {
		return err
	}
	c, err := mysql.GetMySQLFactoryOr(nil)
	if err != nil {
		return err
	}
	return c.Close()
}
