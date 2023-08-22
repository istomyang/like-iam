package authzserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"istomyang.github.com/like-iam/component/pkg/conn"
	"istomyang.github.com/like-iam/component/pkg/server"
	"istomyang.github.com/like-iam/component/pkg/shutdown"
	"istomyang.github.com/like-iam/iam/internal/authzserver/analytics"
	"istomyang.github.com/like-iam/iam/internal/authzserver/authorization"
	"istomyang.github.com/like-iam/iam/internal/authzserver/service"
	"istomyang.github.com/like-iam/iam/internal/authzserver/store"
	"istomyang.github.com/like-iam/iam/internal/authzserver/store/apiserver"
	"istomyang.github.com/like-iam/log"
)

type authzServer struct {
	svr *server.GeneralApiServer

	ctx    context.Context
	cancel context.CancelFunc

	shutdown *shutdown.Shutdown
}

func newAuthzServer(options *Options) *authzServer {
	s := authzServer{}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	_ = conn.NewRedisClientOr(options.redisOptions)

	initSingletonStore(s.ctx, options)

	if _, err := service.NewService(s.ctx); err != nil {
		log.Fatal(err.Error())
		panic(err)
		return nil
	}
	if _, err := analytics.NewAnalytics(s.ctx, options.analyticsOptions); err != nil {
		log.Fatal(err.Error())
		panic(err)
		return nil
	}
	if _, err := authorization.NewAuthorizator(s.ctx); err != nil {
		log.Fatal(err.Error())
		panic(err)
		return nil
	}

	s.svr = createSvr(options)

	s.shutdown = shutdown.CreateDefaultShutdown(s.close)

	return &s
}

func (s *authzServer) run() error {

	s.shutdown.Run()

	if err := conn.GetRedisClient().Run(); err != nil {
		return err
	}

	if err := store.Client().Run(); err != nil {
		return err
	}

	// Here will run once to test ok, must ensure all components are ready.
	if err := service.GetService().Run(); err != nil {
		return err
	}

	if err := analytics.GetAnalytics().Run(); err != nil {
		return err
	}

	// http server must at last and block main thread.
	if err := s.svr.Run(); err != nil {
		return err
	}

	return nil
}

func (s *authzServer) close() error {
	s.cancel()

	if err := s.svr.Close(); err != nil {
		return err
	}

	if err := analytics.GetAnalytics().Close(); err != nil {
		return err
	}
	if err := service.GetService().Close(); err != nil {
		return err
	}
	if err := store.Client().Close(); err != nil {
		return err
	}
	if err := conn.GetRedisClient().Close(); err != nil {
		return err
	}

	return nil
}

func initSingletonStore(ctx context.Context, options *Options) {
	addr := fmt.Sprintf("%s:%d", options.gRPCOptions.Address, options.gRPCOptions.Port)
	factory, err := apiserver.GetApiServerStoreOr(ctx, addr, options.clientCA)
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
