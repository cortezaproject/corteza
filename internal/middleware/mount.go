package middleware

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/version"
)

func Mount(ctx context.Context, r chi.Router, opts *config.HTTP) {
	r.Use(handleCORS)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(ContextLogger(logger.ContextValue(ctx)))

	if opts.Logging {
		r.Use(LogRequest)
		r.Use(LogResponse)
	}

	if opts.Metrics {
		r.Use(metrics.Middleware("crust"))
	}
}

func MountSystemRoutes(ctx context.Context, r chi.Router, opts *config.HTTP) {
	metrics.MountRoutes(r, opts)
	r.Mount("/debug", middleware.Profiler())
	r.Get("/version", version.HttpHandler)

	r.Get("/routes", func(w http.ResponseWriter, req *http.Request) {
		var printRoutes func(chi.Routes, string)

		printRoutes = func(r chi.Routes, pfix string) {
			routes := r.Routes()
			for _, route := range routes {
				if route.SubRoutes != nil && len(route.SubRoutes.Routes()) > 0 {
					printRoutes(route.SubRoutes, pfix+route.Pattern[:len(route.Pattern)-2])
				} else {
					for method, fn := range route.Handlers {
						fmt.Fprintf(w, "%-8s %-80s -> %s\n", method, pfix+route.Pattern, runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
					}
				}
			}
		}

		printRoutes(r, "")
	})
}
