package handlers

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
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/sam/rest/request"
)

// Internal API interface
type UserAPI interface {
	Search(context.Context, *request.UserSearch) (interface{}, error)
	Message(context.Context, *request.UserMessage) (interface{}, error)
}

// HTTP API interface
type User struct {
	Search  func(http.ResponseWriter, *http.Request)
	Message func(http.ResponseWriter, *http.Request)
}

func NewUser(uh UserAPI) *User {
	return &User{
		Search: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserSearch()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Search(r.Context(), params)
			})
		},
		Message: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewUserMessage()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return uh.Message(r.Context(), params)
			})
		},
	}
}

func (uh *User) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/users", func(r chi.Router) {
			r.Get("/search", uh.Search)
			r.Post("/{userID}/message", uh.Message)
		})
	})
}
