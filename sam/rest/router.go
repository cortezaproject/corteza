package rest

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `user.go`, `user.util.go` or `user_test.go` to
	implement your API calls, helper functions and tests. The file `user.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/sam/rest/server"
)

func MountRoutes(r chi.Router) {
	channel := &server.ChannelHandlers{Channel: Channel{}.New()}
	message := &server.MessageHandlers{Message: Message{}.New()}
	organisation := &server.OrganisationHandlers{Organisation: Organisation{}.New()}
	team := &server.TeamHandlers{Team: Team{}.New()}
	user := &server.UserHandlers{User: User{}.New()}
	r.Group(func(r chi.Router) {
		r.Use(channel.Channel.Authenticator())
		r.Route("/channel", func(r chi.Router) {
			r.Get("/", channel.List)
			r.Put("/", channel.Create)
			r.Post("/edit", channel.Edit)
			r.Get("/read", channel.Read)
			r.Delete("/delete", channel.Delete)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(message.Message.Authenticator())
		r.Route("/channels/{channelId}/messages", func(r chi.Router) {
			r.Post("/", message.Create)
			r.Put("/{messageId}", message.Edit)
			r.Delete("/{messageId}", message.Delete)
			r.Put("/{messageId}/attach", message.Attach)
			r.Get("/search", message.Search)
			r.Post("/{messageId}/pin", message.Pin)
			r.Delete("/{messageId}/pin", message.Unpin)
			r.Post("/{messageId}/flag", message.Flag)
			r.Delete("/{messageId}/flag", message.Deflag)
			r.Post("/{messageId}/react", message.React)
			r.Delete("/{messageId}/react/{reactionId}", message.Unreact)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(organisation.Organisation.Authenticator())
		r.Route("/organisation", func(r chi.Router) {
			r.Get("/", organisation.List)
			r.Put("/", organisation.Create)
			r.Post("/{id}", organisation.Edit)
			r.Delete("/{id}", organisation.Remove)
			r.Get("/{id}", organisation.Read)
			r.Post("/{id}/archive", organisation.Archive)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(team.Team.Authenticator())
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
	})
	r.Group(func(r chi.Router) {
		r.Use(user.User.Authenticator())
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", user.Login)
			r.Get("/search", user.Search)
		})
	})
}
