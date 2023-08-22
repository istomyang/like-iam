package apiserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	v1 "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/server"
	"istomyang.github.com/like-iam/component/pkg/shutdown"
	"istomyang.github.com/like-iam/iam/internal/apiserver/auth"
	"istomyang.github.com/like-iam/iam/internal/apiserver/controller/v1/cache"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store"
	"istomyang.github.com/like-iam/iam/internal/apiserver/store/mysql"
	"istomyang.github.com/like-iam/log"
)

type apiServer struct {
	svr   *server.GeneralApiServer
	redis *conn.RedisClient
	grpc  *server.GeneralGRPCServer

	ctx    context.Context
	cancel context.CancelFunc

	shutdown *shutdown.Shutdown
}

func newApiServer(options *Options) *apiServer {
	s := &apiServer{}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	initSingletonStore(options)

	// in create stage.
	auth.GetJwtSchemeOr(options.jwtOptions)

	s.svr = createSvr(options)
	s.redis = createRedis(options)
	s.grpc = createGRpc(options)
	s.shutdown = shutdown.CreateDefaultShutdown(s.close)

	return s
}

func (s *apiServer) run() error {
	s.shutdown.Run()

	if err := s.redis.Run(); err != nil {
		return err
	}
	if err := s.grpc.Run(); err != nil {
		return err
	}

	// http server must at last and block main thread.
	if err := s.svr.Run(); err != nil {
		return err
	}

	return nil
}

// close does some clean work when app exits.
func (s *apiServer) close() error {
	s.cancel()

	if err := s.svr.Close(); err != nil {
		return err
	}
	if err := s.redis.Close(); err != nil {
		return err
	}
	if err := s.grpc.Close(); err != nil {
		return err
	}

	if err := store.Client().Close(); err != nil {
		return err
	}

	return nil
}

func initSingletonStore(options *Options) {
	factory, err := mysql.GetMySQLFactoryOr(options.mysqlOptions)
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
	store.SetClient(factory)
}

func createSvr(options *Options) *server.GeneralApiServer {
	engine := gin.Default()

	svr := &server.GeneralApiServer{
		ServerOpts:         options.httpSvrOptions,
		InsecureServerOpts: options.insecureSvrOptions,
		SecureServerOpts:   options.secureSvrOptions,
		FeatureOptions:     options.featureOptions,
		Engine:             engine,
	}
	svr.Install()

	installMiddlewares(engine)
	installRouter(engine)

	return svr
}

func createRedis(options *Options) *conn.RedisClient {
	return conn.NewRedisClientOr(options.redisOptions)
}

func createGRpc(options *Options) *server.GeneralGRPCServer {
	addr := fmt.Sprintf("%s:%d", options.gRPCOptions.Address, options.gRPCOptions.Port)
	file, err := credentials.NewServerTLSFromFile(options.secureSvrOptions.Tls.CertFile, options.secureSvrOptions.Tls.KeyFile)
	if err != nil {
		log.Error(err.Error())
		panic(err.Error())
	}
	svr := server.NewGeneralGRpcServer(addr, grpc.MaxRecvMsgSize(options.gRPCOptions.MaxMsgSize), grpc.Creds(file))

	if err := svr.Install(func(g *grpc.Server) error {

		v1.RegisterCacheServer(g, cache.NewCache(store.Client()))

		// TODO:
		// To see server status info.
		reflection.Register(g)

		return nil
	}); err != nil {
		log.Error(err.Error())
		panic(err.Error())
	}

	return svr
}
