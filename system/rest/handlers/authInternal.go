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
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	AuthInternalAPI interface {
		Login(context.Context, *request.AuthInternalLogin) (interface{}, error)
		Signup(context.Context, *request.AuthInternalSignup) (interface{}, error)
		RequestPasswordReset(context.Context, *request.AuthInternalRequestPasswordReset) (interface{}, error)
		ExchangePasswordResetToken(context.Context, *request.AuthInternalExchangePasswordResetToken) (interface{}, error)
		ResetPassword(context.Context, *request.AuthInternalResetPassword) (interface{}, error)
		ConfirmEmail(context.Context, *request.AuthInternalConfirmEmail) (interface{}, error)
		ChangePassword(context.Context, *request.AuthInternalChangePassword) (interface{}, error)
	}

	// HTTP API interface
	AuthInternal struct {
		Login                      func(http.ResponseWriter, *http.Request)
		Signup                     func(http.ResponseWriter, *http.Request)
		RequestPasswordReset       func(http.ResponseWriter, *http.Request)
		ExchangePasswordResetToken func(http.ResponseWriter, *http.Request)
		ResetPassword              func(http.ResponseWriter, *http.Request)
		ConfirmEmail               func(http.ResponseWriter, *http.Request)
		ChangePassword             func(http.ResponseWriter, *http.Request)
	}
)

func NewAuthInternal(h AuthInternalAPI) *AuthInternal {
	return &AuthInternal{
		Login: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalLogin()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Login(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Signup: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalSignup()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Signup(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RequestPasswordReset: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalRequestPasswordReset()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RequestPasswordReset(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ExchangePasswordResetToken: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalExchangePasswordResetToken()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ExchangePasswordResetToken(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ResetPassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalResetPassword()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ResetPassword(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ConfirmEmail: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalConfirmEmail()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ConfirmEmail(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ChangePassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalChangePassword()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ChangePassword(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h AuthInternal) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/auth/internal/login", h.Login)
		r.Post("/auth/internal/signup", h.Signup)
		r.Post("/auth/internal/request-password-reset", h.RequestPasswordReset)
		r.Post("/auth/internal/exchange-password-reset-token", h.ExchangePasswordResetToken)
		r.Post("/auth/internal/reset-password", h.ResetPassword)
		r.Post("/auth/internal/confirm-email", h.ConfirmEmail)
		r.Post("/auth/internal/change-password", h.ChangePassword)
	})
}
