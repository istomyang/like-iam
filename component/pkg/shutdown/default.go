package shutdown

import (
	"istomyang.github.com/like-iam/component/pkg/shutdown/receiver"
	"istomyang.github.com/like-iam/log"
)

// CreateDefaultShutdown create a default Shutdown with close func came from App.
func CreateDefaultShutdown(close SignalHandler) *Shutdown {
	sd := NewShutdown(func(err error) {
		log.Error(err.Error())
	})

	sd.RegisterReceiver(receiver.NewSysSignalReceiver())
	sd.RegisterHandler(close)

	return sd
}
