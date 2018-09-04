package auth

import (
	"net/http"

	"github.com/766b/chi-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct{}

// Middleware is the request logger that provides metrics to prometheus
func (metrics) Middleware(name string) func(http.Handler) http.Handler {
	return chiprometheus.NewMiddleware(name)
}

// Handler exports prometheus metrics for /metrics requests
func (metrics) Handler() http.Handler {
	return prometheus.Handler()
}
