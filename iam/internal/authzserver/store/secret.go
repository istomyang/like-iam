package store

import pb "istomyang.github.com/like-iam/api/proto/v1"

// SecretStore lists data from apiserver server.
type SecretStore interface {
	List() (map[string]*pb.SecretInfo, error)
}
