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
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.List(r.Context(), params); err != nil {
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
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Create(r.Context(), params); err != nil {
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelUpdate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Update(r.Context(), params); err != nil {
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
		State: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelState()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.State(r.Context(), params); err != nil {
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
		SetFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelSetFlag()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.SetFlag(r.Context(), params); err != nil {
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
		RemoveFlag: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRemoveFlag()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.RemoveFlag(r.Context(), params); err != nil {
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelRead()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Read(r.Context(), params); err != nil {
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
		Members: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelMembers()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Members(r.Context(), params); err != nil {
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
		Join: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelJoin()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Join(r.Context(), params); err != nil {
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
		Part: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelPart()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Part(r.Context(), params); err != nil {
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
		Invite: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelInvite()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Invite(r.Context(), params); err != nil {
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
		Attach: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChannelAttach()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Attach(r.Context(), params); err != nil {
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

func (ch *Channel) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/channels/", ch.List)
		r.Post("/channels/", ch.Create)
		r.Put("/channels/{channelID}", ch.Update)
		r.Put("/channels/{channelID}/state", ch.State)
		r.Put("/channels/{channelID}/flag", ch.SetFlag)
		r.Delete("/channels/{channelID}/flag", ch.RemoveFlag)
		r.Get("/channels/{channelID}", ch.Read)
		r.Get("/channels/{channelID}/members", ch.Members)
		r.Put("/channels/{channelID}/members/{userID}", ch.Join)
		r.Delete("/channels/{channelID}/members/{userID}", ch.Part)
		r.Post("/channels/{channelID}/invite", ch.Invite)
		r.Post("/channels/{channelID}/attach", ch.Attach)
	})
}
