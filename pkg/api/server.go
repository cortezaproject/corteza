package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory/resputil"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/version"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
)

type (
	Server struct {
		name string

		log *zap.Logger

		httpOpt    *flags.HTTPOpt
		monitorOpt *flags.MonitorOpt

		endpoints []func(r chi.Router)
	}
)

func NewServer(log *zap.Logger) *Server {
	return &Server{
		endpoints: make([]func(r chi.Router), 0),
		log:       log.Named("http"),
	}
}

func (s *Server) Command(ctx context.Context, cmdName, prefix string, preRun func(context.Context) error) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   cmdName,
		Short: "Start HTTP Server with REST API",

		// Connect all the wires, prepare services, run watchers, bind endpoints
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if s.monitorOpt.Interval > 0 {
				go NewMonitor(int(s.monitorOpt.Interval / time.Second))
			}

			return preRun(ctx)
		},

		// Run the server
		Run: func(cmd *cobra.Command, args []string) {
			s.Serve(ctx)
		},
	}

	s.BindApiServerFlags(cmd, prefix)
	return
}

func (s *Server) BindApiServerFlags(cmd *cobra.Command, prefix string) {
	s.httpOpt = flags.HTTP(cmd, prefix)
	s.monitorOpt = flags.Monitor(cmd, prefix)
}

func (s *Server) MountRoutes(mm ...func(chi.Router)) {
	s.endpoints = append(s.endpoints, mm...)
}

func (s Server) Serve(ctx context.Context) {
	s.log.Info("Starting HTTP server with REST API", zap.String("address", s.httpOpt.Addr))

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: s.httpOpt.Pretty,
		Trace:  s.httpOpt.Tracing,
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

	router.Use(Base()...)

	if s.httpOpt.Logging {
		router.Use(Logging(s.log)...)
	}

	if s.httpOpt.EnableMetrics {
		router.Use(Middleware(s.httpOpt.MetricsServiceLabel))
	}

	router.Group(func(r chi.Router) {
		r.Use(
			auth.DefaultJwtHandler.Verifier(),
			auth.DefaultJwtHandler.Authenticator(),
		)

		for _, mountRoutes := range s.endpoints {
			mountRoutes(r)
		}
	})

	if s.httpOpt.EnableMetrics {
		Mount(router, s.httpOpt.MetricsUsername, s.httpOpt.MetricsPassword)
	}

	if s.httpOpt.EnableDebugRoute {
		Debug(router)
	}

	if s.httpOpt.EnableVersionRoute {
		router.Get("/version", version.HttpHandler)
	}

	go func() {
		err = http.Serve(listener, router)
	}()
	<-ctx.Done()

	if err == nil {
		err = ctx.Err()
		if err == context.Canceled {
			err = nil
		}
	}

	s.log.Info("HTTP server stopped", zap.Error(err))

	return
}
