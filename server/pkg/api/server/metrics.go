package server

import (
	"github.com/99designs/basicauth-go"
	"net/http"

	"github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsMiddleware is the request logger that provides metrics to prometheus
func metricsMiddleware(name string) func(http.Handler) http.Handler {
	return chiprometheus.NewMiddleware(name)
}

func metricsMount(r chi.Router, username, password string) {
	r.Route("/metrics", func(r chi.Router) {
		r.Use(basicauth.New("Metrics", map[string][]string{
			username: {password},
		}))
		r.Handle("/", promhttp.Handler())
	})
}
