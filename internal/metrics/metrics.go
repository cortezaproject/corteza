package metrics

import (
	"net/http"

	"github.com/766b/chi-prometheus"
	"github.com/99designs/basicauth-go"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/crusttech/crust/internal/config"
)

// Middleware is the request logger that provides metrics to prometheus
func Middleware(name string) func(http.Handler) http.Handler {
	return chiprometheus.NewMiddleware(name)
}

// Handler exports prometheus metrics for /metrics requests
func Handler() http.Handler {
	return prometheus.Handler()
}

func MountRoutes(r chi.Router, opts *config.HTTP) {
	if opts.Metrics {
		r.Group(func(r chi.Router) {
			r.Use(basicauth.New("Metrics", map[string][]string{
				opts.MetricsUsername: {opts.MetricsPassword},
			}))
			r.Handle("/metrics", Handler())
		})
	}
}
