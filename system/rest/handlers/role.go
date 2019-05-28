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

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
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

func NewRole(h RoleAPI) *Role {
	return &Role{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Archive: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleArchive()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Archive", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Archive(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Archive", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Archive", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Move: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Move", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Move(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Move", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Move", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Merge: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMerge()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.Merge", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Merge(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.Merge", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.Merge", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MemberList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.MemberList", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MemberList(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.MemberList", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.MemberList", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MemberAdd: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberAdd()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.MemberAdd", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MemberAdd(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.MemberAdd", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.MemberAdd", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		MemberRemove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewRoleMemberRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Role.MemberRemove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.MemberRemove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Role.MemberRemove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Role.MemberRemove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Role) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/roles/", h.List)
		r.Post("/roles/", h.Create)
		r.Put("/roles/{roleID}", h.Update)
		r.Get("/roles/{roleID}", h.Read)
		r.Delete("/roles/{roleID}", h.Delete)
		r.Post("/roles/{roleID}/archive", h.Archive)
		r.Post("/roles/{roleID}/move", h.Move)
		r.Post("/roles/{roleID}/merge", h.Merge)
		r.Get("/roles/{roleID}/members", h.MemberList)
		r.Post("/roles/{roleID}/member/{userID}", h.MemberAdd)
		r.Delete("/roles/{roleID}/member/{userID}", h.MemberRemove)
	})
}
