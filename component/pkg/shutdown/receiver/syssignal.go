package receiver

import (
	"os"
	"os/signal"
	"syscall"
)

// SysSignalReceiver receives system signal once.
// If receives one, then tell can close.
type SysSignalReceiver struct {
	signs []os.Signal
}

func NewSysSignalReceiver(signals ...os.Signal) *SysSignalReceiver {
	if signals == nil {
		signals = make([]os.Signal, 2)
		signals[0] = syscall.SIGINT
		signals[1] = syscall.SIGTERM
	}
	return &SysSignalReceiver{signs: signals}
}

func (r *SysSignalReceiver) Start(sig chan<- bool) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, r.signs...)
	<-ch

	sig <- true
	return nil
}

func (r *SysSignalReceiver) Before() error {
	return nil
}

func (r *SysSignalReceiver) After() error {
	return nil
}
