package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `team.go`, `team.util.go` or `team_test.go` to
	implement your API calls, helper functions and tests. The file `team.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/system/rest/request"
)

// Internal API interface
type TeamAPI interface {
	List(context.Context, *request.TeamList) (interface{}, error)
	Create(context.Context, *request.TeamCreate) (interface{}, error)
	Edit(context.Context, *request.TeamEdit) (interface{}, error)
	Read(context.Context, *request.TeamRead) (interface{}, error)
	Remove(context.Context, *request.TeamRemove) (interface{}, error)
	Archive(context.Context, *request.TeamArchive) (interface{}, error)
	Move(context.Context, *request.TeamMove) (interface{}, error)
	Merge(context.Context, *request.TeamMerge) (interface{}, error)
	MemberAdd(context.Context, *request.TeamMemberAdd) (interface{}, error)
	MemberRemove(context.Context, *request.TeamMemberRemove) (interface{}, error)
}

// HTTP API interface
type Team struct {
	List         func(http.ResponseWriter, *http.Request)
	Create       func(http.ResponseWriter, *http.Request)
	Edit         func(http.ResponseWriter, *http.Request)
	Read         func(http.ResponseWriter, *http.Request)
	Remove       func(http.ResponseWriter, *http.Request)
	Archive      func(http.ResponseWriter, *http.Request)
	Move         func(http.ResponseWriter, *http.Request)
	Merge        func(http.ResponseWriter, *http.Request)
	MemberAdd    func(http.ResponseWriter, *http.Request)
	MemberRemove func(http.ResponseWriter, *http.Request)
}

func NewTeam(th TeamAPI) *Team {
	return &Team{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Create(r.Context(), params)
			})
		},
		Edit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamEdit()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Edit(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Read(r.Context(), params)
			})
		},
		Remove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Remove(r.Context(), params)
			})
		},
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamArchive()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Archive(r.Context(), params)
			})
		},
		Move: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamMove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Move(r.Context(), params)
			})
		},
		Merge: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamMerge()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Merge(r.Context(), params)
			})
		},
		MemberAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamMemberAdd()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.MemberAdd(r.Context(), params)
			})
		},
		MemberRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTeamMemberRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.MemberRemove(r.Context(), params)
			})
		},
	}
}

func (th *Team) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/teams", func(r chi.Router) {
			r.Get("/", th.List)
			r.Post("/", th.Create)
			r.Put("/{teamID}", th.Edit)
			r.Get("/{teamID}", th.Read)
			r.Delete("/{teamID}", th.Remove)
			r.Post("/{teamID}/archive", th.Archive)
			r.Post("/{teamID}/move", th.Move)
			r.Post("/{teamID}/merge", th.Merge)
			r.Post("/{teamID}/memberAdd", th.MemberAdd)
			r.Post("/{teamID}/memberRemove", th.MemberRemove)
		})
	})
}
