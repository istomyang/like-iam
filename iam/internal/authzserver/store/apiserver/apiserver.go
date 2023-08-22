package apiserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/iam/internal/authzserver/store"
	"sync"
)

type datastore struct {
	pb   pb.CacheClient
	conn *grpc.ClientConn

	address string
	cert    string

	ctx context.Context
}

var (
	factory store.Factory
	once    sync.Once
)

func GetApiServerStoreOr(ctx context.Context, address, cert string) (store.Factory, error) {
	if factory == nil && address == "" {
		return nil, fmt.Errorf("must assign params to create factory")
	}

	once.Do(func() {
		factory = &datastore{address: address, cert: cert, ctx: ctx}
	})

	return factory, nil
}

func (s *datastore) Secrets() store.SecretStore {
	return newSecret(s.ctx, s.pb)
}

func (s *datastore) Policies() store.PolicyStore {
	return newPolicy(s.ctx, s.pb)
}

func (s *datastore) Run() error {
	// allow empty.
	credential, _ := credentials.NewClientTLSFromFile(s.cert, "")

	conn, err := grpc.Dial(s.address, grpc.WithTransportCredentials(credential), grpc.WithBlock())
	if err != nil {
		return err
	}

	s.conn = conn
	s.pb = pb.NewCacheClient(conn)

	return nil
}

func (s *datastore) Close() error {
	return s.conn.Close()
}
