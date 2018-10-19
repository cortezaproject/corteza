package handlers

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
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/sam/rest/request"
)

// Internal API interface
type MessageAPI interface {
	Create(context.Context, *request.MessageCreate) (interface{}, error)
	History(context.Context, *request.MessageHistory) (interface{}, error)
	Edit(context.Context, *request.MessageEdit) (interface{}, error)
	Delete(context.Context, *request.MessageDelete) (interface{}, error)
	Search(context.Context, *request.MessageSearch) (interface{}, error)
	Pin(context.Context, *request.MessagePin) (interface{}, error)
	GetReplies(context.Context, *request.MessageGetReplies) (interface{}, error)
	CreateReply(context.Context, *request.MessageCreateReply) (interface{}, error)
	Unpin(context.Context, *request.MessageUnpin) (interface{}, error)
	Flag(context.Context, *request.MessageFlag) (interface{}, error)
	Unflag(context.Context, *request.MessageUnflag) (interface{}, error)
	React(context.Context, *request.MessageReact) (interface{}, error)
	Unreact(context.Context, *request.MessageUnreact) (interface{}, error)
}

// HTTP API interface
type Message struct {
	Create      func(http.ResponseWriter, *http.Request)
	History     func(http.ResponseWriter, *http.Request)
	Edit        func(http.ResponseWriter, *http.Request)
	Delete      func(http.ResponseWriter, *http.Request)
	Search      func(http.ResponseWriter, *http.Request)
	Pin         func(http.ResponseWriter, *http.Request)
	GetReplies  func(http.ResponseWriter, *http.Request)
	CreateReply func(http.ResponseWriter, *http.Request)
	Unpin       func(http.ResponseWriter, *http.Request)
	Flag        func(http.ResponseWriter, *http.Request)
	Unflag      func(http.ResponseWriter, *http.Request)
	React       func(http.ResponseWriter, *http.Request)
	Unreact     func(http.ResponseWriter, *http.Request)
}

func NewMessage(mh MessageAPI) *Message {
	return &Message{
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Create(r.Context(), params)
			})
		},
		History: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageHistory()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.History(r.Context(), params)
			})
		},
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageEdit()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Edit(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Delete(r.Context(), params)
			})
		},
		Search: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageSearch()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Search(r.Context(), params)
			})
		},
		Pin: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Pin(r.Context(), params)
			})
		},
		GetReplies: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageGetReplies()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.GetReplies(r.Context(), params)
			})
		},
		CreateReply: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageCreateReply()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.CreateReply(r.Context(), params)
			})
		},
		Unpin: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageUnpin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Unpin(r.Context(), params)
			})
		},
		Flag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageFlag()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Flag(r.Context(), params)
			})
		},
		Unflag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageUnflag()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Unflag(r.Context(), params)
			})
		},
		React: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReact()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.React(r.Context(), params)
			})
		},
		Unreact: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageUnreact()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.Unreact(r.Context(), params)
			})
		},
	}
}

func (mh *Message) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/channels/{channelID}/messages", func(r chi.Router) {
			r.Post("/", mh.Create)
			r.Get("/", mh.History)
			r.Put("/{messageID}", mh.Edit)
			r.Delete("/{messageID}", mh.Delete)
			r.Get("/search", mh.Search)
			r.Post("/{messageID}/pin", mh.Pin)
			r.Get("/{messageID}/replies", mh.GetReplies)
			r.Post("/{messageID}/replies", mh.CreateReply)
			r.Delete("/{messageID}/pin", mh.Unpin)
			r.Post("/{messageID}/flag", mh.Flag)
			r.Delete("/{messageID}/flag", mh.Unflag)
			r.Put("/{messageID}/reaction/{reaction}", mh.React)
			r.Delete("/{messageID}/react/{reaction}", mh.Unreact)
		})
	})
}
