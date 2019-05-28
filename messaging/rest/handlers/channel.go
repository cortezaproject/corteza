package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `channel.go`, `channel.util.go` or `channel_test.go` to
	implement your API calls, helper functions and tests. The file `channel.go`
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
type ChannelAPI interface {
	List(context.Context, *request.ChannelList) (interface{}, error)
	Create(context.Context, *request.ChannelCreate) (interface{}, error)
	Update(context.Context, *request.ChannelUpdate) (interface{}, error)
	State(context.Context, *request.ChannelState) (interface{}, error)
	SetFlag(context.Context, *request.ChannelSetFlag) (interface{}, error)
	RemoveFlag(context.Context, *request.ChannelRemoveFlag) (interface{}, error)
	Read(context.Context, *request.ChannelRead) (interface{}, error)
	Members(context.Context, *request.ChannelMembers) (interface{}, error)
	Join(context.Context, *request.ChannelJoin) (interface{}, error)
	Part(context.Context, *request.ChannelPart) (interface{}, error)
	Invite(context.Context, *request.ChannelInvite) (interface{}, error)
	Attach(context.Context, *request.ChannelAttach) (interface{}, error)
}

// HTTP API interface
type Channel struct {
	List       func(http.ResponseWriter, *http.Request)
	Create     func(http.ResponseWriter, *http.Request)
	Update     func(http.ResponseWriter, *http.Request)
	State      func(http.ResponseWriter, *http.Request)
	SetFlag    func(http.ResponseWriter, *http.Request)
	RemoveFlag func(http.ResponseWriter, *http.Request)
	Read       func(http.ResponseWriter, *http.Request)
	Members    func(http.ResponseWriter, *http.Request)
	Join       func(http.ResponseWriter, *http.Request)
	Part       func(http.ResponseWriter, *http.Request)
	Invite     func(http.ResponseWriter, *http.Request)
	Attach     func(http.ResponseWriter, *http.Request)
}

func NewChannel(h ChannelAPI) *Channel {
	return &Channel{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		State: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelState()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.State", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.State(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.State", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.State", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		SetFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelSetFlag()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.SetFlag", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.SetFlag(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.SetFlag", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.SetFlag", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		RemoveFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRemoveFlag()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.RemoveFlag", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.RemoveFlag(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.RemoveFlag", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.RemoveFlag", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Members: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelMembers()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Members", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Members(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Members", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Members", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Join: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelJoin()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Join", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Join(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Join", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Join", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Part: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelPart()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Part", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Part(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Part", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Part", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Invite: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelInvite()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Invite", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Invite(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Invite", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Invite", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Attach: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelAttach()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Channel.Attach", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Attach(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Channel.Attach", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Channel.Attach", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Channel) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/channels/", h.List)
		r.Post("/channels/", h.Create)
		r.Put("/channels/{channelID}", h.Update)
		r.Put("/channels/{channelID}/state", h.State)
		r.Put("/channels/{channelID}/flag", h.SetFlag)
		r.Delete("/channels/{channelID}/flag", h.RemoveFlag)
		r.Get("/channels/{channelID}", h.Read)
		r.Get("/channels/{channelID}/members", h.Members)
		r.Put("/channels/{channelID}/members/{userID}", h.Join)
		r.Delete("/channels/{channelID}/members/{userID}", h.Part)
		r.Post("/channels/{channelID}/invite", h.Invite)
		r.Post("/channels/{channelID}/attach", h.Attach)
	})
}
