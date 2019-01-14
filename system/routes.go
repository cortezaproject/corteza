package service

import (
	"fmt"
	"reflect"
	"runtime"

	"net/http"

	"github.com/99designs/basicauth-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/metrics"
	"github.com/crusttech/crust/internal/version"
	"github.com/crusttech/crust/system/rest"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(handleCORS)

	// Only protect application routes with JWT
	r.Group(func(r chi.Router) {
		r.Use(jwtVerifier, jwtAuthenticator)
		mountRoutes(r, flags.http, rest.MountRoutes(flags.oidc, flags.social, jwtEncoder))
	})

	printRoutes(r, flags.http)
	mountSystemRoutes(r, flags.http)
	return r
}

func mountRoutes(r chi.Router, opts *config.HTTP, mounts ...func(r chi.Router)) {
	if opts.Logging {
		r.Use(middleware.Logger)
	}
	if opts.Metrics {
		r.Use(metrics.Middleware("auth"))
	}

	for _, mount := range mounts {
		mount(r)
	}
}

func mountSystemRoutes(r chi.Router, opts *config.HTTP) {
	if opts.Metrics {
		r.Group(func(r chi.Router) {
			r.Use(basicauth.New("Metrics", map[string][]string{
				opts.MetricsUsername: {opts.MetricsPassword},
			}))
			r.Handle("/metrics", metrics.Handler())
		})
	}
	r.Mount("/debug", middleware.Profiler())
	r.Get("/version", version.HttpHandler)
}

func printRoutes(r chi.Router, opts *config.HTTP) {
	var printRoutes func(chi.Routes, string, string)
	printRoutes = func(r chi.Routes, indent string, prefix string) {
		routes := r.Routes()
		for _, route := range routes {
			if route.SubRoutes != nil && len(route.SubRoutes.Routes()) > 0 {
				fmt.Printf(indent+"%s - with %d handlers, %d subroutes\n", route.Pattern, len(route.Handlers), len(route.SubRoutes.Routes()))
				printRoutes(route.SubRoutes, indent+"\t", prefix+route.Pattern[:len(route.Pattern)-2])
			} else {
				for key, fn := range route.Handlers {
					fmt.Printf("%s%s\t%s -> %s\n", indent, key, prefix+route.Pattern, runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
				}
			}
		}
	}
	printRoutes(r, "", "")
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
