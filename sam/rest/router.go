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
			r.Post("/{channelID}", channel.Edit)
			r.Get("/{channelID}", channel.Read)
			r.Delete("/{channelID}", channel.Delete)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(message.Message.Authenticator())
		r.Route("/channels/{channelID}/messages", func(r chi.Router) {
			r.Post("/", message.Create)
			r.Get("/", message.History)
			r.Put("/{messageID}", message.Edit)
			r.Delete("/{messageID}", message.Delete)
			r.Put("/{messageID}/attach", message.Attach)
			r.Get("/search", message.Search)
			r.Post("/{messageID}/pin", message.Pin)
			r.Delete("/{messageID}/pin", message.Unpin)
			r.Post("/{messageID}/flag", message.Flag)
			r.Delete("/{messageID}/flag", message.Unflag)
			r.Put("/{messageID}/reaction/{reaction}", message.React)
			r.Delete("/{messageID}/react/{reaction}", message.Unreact)
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
			r.Post("/{teamID}", team.Edit)
			r.Get("/{teamID}", team.Read)
			r.Delete("/{teamID}", team.Remove)
			r.Post("/{teamID}/archive", team.Archive)
			r.Post("/{teamID}/move", team.Move)
			r.Post("/{teamID}/merge", team.Merge)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(user.User.Authenticator())
		r.Route("/users", func(r chi.Router) {
			r.Get("/search", user.Search)
		})
	})
}
