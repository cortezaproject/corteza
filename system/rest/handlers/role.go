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

	"net/http"

	"github.com/go-chi/chi"
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
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.List(r.Context(), params); err != nil {
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
			params := request.NewRoleCreate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Create(r.Context(), params); err != nil {
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
			params := request.NewRoleUpdate()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Update(r.Context(), params); err != nil {
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
			params := request.NewRoleRead()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Read(r.Context(), params); err != nil {
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
			params := request.NewRoleDelete()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Delete(r.Context(), params); err != nil {
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
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleArchive()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Archive(r.Context(), params); err != nil {
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
		Move: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMove()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Move(r.Context(), params); err != nil {
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
		Merge: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMerge()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.Merge(r.Context(), params); err != nil {
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
		MemberList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberList()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.MemberList(r.Context(), params); err != nil {
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
		MemberAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberAdd()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.MemberAdd(r.Context(), params); err != nil {
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
		MemberRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberRemove()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := rh.MemberRemove(r.Context(), params); err != nil {
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

func (rh *Role) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/roles/", rh.List)
		r.Post("/roles/", rh.Create)
		r.Put("/roles/{roleID}", rh.Update)
		r.Get("/roles/{roleID}", rh.Read)
		r.Delete("/roles/{roleID}", rh.Delete)
		r.Post("/roles/{roleID}/archive", rh.Archive)
		r.Post("/roles/{roleID}/move", rh.Move)
		r.Post("/roles/{roleID}/merge", rh.Merge)
		r.Get("/roles/{roleID}/members", rh.MemberList)
		r.Post("/roles/{roleID}/member/{userID}", rh.MemberAdd)
		r.Delete("/roles/{roleID}/member/{userID}", rh.MemberRemove)
	})
}
