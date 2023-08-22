package subscribe

import (
	"context"
	"github.com/go-redis/redis/v8"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service"
	"istomyang.github.com/like-iam/iam/internal/pkg"
	"istomyang.github.com/like-iam/log"
	"time"
)

type redisSub struct {
	reload chan<- bool
	pubSub *redis.PubSub
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisSubClient(ctx context.Context) (service.Subscribe, error) {
	r := &redisSub{}
	r.ctx, r.cancel = context.WithCancel(ctx)
	return r, nil
}

func (r *redisSub) Run() error {
	go r.listen()
	return nil
}

func (r *redisSub) listen() {
	client := conn.GetRedisClient()
	pubSub := client.UniversalClient().Subscribe(r.ctx, pkg.PubSubChannel)
	pubSub.Channel(
		redis.WithChannelSize(1),
		redis.WithChannelHealthCheckInterval(time.Second*5),
		redis.WithChannelSendTimeout(time.Second*5)) // maybe consumer will fail if message left for x seconds.
	r.pubSub = pubSub

	// every x seconds handle the message.
	// note: this func has no way to release, which cause leak, but it's ok in this.
	ticker := time.Tick(time.Second * 2)

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker:
			message := <-pubSub.Channel()
			switch message.Payload {
			case pkg.MessageSecret, pkg.MessagePolicy:
				// TODO: more research, if it has UUID in Payload for debug pub and sub system.
				r.reload <- true
			default:
				log.Warnf("redis pubSub receive wrong Payload, which body: %v", message)
				time.Sleep(time.Second * 10)
			}
		}
	}
}

func (r *redisSub) Close() error {
	if err := r.pubSub.Close(); err != nil {
		return err
	}
	return nil
}

func (r *redisSub) OnReceive(message chan<- bool) error {
	r.reload = message
	return nil
}

var _ service.Subscribe = &redisSub{}
