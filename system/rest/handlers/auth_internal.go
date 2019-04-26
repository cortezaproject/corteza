package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `auth_internal.go`, `auth_internal.util.go` or `auth_internal_test.go` to
	implement your API calls, helper functions and tests. The file `auth_internal.go`
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
type AuthInternalAPI interface {
	Login(context.Context, *request.AuthInternalLogin) (interface{}, error)
	Signup(context.Context, *request.AuthInternalSignup) (interface{}, error)
	RequestPasswordReset(context.Context, *request.AuthInternalRequestPasswordReset) (interface{}, error)
	ExchangePasswordResetToken(context.Context, *request.AuthInternalExchangePasswordResetToken) (interface{}, error)
	ResetPassword(context.Context, *request.AuthInternalResetPassword) (interface{}, error)
	ConfirmEmail(context.Context, *request.AuthInternalConfirmEmail) (interface{}, error)
	ChangePassword(context.Context, *request.AuthInternalChangePassword) (interface{}, error)
}

// HTTP API interface
type AuthInternal struct {
	Login                      func(http.ResponseWriter, *http.Request)
	Signup                     func(http.ResponseWriter, *http.Request)
	RequestPasswordReset       func(http.ResponseWriter, *http.Request)
	ExchangePasswordResetToken func(http.ResponseWriter, *http.Request)
	ResetPassword              func(http.ResponseWriter, *http.Request)
	ConfirmEmail               func(http.ResponseWriter, *http.Request)
	ChangePassword             func(http.ResponseWriter, *http.Request)
}

func NewAuthInternal(ah AuthInternalAPI) *AuthInternal {
	return &AuthInternal{
		Login: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalLogin()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.Login(r.Context(), params); err != nil {
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
		Signup: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalSignup()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.Signup(r.Context(), params); err != nil {
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
		RequestPasswordReset: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalRequestPasswordReset()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.RequestPasswordReset(r.Context(), params); err != nil {
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
		ExchangePasswordResetToken: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalExchangePasswordResetToken()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.ExchangePasswordResetToken(r.Context(), params); err != nil {
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
		ResetPassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalResetPassword()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.ResetPassword(r.Context(), params); err != nil {
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
		ConfirmEmail: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalConfirmEmail()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.ConfirmEmail(r.Context(), params); err != nil {
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
		ChangePassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalChangePassword()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := ah.ChangePassword(r.Context(), params); err != nil {
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

func (ah *AuthInternal) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/auth/internal/login", ah.Login)
		r.Post("/auth/internal/signup", ah.Signup)
		r.Post("/auth/internal/request-password-reset", ah.RequestPasswordReset)
		r.Post("/auth/internal/exchange-password-reset-token", ah.ExchangePasswordResetToken)
		r.Post("/auth/internal/reset-password", ah.ResetPassword)
		r.Post("/auth/internal/confirm-email", ah.ConfirmEmail)
		r.Post("/auth/internal/change-password", ah.ChangePassword)
	})
}
