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
	MarkAsRead(context.Context, *request.MessageMarkAsRead) (interface{}, error)
	Edit(context.Context, *request.MessageEdit) (interface{}, error)
	Delete(context.Context, *request.MessageDelete) (interface{}, error)
	ReplyGet(context.Context, *request.MessageReplyGet) (interface{}, error)
	ReplyCreate(context.Context, *request.MessageReplyCreate) (interface{}, error)
	PinCreate(context.Context, *request.MessagePinCreate) (interface{}, error)
	PinRemove(context.Context, *request.MessagePinRemove) (interface{}, error)
	BookmarkCreate(context.Context, *request.MessageBookmarkCreate) (interface{}, error)
	BookmarkRemove(context.Context, *request.MessageBookmarkRemove) (interface{}, error)
	ReactionCreate(context.Context, *request.MessageReactionCreate) (interface{}, error)
	ReactionRemove(context.Context, *request.MessageReactionRemove) (interface{}, error)
}

// HTTP API interface
type Message struct {
	Create         func(http.ResponseWriter, *http.Request)
	History        func(http.ResponseWriter, *http.Request)
	MarkAsRead     func(http.ResponseWriter, *http.Request)
	Edit           func(http.ResponseWriter, *http.Request)
	Delete         func(http.ResponseWriter, *http.Request)
	ReplyGet       func(http.ResponseWriter, *http.Request)
	ReplyCreate    func(http.ResponseWriter, *http.Request)
	PinCreate      func(http.ResponseWriter, *http.Request)
	PinRemove      func(http.ResponseWriter, *http.Request)
	BookmarkCreate func(http.ResponseWriter, *http.Request)
	BookmarkRemove func(http.ResponseWriter, *http.Request)
	ReactionCreate func(http.ResponseWriter, *http.Request)
	ReactionRemove func(http.ResponseWriter, *http.Request)
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
		MarkAsRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageMarkAsRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.MarkAsRead(r.Context(), params)
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
		ReplyGet: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReplyGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ReplyGet(r.Context(), params)
			})
		},
		ReplyCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReplyCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ReplyCreate(r.Context(), params)
			})
		},
		PinCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.PinCreate(r.Context(), params)
			})
		},
		PinRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.PinRemove(r.Context(), params)
			})
		},
		BookmarkCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.BookmarkCreate(r.Context(), params)
			})
		},
		BookmarkRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.BookmarkRemove(r.Context(), params)
			})
		},
		ReactionCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ReactionCreate(r.Context(), params)
			})
		},
		ReactionRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return mh.ReactionRemove(r.Context(), params)
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
			r.Get("/mark-as-read", mh.MarkAsRead)
			r.Put("/{messageID}", mh.Edit)
			r.Delete("/{messageID}", mh.Delete)
			r.Get("/{messageID}/replies", mh.ReplyGet)
			r.Post("/{messageID}/replies", mh.ReplyCreate)
			r.Post("/{messageID}/pin", mh.PinCreate)
			r.Delete("/{messageID}/pin", mh.PinRemove)
			r.Post("/{messageID}/bookmark", mh.BookmarkCreate)
			r.Delete("/{messageID}/bookmark", mh.BookmarkRemove)
			r.Post("/{messageID}/reaction/{reaction}", mh.ReactionCreate)
			r.Delete("/{messageID}/reaction/{reaction}", mh.ReactionRemove)
		})
	})
}
