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

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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

func NewMessage(h MessageAPI) *Message {
	return &Message{
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ExecuteCommand: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageExecuteCommand()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.ExecuteCommand", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ExecuteCommand(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.ExecuteCommand", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.ExecuteCommand", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MarkAsRead: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageMarkAsRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.MarkAsRead", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MarkAsRead(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.MarkAsRead", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.MarkAsRead", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageEdit()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.Edit", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Edit(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.Edit", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.Edit", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ReplyCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReplyCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.ReplyCreate", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReplyCreate(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.ReplyCreate", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.ReplyCreate", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		PinCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.PinCreate", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.PinCreate(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.PinCreate", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.PinCreate", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		PinRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessagePinRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.PinRemove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.PinRemove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.PinRemove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.PinRemove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		BookmarkCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.BookmarkCreate", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.BookmarkCreate(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.BookmarkCreate", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.BookmarkCreate", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		BookmarkRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageBookmarkRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.BookmarkRemove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.BookmarkRemove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.BookmarkRemove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.BookmarkRemove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ReactionCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.ReactionCreate", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReactionCreate(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.ReactionCreate", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.ReactionCreate", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ReactionRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMessageReactionRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Message.ReactionRemove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReactionRemove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Message.ReactionRemove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Message.ReactionRemove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
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
