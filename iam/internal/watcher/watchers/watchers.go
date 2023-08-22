package watchers

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/robfig/cron/v3"
	"istomyang.github.com/like-iam/component-base/errors"
	"sync"
)

// Watcher defines impl of a watcher.
type Watcher interface {
	// Init uses context to close cron, and uses redisOptions to create redis sync mutex.
	Init(ctx context.Context, mut *redsync.Mutex, config interface{})
	// Schedules returns spec string from cron, like `@yearly`, `0 0 1 1 *` and so on.
	Schedules() string
	// Job is to connect cron lib.
	cron.Job
}

var (
	ErrDuplicate = errors.New("register name has already registered.")
	ErrConfig    = errors.New("watcher's config is invalid.")
)

var (
	watchers = make(map[string]Watcher)
	mut      = new(sync.Mutex)
)

// Register will be called in impl init function, which don't need to validate.
func Register(name string, watcher Watcher) {
	mut.Lock()
	defer mut.Unlock()

	if _, ex := watchers[name]; ex {
		panic(ErrDuplicate)
	}
	watchers[name] = watcher
}

func Find(name string) Watcher {
	mut.Lock()
	defer mut.Unlock()

	return watchers[name]
}

func List() []Watcher {
	mut.Lock()
	defer mut.Unlock()

	var l = make([]Watcher, len(watchers))
	for _, watcher := range watchers {
		l = append(l, watcher)
	}
	return l
}

func ListMap() map[string]Watcher {
	mut.Lock()
	defer mut.Unlock()
	return watchers
}
