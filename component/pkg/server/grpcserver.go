package server

import (
	"fmt"
	"google.golang.org/grpc"
	"istomyang.github.com/like-iam/log"
	"net"
)

type GeneralGRPCServer struct {
	svr     *grpc.Server
	address string
}

func NewGeneralGRpcServer(address string, opt ...grpc.ServerOption) *GeneralGRPCServer {
	return &GeneralGRPCServer{svr: grpc.NewServer(opt...), address: address}
}

func (s *GeneralGRPCServer) Install(register func(*grpc.Server) error) error {
	return register(s.svr)
}

func (s *GeneralGRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Infof("server listening at %v", lis.Addr())
	go func() {
		if err := s.svr.Serve(lis); err != nil {
			log.Fatalf("fail to serve apiserver: &v", err.Error())
		}
	}()

	return nil
}

func (s *GeneralGRPCServer) Close() error {
	log.Infof("GRPC server on %s stopped", s.address)
	s.svr.GracefulStop()
	return nil
}
