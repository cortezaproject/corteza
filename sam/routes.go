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
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
	channel := ChannelHandlers{}.new()
	message := MessageHandlers{}.new()
	organisation := OrganisationHandlers{}.new()
	team := TeamHandlers{}.new()
	user := UserHandlers{}.new()
	websocket := WebsocketHandlers{}.new()
	r.Route("/channel", func(r chi.Router) {
		r.Post("/edit", channel.Edit)
		r.Delete("/remove", channel.Remove)
		r.Get("/read", channel.Read)
		r.Get("/search", channel.Search)
		r.Post("/archive", channel.Archive)
		r.Post("/move", channel.Move)
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
		r.Post("/edit", organisation.Edit)
		r.Delete("/remove", organisation.Remove)
		r.Get("/read", organisation.Read)
		r.Get("/search", organisation.Search)
		r.Post("/archive", organisation.Archive)
	})
	r.Route("/team", func(r chi.Router) {
		r.Post("/edit", team.Edit)
		r.Delete("/remove", team.Remove)
		r.Get("/read", team.Read)
		r.Get("/search", team.Search)
		r.Post("/archive", team.Archive)
		r.Post("/move", team.Move)
		r.Post("/merge", team.Merge)
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/login", user.Login)
		r.Get("/search", user.Search)
	})
	r.Route("/websocket", func(r chi.Router) {
		r.Get("/client", websocket.Client)
	})
}
