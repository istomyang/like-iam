package cache

import (
	"context"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/iam/internal/apiserver/service"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
)

type Cache struct {
	svc service.Service
	pb.UnimplementedCacheServer
}

func NewCache(store store.Factory) *Cache {
	return &Cache{svc: service.NewService(store)}
}

func (c *Cache) ListSecrets(context.Context, *pb.ListRequest) (*pb.ListSecretsReply, error) {
	return nil, nil
}

func (c *Cache) ListPolicies(context.Context, *pb.ListRequest) (*pb.ListPoliciesReply, error) {
	return nil, nil
}
