package store

var client Factory

type Factory interface {
	Secrets() SecretStore
	Policies() PolicyStore

	Run() error
	Close() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
