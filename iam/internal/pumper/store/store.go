package store

import "istomyang.github.com/like-iam/component/pkg/interfaces"

var client Factory

type Factory interface {
	Pop(key string) ([][]byte, error)

	interfaces.ComponentCommon
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
