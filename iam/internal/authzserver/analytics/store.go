package analytics

import (
	"istomyang.github.com/like-iam/component/pkg/interfaces"
	"time"
)

type Store interface {
	WithExpire(duration time.Duration) // Set value's expiration.

	Send(k string, batch [][]byte) error

	interfaces.ComponentCommon
}
