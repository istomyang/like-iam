package server

import (
	"context"
	stdError "errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
	"istomyang.github.com/like-iam/component-base/web"
	"istomyang.github.com/like-iam/component/pkg/middleware"
	"istomyang.github.com/like-iam/component/pkg/options"
	"istomyang.github.com/like-iam/log"
	"net/http"
	"time"
)

type Option func(server *GeneralApiServer)

type GeneralApiServer struct {
	*options.ServerOpts
	*options.InsecureServerOpts
	*options.SecureServerOpts
	*options.FeatureOptions

	// gin.Engine satisfies http.Handler interface.
	*gin.Engine

	insecure, secure *http.Server
}

// Install initializes config for server.
func (s *GeneralApiServer) Install() (errors []error) {

	errors = append(errors, s.validate()...)

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	for _, m := range middleware.MustMiddlewares {
		s.Engine.Use(m)
	}
	for _, m := range s.PresetMiddlewares {
		s.Engine.Use(middleware.PresetMiddlewares[m])
	}

	if s.Healthz {
		s.GET("/healthz", func(c *gin.Context) {
			web.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}

	if s.Metrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	if s.Profile {
		pprof.Register(s.Engine)
	}

	gin.SetMode(s.Mode)

	// TODO: Version tag

	return
}

// Run a server with a context which can be canceled by father context.
// Must be last to run.
func (s *GeneralApiServer) Run() error {

	s.insecure = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.InsecureServerOpts.Address, s.InsecureServerOpts.Port),
		Handler: s,
	}
	s.secure = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.SecureServerOpts.Address, s.SecureServerOpts.Port),
		Handler: s,
	}

	eg, c := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s:%d",
			s.InsecureServerOpts.Address, s.InsecureServerOpts.Port)

		if err := s.insecure.ListenAndServe(); err != nil && !stdError.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
			return err
		}

		log.Infof("Server on %s:%d stopped", s.InsecureServerOpts.Address, s.InsecureServerOpts.Port)
		return nil
	})

	eg.Go(func() (err error) {
		log.Infof("Start to listening the incoming requests on http address: %s:%d",
			s.SecureServerOpts.Address, s.SecureServerOpts.Port)

		if err := s.secure.ListenAndServeTLS(
			s.SecureServerOpts.Tls.CertFile,
			s.SecureServerOpts.Tls.KeyFile); err != nil && !stdError.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
			return err
		}

		log.Infof("Server on %s:%d stopped", s.SecureServerOpts.Address, s.SecureServerOpts.Port)
		return nil
	})

	if s.Healthz {
		host := s.InsecureServerOpts.Address
		if host == "0.0.0.0" {
			host = "127.0.0.1"
		}
		url := fmt.Sprintf("%s:%d", host, s.InsecureServerOpts.Port)
		if err := web.Ping(c, url, time.Second*10); err != nil {
			return err
		}
	}

	s.Release()

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

// Close mainly close this server.
func (s *GeneralApiServer) Close() error {
	c, cancel := context.WithTimeout(context.Background(), s.ShutdownTime)
	defer cancel()

	if err := s.secure.Shutdown(c); err != nil {
		log.Warnf("http api server shutdown got error: %s", err.Error())
		return err
	}
	if err := s.insecure.Shutdown(c); err != nil {
		log.Warnf("https api server shutdown got error: %s", err.Error())
		return err
	}
	return nil
}

// Release releases heap objects not use anymore.
func (s *GeneralApiServer) Release() {
	s.ServerOpts = nil
	s.InsecureServerOpts = nil
	s.SecureServerOpts = nil
}

func (s *GeneralApiServer) validate() (errors []error) {
	errors = append(errors, s.ServerOpts.Validate()...)
	errors = append(errors, s.InsecureServerOpts.Validate()...)
	errors = append(errors, s.SecureServerOpts.Validate()...)
	return
}
