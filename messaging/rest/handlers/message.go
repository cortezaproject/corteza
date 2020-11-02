package handlers

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	MessageAPI interface {
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
	Message struct {
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
)

func NewMessage(h MessageAPI) *Message {
	return &Message{
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ExecuteCommand: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageExecuteCommand()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ExecuteCommand(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		MarkAsRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageMarkAsRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.MarkAsRead(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageEdit()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Edit(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReplyCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReplyCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReplyCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		PinCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.PinCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		PinRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.PinRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		BookmarkCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.BookmarkCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		BookmarkRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.BookmarkRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReactionCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReactionCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReactionRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionRemove()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReactionRemove(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Message) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/channels/{channelID}/messages/", h.Create)
		r.Post("/channels/{channelID}/messages/command/{command}/exec", h.ExecuteCommand)
		r.Get("/channels/{channelID}/messages/mark-as-read", h.MarkAsRead)
		r.Put("/channels/{channelID}/messages/{messageID}", h.Edit)
		r.Delete("/channels/{channelID}/messages/{messageID}", h.Delete)
		r.Post("/channels/{channelID}/messages/{messageID}/replies", h.ReplyCreate)
		r.Post("/channels/{channelID}/messages/{messageID}/pin", h.PinCreate)
		r.Delete("/channels/{channelID}/messages/{messageID}/pin", h.PinRemove)
		r.Post("/channels/{channelID}/messages/{messageID}/bookmark", h.BookmarkCreate)
		r.Delete("/channels/{channelID}/messages/{messageID}/bookmark", h.BookmarkRemove)
		r.Post("/channels/{channelID}/messages/{messageID}/reaction/{reaction}", h.ReactionCreate)
		r.Delete("/channels/{channelID}/messages/{messageID}/reaction/{reaction}", h.ReactionRemove)
	})
}
