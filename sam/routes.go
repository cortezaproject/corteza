package sam

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
	channel := ChannelHandlers{}.new()
	message := MessageHandlers{}.new()
	organisation := OrganisationHandlers{}.new()
	team := TeamHandlers{}.new()
	user := UserHandlers{}.new()
	websocket := WebsocketHandlers{}.new()
	r.Route("/channel", func(r chi.Router) {
		r.Get("/", channel.List)
		r.Put("/", channel.Create)
		r.Post("/edit", channel.Edit)
		r.Get("/read", channel.Read)
		r.Delete("/delete", channel.Delete)
	})
	r.Route("/message", func(r chi.Router) {
		r.Post("/edit", message.Edit)
		r.Put("/attach", message.Attach)
		r.Delete("/remove", message.Remove)
		r.Get("/read", message.Read)
		r.Get("/search", message.Search)
		r.Post("/pin", message.Pin)
		r.Post("/flag", message.Flag)
	})
	r.Route("/organisation", func(r chi.Router) {
		r.Get("/", organisation.List)
		r.Put("/", organisation.Create)
		r.Post("/{id}", organisation.Edit)
		r.Delete("/{id}", organisation.Remove)
		r.Get("/{id}", organisation.Read)
		r.Post("/{id}/archive", organisation.Archive)
	})
	r.Route("/team", func(r chi.Router) {
		r.Get("/", team.List)
		r.Put("/", team.Create)
		r.Post("/{id}", team.Edit)
		r.Get("/{id}", team.Read)
		r.Delete("/{id}", team.Remove)
		r.Post("/{id}/archive", team.Archive)
		r.Post("/{id}/move", team.Move)
		r.Post("/{id}/merge", team.Merge)
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/login", user.Login)
		r.Get("/search", user.Search)
	})
	r.Route("/websocket", func(r chi.Router) {
		r.Get("/client", websocket.Client)
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
