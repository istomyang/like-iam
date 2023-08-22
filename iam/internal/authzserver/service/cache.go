package service

import (
	"github.com/ory/ladon"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/component/pkg/interfaces"
)

type Cache interface {
	GetPolicy(k string) ([]ladon.Policy, error)
	GetSecret(k string) (*pb.SecretInfo, error)

	// Sync reloads data through store.Factory when sync signal is coming.
	Sync() error
	Clear() error

	interfaces.ComponentCommon
}
