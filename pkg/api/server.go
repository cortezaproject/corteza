package api

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/go-chi/chi/middleware"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	server struct {
		log        *zap.Logger
		httpOpt    options.HTTPServerOpt
		waitForOpt options.WaitForOpt
		endpoints  []func(r chi.Router)
	}
)

func New(log *zap.Logger, httpOpt options.HTTPServerOpt, waitForOpt options.WaitForOpt) *server {
	return &server{
		endpoints:  make([]func(r chi.Router), 0),
		log:        log.Named("http"),
		httpOpt:    httpOpt,
		waitForOpt: waitForOpt,
	}
}

func (s *server) MountRoutes(mm ...func(chi.Router)) {
	s.endpoints = append(s.endpoints, mm...)
}

func (s server) Serve(ctx context.Context) {
	s.log.Info("Starting HTTP server with REST API", zap.String("address", s.httpOpt.Addr))

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Trace: s.httpOpt.Tracing,
		Logger: func(err error) {
			// @todo: error logging
		},
	})

	listener, err := net.Listen("tcp", s.httpOpt.Addr)
	if err != nil {
		s.log.Error("Can not start server", zap.Error(err))
		return
	}

	router := chi.NewRouter()

	// Base middleware, CORS, RealIP, RequestID, context-logger
	router.Use(BaseMiddleware(s.log)...)

	router.Group(func(r chi.Router) {
		s.bindMiscRoutes(r)
	})

	router.Group(func(r chi.Router) {
		// Logging request if enabled
		if s.httpOpt.LogRequest {
			r.Use(LogRequest)
		}

		// Logging response if enabled
		if s.httpOpt.LogResponse {
			r.Use(LogResponse)
		}

		r.Use(
			auth.DefaultJwtHandler.HttpVerifier(),
			auth.DefaultJwtHandler.HttpAuthenticator(),
		)

		for _, mountRoutes := range s.endpoints {
			mountRoutes(r)
		}
	})

	go func() {
		srv := http.Server{
			Handler: router,
			BaseContext: func(listener net.Listener) context.Context {
				return ctx
			},
		}
		err = srv.Serve(listener)
	}()
	<-ctx.Done()

	if err == nil {
		err = ctx.Err()
		if err == context.Canceled {
			err = nil
		}
	}

	s.log.Info("Server stopped", zap.Error(err))
}

func (s server) bindMiscRoutes(router chi.Router) {
	if s.httpOpt.EnableMetrics {
		metricsMount(router, s.httpOpt.MetricsUsername, s.httpOpt.MetricsPassword)
	}

	// Metrics tracking middleware
	if s.httpOpt.EnableMetrics {
		router.Use(metricsMiddleware(s.httpOpt.MetricsServiceLabel))
	}

	// Handle panic (sets 500 server error headers)
	router.Use(handlePanic)

	// Reports error to Sentry if enabled
	if s.httpOpt.EnablePanicReporting {
		router.Use(sentryMiddleware())
	}

	if s.httpOpt.EnableDebugRoute {
		s.log.Debug("profiler: /__profiler")
		router.Mount("/__profiler", middleware.Profiler())

		s.log.Debug("list of routes: /__routes")
		router.Get("/__routes", debugRoutes(router))

		s.log.Debug("eventbus handlers: /__eventbus")
		router.Get("/__eventbus", debugEventbus())

		s.log.Debug("corredor service: /__corredor")
		router.Get("/__corredor", debugCorredor())
	}

	if s.httpOpt.EnableVersionRoute {
		router.Get("/version", version.HttpHandler)
	}

	if s.httpOpt.EnableHealthcheckRoute {
		router.Get("/healthcheck", healthcheck.HttpHandler())
	}
}
