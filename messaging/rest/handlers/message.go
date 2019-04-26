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

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
)

// Internal API interface
type MessageAPI interface {
	Create(context.Context, *request.MessageCreate) (interface{}, error)
	ExecuteCommand(context.Context, *request.MessageExecuteCommand) (interface{}, error)
	MarkAsRead(context.Context, *request.MessageMarkAsRead) (interface{}, error)
	Edit(context.Context, *request.MessageEdit) (interface{}, error)
	Delete(context.Context, *request.MessageDelete) (interface{}, error)
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
	ExecuteCommand func(http.ResponseWriter, *http.Request)
	MarkAsRead     func(http.ResponseWriter, *http.Request)
	Edit           func(http.ResponseWriter, *http.Request)
	Delete         func(http.ResponseWriter, *http.Request)
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
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.Create(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		ExecuteCommand: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageExecuteCommand()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.ExecuteCommand(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		MarkAsRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageMarkAsRead()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.MarkAsRead(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageEdit()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.Edit(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageDelete()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.Delete(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		ReplyCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReplyCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.ReplyCreate(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		PinCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.PinCreate(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		PinRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinRemove()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.PinRemove(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		BookmarkCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.BookmarkCreate(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		BookmarkRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkRemove()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.BookmarkRemove(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		ReactionCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.ReactionCreate(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		ReactionRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionRemove()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := mh.ReactionRemove(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
	}
}

func (mh *Message) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/channels/{channelID}/messages/", mh.Create)
		r.Post("/channels/{channelID}/messages/command/{command}/exec", mh.ExecuteCommand)
		r.Get("/channels/{channelID}/messages/mark-as-read", mh.MarkAsRead)
		r.Put("/channels/{channelID}/messages/{messageID}", mh.Edit)
		r.Delete("/channels/{channelID}/messages/{messageID}", mh.Delete)
		r.Post("/channels/{channelID}/messages/{messageID}/replies", mh.ReplyCreate)
		r.Post("/channels/{channelID}/messages/{messageID}/pin", mh.PinCreate)
		r.Delete("/channels/{channelID}/messages/{messageID}/pin", mh.PinRemove)
		r.Post("/channels/{channelID}/messages/{messageID}/bookmark", mh.BookmarkCreate)
		r.Delete("/channels/{channelID}/messages/{messageID}/bookmark", mh.BookmarkRemove)
		r.Post("/channels/{channelID}/messages/{messageID}/reaction/{reaction}", mh.ReactionCreate)
		r.Delete("/channels/{channelID}/messages/{messageID}/reaction/{reaction}", mh.ReactionRemove)
	})
}
