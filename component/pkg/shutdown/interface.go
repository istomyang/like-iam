package shutdown

// SignalHandler has error return because Shutdown is an independent module to handle shutdown business,
// If error occurs, then itself may handle this error such as report to logger server.
type SignalHandler func() error

type ErrorHandler func(err error)

// Receiver receives signal from system or others places, and tells Shutdown do the work such as closing resources.
type Receiver interface {
	// Start starts listening and receiving signal from external system, using chan type to notify Shutdown do close work.
	Start(signal chan<- bool) error

	// Before is a hook you can to some work.
	Before() error

	// After is a hook you can to some work such as close this receiver.
	After() error
}
