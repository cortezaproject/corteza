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
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
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

func NewChannel(ch ChannelAPI) *Channel {
	return &Channel{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Create(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Update(r.Context(), params)
			})
		},
		State: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelState()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.State(r.Context(), params)
			})
		},
		SetFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelSetFlag()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.SetFlag(r.Context(), params)
			})
		},
		RemoveFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRemoveFlag()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.RemoveFlag(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Read(r.Context(), params)
			})
		},
		Members: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelMembers()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Members(r.Context(), params)
			})
		},
		Join: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelJoin()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Join(r.Context(), params)
			})
		},
		Part: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelPart()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Part(r.Context(), params)
			})
		},
		Invite: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelInvite()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Invite(r.Context(), params)
			})
		},
		Attach: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelAttach()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Attach(r.Context(), params)
			})
		},
	}
}

func (ch *Channel) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/channels", func(r chi.Router) {
			r.Get("/", ch.List)
			r.Post("/", ch.Create)
			r.Put("/{channelID}", ch.Update)
			r.Put("/{channelID}/state", ch.State)
			r.Put("/{channelID}/flag", ch.SetFlag)
			r.Delete("/{channelID}/flag", ch.RemoveFlag)
			r.Get("/{channelID}", ch.Read)
			r.Get("/{channelID}/members", ch.Members)
			r.Put("/{channelID}/members/{userID}", ch.Join)
			r.Delete("/{channelID}/members/{userID}", ch.Part)
			r.Post("/{channelID}/invite", ch.Invite)
			r.Post("/{channelID}/attach", ch.Attach)
		})
	})
}
