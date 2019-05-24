package api

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Debug(r chi.Router) {
	r.Mount("/debug", middleware.Profiler())
	DebugRoutes(r)
}

func DebugRoutes(r chi.Router) {
	r.Get("/debug/routes", func(w http.ResponseWriter, req *http.Request) {
		var printRoutes func(chi.Routes, string)

		printRoutes = func(r chi.Routes, pfix string) {
			routes := r.Routes()
			for _, route := range routes {
				if route.SubRoutes != nil && len(route.SubRoutes.Routes()) > 0 {
					printRoutes(route.SubRoutes, pfix+route.Pattern[:len(route.Pattern)-2])
				} else {
					for method, fn := range route.Handlers {
						fmt.Fprintf(w, "%-8s %-80s -> %s\n", method, pfix+route.Pattern, runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
					}
				}
			}
		}

		printRoutes(r, "")
	})
}
