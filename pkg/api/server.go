package api

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory/resputil"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/version"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

type (
	Server struct {
		name string

		log *zap.Logger

		httpOpt    *options.HTTPOpt
		monitorOpt *options.MonitorOpt

		endpoints []func(r chi.Router)
	}
)

var (
	Monolith = false
	BaseURL  = "/"
)

func NewServer(log *zap.Logger) *Server {
	return &Server{
		endpoints: make([]func(r chi.Router), 0),
		log:       log.Named("http"),
	}
}

func (s *Server) Command(ctx context.Context, cmdName, prefix string, preRun func(context.Context) error) (cmd *cobra.Command) {
	s.httpOpt = options.HTTP(prefix)
	s.monitorOpt = options.Monitor(prefix)

	cmd = &cobra.Command{
		Use:   cmdName,
		Short: "Start HTTP Server with REST API",

		// Connect all the wires, prepare services, run watchers, bind endpoints
		PreRunE: func(cmd *cobra.Command, args []string) error {
			s.waitFor(ctx, options.WaitFor(prefix))

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

	return
}

func (s *Server) MountRoutes(mm ...func(chi.Router)) {
	s.endpoints = append(s.endpoints, mm...)
}

func (s Server) Serve(ctx context.Context) {
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
	router.Use(Base(s.log)...)

	// Logging request if enabled
	if s.httpOpt.LogRequest {
		router.Use(LogRequest)
	}

	// Logging response if enabled
	if s.httpOpt.LogResponse {
		router.Use(LogResponse)
	}

	// Handle panic (sets 500 Server error headers)
	router.Use(HandlePanic)

	// Reports error to Sentry if enabled
	if s.httpOpt.EnablePanicReporting {
		router.Use(Sentry())
	}

	// Metrics tracking middleware
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

// waitFor sets up a simple status page, delays execution and probes services
func (s Server) waitFor(ctx context.Context, opt *options.WaitForOpt) {
	var (
		services = opt.GetServices()
	)

	if len(services) == 0 && opt.Delay == 0 {
		// Nothing to do here..
		return
	}

	var (
		log         = s.log.Named("wait-for")
		depChan     = make(chan struct{})
		wg          sync.WaitGroup
		serviceAddr string
		serviceURL  *url.URL
		err         error
	)

	// Setup a simple HTTP server that will inform the impatent users
	listener, err := net.Listen("tcp", s.httpOpt.Addr)
	if err != nil {
		s.log.Error("Can not start server", zap.Error(err))
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

	if opt.Delay > 0 {
		s.log.Info("delaying", zap.Duration("delay", opt.Delay))

		// First delay execution
		select {
		case <-ctx.Done():
			log.Debug("canceled")
			return
		case <-time.After(opt.Delay):
			// all good...
		}
	}

	if len(services) == 0 {
		return
	}

	log.Info("waiting for services", zap.Strings("services", services))
	// Probe services
	wg.Add(len(services))
	go func() {

		for _, service := range services {
			slog := log.With(zap.String("service", service))

			go func(ctx context.Context, service string) {
				defer wg.Done()

				if serviceAddr, serviceURL, err = s.resolveService(service); err != nil {
					log.Error("could not resolve service", zap.Error(err))
				}

				for {
					ctx, cancelFn := context.WithTimeout(ctx, opt.ServicesProbeTimeout)
					defer cancelFn()

					if serviceURL == nil {
						if err = s.probeService(ctx, serviceAddr); err != nil {
							slog.Warn("service probe failed", zap.Error(err))
							time.Sleep(opt.ServicesProbeInterval)
							continue
						}
					} else {
						if err = s.probeServiceURL(ctx, serviceURL); err != nil {
							slog.Warn("service URL probe failed", zap.Error(err))
							time.Sleep(opt.ServicesProbeInterval)
							continue
						}
					}

					slog.Debug("service ready")
					return
				}
			}(ctx, service)
		}
		wg.Wait()
		close(depChan)
	}()

	select {
	case <-ctx.Done():
		log.Debug("canceled")
		return
	case <-depChan: // services are ready
		log.Debug("all services ready")
		return
	case <-time.After(opt.ServicesTimeout):
		log.Debug("services not ready")
		os.Exit(1)
	}
}

func (s Server) resolveService(service string) (addr string, u *url.URL, err error) {
	addr = service

	if strings.Contains(addr, "://") {
		// Is service an URL?
		u, err = url.Parse(addr)
		if err != nil {
			return
		}

		addr = u.Host

		if u.Port() == "" {
			if u.Scheme == "https" {
				addr += ":443"
			}
		}
	}

	// Default to port 80
	if !strings.Contains(addr, ":") {
		addr += ":80"
	}

	return
}

func (s Server) probeService(ctx context.Context, addr string) (err error) {
	if err != nil {
		return err
	}

	dialer := net.Dialer{}
	_, err = dialer.DialContext(ctx, "tcp", addr)
	return
}

func (s Server) probeServiceURL(ctx context.Context, u *url.URL) error {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to assemble service request")
	}

	rsp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "service URL request failed")
	}

	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.Errorf("service responded with unexpected status '%s'", rsp.Status)
}
