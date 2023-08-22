package store

var client Factory

type Factory interface {
	User() UserStore
	Secret() SecretStore
	Policy() PolicyStore

	Run() error
	Close() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
