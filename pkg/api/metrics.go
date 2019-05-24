package api

import (
	"net/http"

	"github.com/766b/chi-prometheus"
	"github.com/99designs/basicauth-go"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Middleware is the request logger that provides metrics to prometheus
func Middleware(name string) func(http.Handler) http.Handler {
	return chiprometheus.NewMiddleware(name)
}

func Mount(r chi.Router, username, password string) {
	r.Group(func(r chi.Router) {
		r.Use(basicauth.New("Metrics", map[string][]string{
			username: {password},
		}))
		r.Handle("/metrics", promhttp.Handler())
	})
}
