package middleware

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/version"
)

func Mount(ctx context.Context, r chi.Router, opts *config.HTTP) {
	r.Use(handleCORS)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	if opts.Logging {
		r.Use(middleware.Logger)
	}
	if opts.Metrics {
		r.Use(metrics.Middleware("crust"))
	}
}

func MountSystemRoutes(ctx context.Context, r chi.Router, opts *config.HTTP) {
	metrics.MountRoutes(r, opts)
	r.Mount("/debug", middleware.Profiler())
	r.Get("/version", version.HttpHandler)
}
