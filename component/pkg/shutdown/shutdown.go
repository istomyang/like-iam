package shutdown

import (
	"os"
	"sync"
)

type Shutdown struct {
	receivers []Receiver
	handlers  []SignalHandler
	eh        ErrorHandler

	// ch should become a 1 capacity channel.
	ch chan bool
}

// NewShutdown return Shutdown with ErrorHandler which call when error occurs.
// If ErrorHandler is nil, Shutdown will skip error reporting.
func NewShutdown(eh ErrorHandler) *Shutdown {
	return &Shutdown{
		receivers: make([]Receiver, 5),
		handlers:  make([]SignalHandler, 5),
		eh:        eh,
		ch:        make(chan bool, 1),
	}
}

func (s *Shutdown) RegisterReceiver(receiver Receiver) {
	s.receivers = append(s.receivers, receiver)
}

func (s *Shutdown) RegisterHandler(handler SignalHandler) {
	s.handlers = append(s.handlers, handler)
}

func (s *Shutdown) Run() {
	go s.start()
}

func (s *Shutdown) start() {

	s.do(func(receiver Receiver) {
		if err := receiver.Start(s.ch); err != nil {
			s.error(err)
		}
	})

	<-s.ch

	s.do(func(receiver Receiver) {
		if err := receiver.Before(); err != nil {
			s.error(err)
		}
	})

	var wg sync.WaitGroup

	for _, handler := range s.handlers {
		wg.Add(1)
		go func(signalHandler SignalHandler) {
			defer wg.Done()

			if err := signalHandler(); err != nil {
				s.error(err)
			}
		}(handler)
	}

	wg.Wait()

	s.do(func(receiver Receiver) {
		if err := receiver.After(); err != nil {
			s.error(err)
		}
	})

	os.Exit(0)
}

// error prevents err and eh is nil.
func (s *Shutdown) error(err error) {
	if s.eh != nil && err != nil {
		s.eh(err)
	}
}

func (s *Shutdown) do(f func(receiver Receiver)) {
	for _, receiver := range s.receivers {
		go func(receiver Receiver) {
			f(receiver)
		}(receiver)
	}
}
