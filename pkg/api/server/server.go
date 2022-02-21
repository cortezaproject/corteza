package server

import (
	"context"
	"net"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	server struct {
		log       *zap.Logger
		opts      *options.Options
		endpoints []func(r chi.Router)

		demux *demux
	}
)

const (
	waiting uint32 = iota
	active
	shutdown
)

// New initializes new HTTP server with special powers
// Server is started as early as possible and with a special request handler
// that demultiplexes request to one of the configured routers according to the server state.
//
// Waiting state
// This is initial state that with some simple route handlers:
//  - /version
//  - /healthcheck
//  - /healthcheck

func New(log *zap.Logger, opts *options.Options) *server {
	s := &server{
		endpoints: make([]func(r chi.Router), 0),
		log:       log.Named("http"),

		opts: opts,
	}

	s.demux = Demux(waiting, waitingRoutes(s.log.Named("waiting"), s.opts.HTTPServer))
	s.demux.Router(shutdown, shutdownRoutes())

	return s
}

// Activate reconfigures server to use active routes
func (s *server) Activate(mm ...func(chi.Router)) {
	s.demux.Router(active, activeRoutes(s.log, mm, s.opts))

	s.log.Debug("entering active state")
	s.demux.State(active)
}

// Shutdown reconfigures server to use shutdown routes
func (s *server) Shutdown() {
	s.log.Debug("entering shutdown state")
	s.demux.State(shutdown)
}

func (s server) Serve(ctx context.Context) {
	s.log.Info(
		"starting HTTP server",

		zap.String("path-prefix", s.opts.HTTPServer.BaseUrl),
		zap.String("address", s.opts.HTTPServer.Addr),
	)

	listener, err := net.Listen("tcp", s.opts.HTTPServer.Addr)
	if err != nil {
		s.log.Error("cannot start server", zap.Error(err))
		return
	}

	go func() {
		srv := http.Server{
			Handler: s.demux,

			// use root context as server's base context and as a basis for
			// context for all requests
			// this enables us to send cancellation down to every request
			BaseContext: func(listener net.Listener) context.Context { return ctx },
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
