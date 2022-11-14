package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
)

// WaitFor sets up a simple status page, delays execution and probes services
func (s server) WaitFor(ctx context.Context) {
	var (
		err error
	)

	// Set up a simple HTTP server that will inform the impatient users
	listener, err := net.Listen("tcp", s.httpOpt.Addr)
	if err != nil {
		s.log.Error("cannot start server", zap.Error(err))
		os.Exit(1)
	}
	defer listener.Close()
	go func() {
		router := chi.NewRouter()
		router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Write([]byte("waiting for services..."))
		})
		_ = http.Serve(listener, router)
	}()
	<-ctx.Done()
}
