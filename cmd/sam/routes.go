package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// MountRoutes will register API routes
func MountRoutes(r chi.Router, opts *RouteOptions, mountRoutes func(r chi.Router)) {
	// CORS for local development...
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	if opts.enableLogging {
		r.Use(middleware.Logger)
	}

	mountRoutes(r)
}
