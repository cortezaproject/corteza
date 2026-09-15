package server

import (
	"net/http"

	"github.com/go-chi/cors"
)

// Sets up CORS rules with allowed origins to use as a middleware
func handleCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-ID",
		},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler
}
