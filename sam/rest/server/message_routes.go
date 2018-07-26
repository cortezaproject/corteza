package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `message.go`, `message.util.go` or `message_test.go` to
	implement your API calls, helper functions and tests. The file `message.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (mh *MessageHandlers) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/channels/{channelID}/messages", func(r chi.Router) {
			r.Post("/", mh.Create)
			r.Get("/", mh.History)
			r.Put("/{messageID}", mh.Edit)
			r.Delete("/{messageID}", mh.Delete)
			r.Put("/{messageID}/attach", mh.Attach)
			r.Get("/search", mh.Search)
			r.Post("/{messageID}/pin", mh.Pin)
			r.Delete("/{messageID}/pin", mh.Unpin)
			r.Post("/{messageID}/flag", mh.Flag)
			r.Delete("/{messageID}/flag", mh.Unflag)
			r.Put("/{messageID}/reaction/{reaction}", mh.React)
			r.Delete("/{messageID}/react/{reaction}", mh.Unreact)
		})
	})
}
