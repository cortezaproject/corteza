package server

import (
	"context"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type (
	server struct {
		log            *zap.Logger
		httpOpt        options.HTTPServerOpt
		waitForOpt     options.WaitForOpt
		environmentOpt options.EnvironmentOpt
		endpoints      []func(r chi.Router)
	}
)

func New(log *zap.Logger, envOpt options.EnvironmentOpt, httpOpt options.HTTPServerOpt, waitForOpt options.WaitForOpt) *server {
	return &server{
		endpoints: make([]func(r chi.Router), 0),
		log:       log.Named("http"),

		environmentOpt: envOpt,
		httpOpt:        httpOpt,
		waitForOpt:     waitForOpt,
	}
}

func (s *server) MountRoutes(mm ...func(chi.Router)) {
	s.endpoints = append(s.endpoints, mm...)
}

func (s server) Serve(ctx context.Context) {
	s.log.Info(
		"starting HTTP server",

		zap.String("path-prefix", s.httpOpt.BaseUrl),
		zap.String("address", s.httpOpt.Addr),
	)

	listener, err := net.Listen("tcp", s.httpOpt.Addr)
	if err != nil {
		s.log.Error("cannot start server", zap.Error(err))
		return
	}

	router := chi.NewRouter()

	router.Route("/"+strings.TrimPrefix(s.httpOpt.BaseUrl, "/"), func(r chi.Router) {
		// Reports error to Sentry if enabled
		if s.httpOpt.EnablePanicReporting {
			r.Use(sentryMiddleware())
		}

		if s.httpOpt.EnableMetrics {
			// Metrics tracking middleware
			r.Use(metricsMiddleware(s.httpOpt.MetricsServiceLabel))
		}

		// Handle panic (sets 500 server error headers)
		r.Use(handlePanic)

		// Base middleware, CORS, RealIP, RequestID, context-logger
		r.Use(BaseMiddleware(s.environmentOpt.IsProduction(), s.log)...)

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

		if s.httpOpt.EnableMetrics {
			metricsMount(r, s.httpOpt.MetricsUsername, s.httpOpt.MetricsPassword)
		}

	})

	if s.httpOpt.BaseUrl != "/" {
		router.Handle("/", http.RedirectHandler(s.httpOpt.BaseUrl, http.StatusTemporaryRedirect))

	}

	s.bindMiscRoutes(router)

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

	s.log.Info("HTTP server stopped", zap.Error(err))
}

func (s server) bindMiscRoutes(r chi.Router) {
	if s.httpOpt.EnableDebugRoute {
		s.log.Debug("route debugger enabled: /__routes")
		r.Get("/__routes", debugRoutes(r))

		s.log.Debug("profiler enabled: /__profiler")
		r.Mount("/__profiler", middleware.Profiler())

		s.log.Debug("eventbus handlers debug enabled: /__eventbus")
		r.Get("/__eventbus", debugEventbus())

		s.log.Debug("corredor service debug enabled: /__corredor")
		r.Get("/__corredor", debugCorredor())
	}

	if s.httpOpt.EnableVersionRoute {
		var (
			dPath = "/version"
			sPath = path.Join(s.httpOpt.BaseUrl, dPath)
			v     = func(w http.ResponseWriter, r *http.Request) {
				api.Send(w, r, struct {
					BuildTime string `json:"buildTime"`
					Version   string `json:"version"`
				}{version.BuildTime, version.Version})
			}
		)

		r.Get(dPath, v)
		if dPath != sPath {
			r.Get(sPath, v)
		}
	}

	if s.httpOpt.EnableHealthcheckRoute {
		// default & sub path for healthcheck endpoint
		var (
			dPath = "/healthcheck"
			sPath = path.Join(s.httpOpt.BaseUrl, dPath)
			log   = s.log.With(zap.String("url", dPath))
		)

		r.Get(dPath, healthcheck.HttpHandler())

		if dPath != sPath {
			r.Get(sPath, healthcheck.HttpHandler())
			log = log.With(zap.String("url", sPath))
		}

		log.Info("healthcheck endpoint enabled")
	}
}
