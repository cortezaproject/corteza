package main

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// MountRoutes will register API routes
func MountRoutes(r chi.Router, opts *RouteOptions, mountRoutes ...func(r chi.Router)) {
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

	for _, mount := range mountRoutes {
		mount(r)
	}

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

	r.Mount("/debug", middleware.Profiler())
}
