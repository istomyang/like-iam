package analytics

import (
	"context"
	"fmt"
	"istomyang.github.com/like-iam/iam/internal/authzserver/analytics/store"
	"istomyang.github.com/like-iam/iam/internal/pkg"
	"istomyang.github.com/like-iam/iam/internal/pkg/analytics"
	"sync"
	"sync/atomic"
	"time"
)

type Analytics struct {
	key    string // key is key in kv storage.
	expire time.Duration

	store         Store
	pool          chan []byte
	workers       int
	batch         uint64
	flushInterval time.Duration
	flag          uint32 // 0: stop, 1: start
	wg            sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

var (
	client *Analytics
	once   sync.Once
)

func GetAnalytics() *Analytics {
	return client
}

func NewAnalytics(ctx context.Context, options *Options) (*Analytics, error) {
	once.Do(func() {
		a := Analytics{
			pool:          make(chan []byte, options.Workers),
			workers:       options.Workers,
			batch:         options.BatchSize,
			flushInterval: options.FlushInterval,
			flag:          0,
			expire:        options.Expire,
		}
		a.ctx, a.cancel = context.WithCancel(ctx)
		a.store = store.NewRedisStore(a.ctx)

		a.key = pkg.AnalyticsKey

		client = &a
	})

	return client, nil
}

func (a *Analytics) Run() error {
	atomic.SwapUint32(&a.flag, 1)
	go func() {
		for i := 0; i < a.workers; i++ {
			go func() {
				a.wg.Add(1)
				defer a.wg.Done()

				if err := a.runWorker(a.ctx); err != nil {
					_ = fmt.Errorf("authz store worker error: %s", err.Error())
					return
				}
			}()
		}
		a.wg.Wait()
	}()
	return nil
}

func (a *Analytics) Close() error {
	// - close all worker.
	a.cancel()

	atomic.SwapUint32(&a.flag, 0)
	close(a.pool)

	if err := a.store.Close(); err != nil {
		return err
	}
	return nil
}

func (a *Analytics) Record(info *analytics.RecordInfo) error {
	if atomic.LoadUint32(&a.flag) == 0 {
		return fmt.Errorf("record service is top")
	}

	info.ExpireAt = time.Now().Add(a.expire)

	bytes, err := info.Marshal()
	if err != nil {
		return err
	}

	a.pool <- bytes

	return nil
}

func (a *Analytics) runWorker(ctx context.Context) error {
	batch := make([][]byte, 0, a.batch)
	var lastSendTime = time.Now()

	for {
		var shouldSend bool

		select {
		case <-ctx.Done():
			if len(batch) != 0 {
				return a.doSend(batch)
			}
			return nil
		case data, ok := <-a.pool:

			// Do last work, general use ctx.Done, this func maybe does not run, but for safe.
			if !ok && len(batch) != 0 {
				return a.doSend(batch)
			}

			// must ensure batch's length smaller than a.batch
			batch = append(batch, data)

			// check at last
			shouldSend = uint64(len(batch)) == a.batch

		case <-time.After(a.flushInterval):
			if len(batch) > 0 {
				shouldSend = true
			}
		}

		// If result greater than, represent batch is expired.
		// drop expiry data.
		if time.Since(lastSendTime) > a.flushInterval {
			lastSendTime = time.Now()
			batch = batch[:0]
			shouldSend = false
		}

		if len(batch) > 0 && shouldSend {
			lastSendTime = time.Now()
			err := a.doSend(batch)
			batch = batch[:0]
			return err
		}
	}
}

func (a *Analytics) doSend(batch [][]byte) error {
	if err := a.store.Send(a.key, batch); err != nil {
		return err
	}
	return nil
}
