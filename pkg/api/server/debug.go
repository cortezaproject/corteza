package server

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/davecgh/go-spew/spew"

	"github.com/go-chi/chi/v5"
)

func debugRoutes(r chi.Routes) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
	}
}

func debugEventbus() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		spew.Fdump(w, eventbus.Service().Debug())
	}
}

func debugCorredor() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		spew.Fdump(w, corredor.Service().Debug())
	}
}
