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
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
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
				logger.LogParamError("AuthInternal.Login", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Login(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.Login", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.Login", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Signup: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalSignup()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AuthInternal.Signup", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Signup(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.Signup", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.Signup", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		RequestPasswordReset: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalRequestPasswordReset()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AuthInternal.RequestPasswordReset", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.RequestPasswordReset(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.RequestPasswordReset", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.RequestPasswordReset", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ExchangePasswordResetToken: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalExchangePasswordResetToken()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AuthInternal.ExchangePasswordResetToken", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ExchangePasswordResetToken(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.ExchangePasswordResetToken", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.ExchangePasswordResetToken", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ResetPassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalResetPassword()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AuthInternal.ResetPassword", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ResetPassword(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.ResetPassword", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.ResetPassword", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ConfirmEmail: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalConfirmEmail()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AuthInternal.ConfirmEmail", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ConfirmEmail(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.ConfirmEmail", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.ConfirmEmail", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ChangePassword: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthInternalChangePassword()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AuthInternal.ChangePassword", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ChangePassword(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AuthInternal.ChangePassword", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AuthInternal.ChangePassword", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
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
