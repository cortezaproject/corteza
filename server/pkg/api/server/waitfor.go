package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// WaitFor sets up a simple status page, delays execution and probes services
func (s server) WaitFor(ctx context.Context) {
	var (
		opt      = s.opts.WaitFor
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

	// Setup a simple HTTP server that will inform the impatient users
	listener, err := net.Listen("tcp", s.opts.HTTPServer.Addr)
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

func (s server) resolveService(service string) (addr string, u *url.URL, err error) {
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

func (s server) probeService(ctx context.Context, addr string) (err error) {
	if err != nil {
		return err
	}

	dialer := net.Dialer{}
	_, err = dialer.DialContext(ctx, "tcp", addr)
	return
}

func (s server) probeServiceURL(ctx context.Context, u *url.URL) error {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to assemble service request: %w", err)
	}

	rsp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("service URL request failed: %w", err)
	}

	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("service responded with unexpected status '%s'", rsp.Status)
}
