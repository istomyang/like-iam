package service

import "istomyang.github.com/like-iam/component/pkg/interfaces"

type Subscribe interface {
	interfaces.ComponentCommon

	// OnReceive use sig channel to reload data through gRPC.
	OnReceive(message chan<- bool) error
}
