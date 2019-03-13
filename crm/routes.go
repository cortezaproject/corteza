package service

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/crusttech/crust/crm/rest"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/routes"
	"github.com/crusttech/crust/internal/version"
)

func Routes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()
	MountRoutes(ctx, r)
	routes.Print(r)
	MountSystemRoutes(r, flags.http)
	return r
}

func MountRoutes(ctx context.Context, r chi.Router) {
	r.Use(handleCORS)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	// Only protect application routes with JWT
	r.Group(func(r chi.Router) {
		r.Use(jwtVerifier, jwtAuthenticator)
		mountRoutes(r, flags.http, rest.MountRoutes())
	})
}

func MountSystemRoutes(r chi.Router, opts *config.HTTP) {
	metrics.MountRoutes(r, opts)
	r.Mount("/debug", middleware.Profiler())
	r.Get("/version", version.HttpHandler)
}

func mountRoutes(r chi.Router, opts *config.HTTP, mounts ...func(r chi.Router)) {
	if opts.Logging {
		r.Use(middleware.Logger)
	}
	if opts.Metrics {
		r.Use(metrics.Middleware("crm"))
	}

	for _, mount := range mounts {
		mount(r)
	}
}

// Sets up default CORS rules to use as a middleware
func handleCORS(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler(next)
}
