package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `.go`, `.util.go` or `_test.go` to
	implement your API calls, helper functions and tests. The file `.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"fmt"
	"github.com/go-chi/chi"
	"reflect"
	"runtime"
)

func MountRoutes(r chi.Router) {
	module := ModuleHandlers{}.new()
	types := TypesHandlers{}.new()
	r.Route("/module", func(r chi.Router) {
		r.Get("/list", module.List)
		r.Post("/edit", module.Edit)
		r.Get("/content/list", module.ContentList)
		r.Post("/content/edit", module.ContentEdit)
		r.Delete("/content/delete", module.ContentDelete)
	})
	r.Route("/types", func(r chi.Router) {
		r.Get("/list", types.List)
		r.Get("/type/{id}", types.Type)
	})

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
}
