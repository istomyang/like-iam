package receiver

import (
	"os"
	"os/signal"
	"syscall"
)

// DoubleSysSignalReceiver receives system signal twice.
// If receives one, do some work, and do really close work before the second signal comes.
type DoubleSysSignalReceiver struct {
	midWork func(chan<- bool) error
	signs   []os.Signal
}

// NewDoubleSysSignalReceiver create DoubleSysSignalReceiver.
// parma midWork's reset channel means whether resetting receiver in true or false.
func NewDoubleSysSignalReceiver(midWork func(reset chan<- bool) error, signals ...os.Signal) *DoubleSysSignalReceiver {
	if signals == nil {
		signals = make([]os.Signal, 2)
		signals[0] = syscall.SIGINT
		signals[1] = syscall.SIGTERM
	}

	return &DoubleSysSignalReceiver{midWork: midWork, signs: signals}
}

func (r *DoubleSysSignalReceiver) Start(sig chan<- bool) error {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, r.signs...)

	<-ch
	if r.midWork != nil {
		reset := make(chan bool, 1)
		if err := r.midWork(reset); err != nil {
			return err
		}
		if <-reset {
			return r.Start(sig)
		}
	}
	<-ch

	sig <- true
	return nil
}

func (r *DoubleSysSignalReceiver) Before() error {
	return nil
}

func (r *DoubleSysSignalReceiver) After() error {
	return nil
}
