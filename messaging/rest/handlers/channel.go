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
	ChannelAPI interface {
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
	Channel struct {
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
)

func NewChannel(h ChannelAPI) *Channel {
	return &Channel{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelCreate()
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		State: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelState()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.State(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		SetFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelSetFlag()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.SetFlag(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RemoveFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRemoveFlag()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RemoveFlag(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Members: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelMembers()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Members(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Join: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelJoin()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Join(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Part: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelPart()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Part(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Invite: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelInvite()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Invite(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Attach: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelAttach()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Attach(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
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
