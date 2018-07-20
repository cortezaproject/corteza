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
	auth := &server.AuthHandlers{Auth: Auth{}.New()}
	channel := &server.ChannelHandlers{Channel: Channel{}.New()}
	message := &server.MessageHandlers{Message: Message{}.New()}
	organisation := &server.OrganisationHandlers{Organisation: Organisation{}.New()}
	team := &server.TeamHandlers{Team: Team{}.New()}
	user := &server.UserHandlers{User: User{}.New()}
	r.Group(func(r chi.Router) {
		r.Use(auth.Auth.Authenticator())
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", auth.Login)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(channel.Channel.Authenticator())
		r.Route("/channels", func(r chi.Router) {
			r.Get("/", channel.List)
			r.Put("/", channel.Create)
			r.Post("/{channelId}", channel.Edit)
			r.Get("/{channelId}", channel.Read)
			r.Delete("/{channelId}", channel.Delete)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(message.Message.Authenticator())
		r.Route("/channels/{channelId}/messages", func(r chi.Router) {
			r.Post("/", message.Create)
			r.Get("/", message.History)
			r.Put("/{messageId}", message.Edit)
			r.Delete("/{messageId}", message.Delete)
			r.Put("/{messageId}/attach", message.Attach)
			r.Get("/search", message.Search)
			r.Post("/{messageId}/pin", message.Pin)
			r.Delete("/{messageId}/pin", message.Unpin)
			r.Post("/{messageId}/flag", message.Flag)
			r.Delete("/{messageId}/flag", message.Unflag)
			r.Put("/{messageId}/reaction/{reaction}", message.React)
			r.Delete("/{messageId}/react/{reaction}", message.Unreact)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(organisation.Organisation.Authenticator())
		r.Route("/organisations", func(r chi.Router) {
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
		r.Route("/teams", func(r chi.Router) {
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
		r.Route("/users", func(r chi.Router) {
			r.Get("/search", user.Search)
		})
	})
}
