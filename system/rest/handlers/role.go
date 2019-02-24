package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `role.go`, `role.util.go` or `role_test.go` to
	implement your API calls, helper functions and tests. The file `role.go`
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
type RoleAPI interface {
	List(context.Context, *request.RoleList) (interface{}, error)
	Create(context.Context, *request.RoleCreate) (interface{}, error)
	Update(context.Context, *request.RoleUpdate) (interface{}, error)
	Read(context.Context, *request.RoleRead) (interface{}, error)
	Delete(context.Context, *request.RoleDelete) (interface{}, error)
	Archive(context.Context, *request.RoleArchive) (interface{}, error)
	Move(context.Context, *request.RoleMove) (interface{}, error)
	Merge(context.Context, *request.RoleMerge) (interface{}, error)
	MemberList(context.Context, *request.RoleMemberList) (interface{}, error)
	MemberAdd(context.Context, *request.RoleMemberAdd) (interface{}, error)
	MemberRemove(context.Context, *request.RoleMemberRemove) (interface{}, error)
}

// HTTP API interface
type Role struct {
	List         func(http.ResponseWriter, *http.Request)
	Create       func(http.ResponseWriter, *http.Request)
	Update       func(http.ResponseWriter, *http.Request)
	Read         func(http.ResponseWriter, *http.Request)
	Delete       func(http.ResponseWriter, *http.Request)
	Archive      func(http.ResponseWriter, *http.Request)
	Move         func(http.ResponseWriter, *http.Request)
	Merge        func(http.ResponseWriter, *http.Request)
	MemberList   func(http.ResponseWriter, *http.Request)
	MemberAdd    func(http.ResponseWriter, *http.Request)
	MemberRemove func(http.ResponseWriter, *http.Request)
}

func NewRole(rh RoleAPI) *Role {
	return &Role{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Create(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Update(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Read(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Delete(r.Context(), params)
			})
		},
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleArchive()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Archive(r.Context(), params)
			})
		},
		Move: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Move(r.Context(), params)
			})
		},
		Merge: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMerge()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.Merge(r.Context(), params)
			})
		},
		MemberList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.MemberList(r.Context(), params)
			})
		},
		MemberAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberAdd()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.MemberAdd(r.Context(), params)
			})
		},
		MemberRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberRemove()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return rh.MemberRemove(r.Context(), params)
			})
		},
	}
}

func (rh *Role) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/roles", func(r chi.Router) {
			r.Get("/", rh.List)
			r.Post("/", rh.Create)
			r.Put("/{roleID}", rh.Update)
			r.Get("/{roleID}", rh.Read)
			r.Delete("/{roleID}", rh.Delete)
			r.Post("/{roleID}/archive", rh.Archive)
			r.Post("/{roleID}/move", rh.Move)
			r.Post("/{roleID}/merge", rh.Merge)
			r.Get("/{roleID}/members", rh.MemberList)
			r.Post("/{roleID}/member/{userID}", rh.MemberAdd)
			r.Delete("/{roleID}/member/{userID}", rh.MemberRemove)
		})
	})
}
