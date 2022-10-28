package server

import (
	"net/http"

	"github.com/99designs/basicauth-go"

	"github.com/go-chi/chi/v5"
	"github.com/ppaanngggg/chi-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsMiddleware is the request logger that provides metrics to prometheus
func metricsMiddleware(name string) func(http.Handler) http.Handler {
	return chiprometheus.NewPatternMiddleware(name)
}

func metricsMount(r chi.Router, username, password string) {
	r.Route("/metrics", func(r chi.Router) {
		if len(password) > 0 {
			r.Use(basicauth.New("Metrics", map[string][]string{
				username: {password},
			}))
		}
		r.Handle("/", promhttp.Handler())
	})
}
