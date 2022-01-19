package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza-server/auth/external"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func copyProviderToContext(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "provider", chi.URLParam(r, "provider")))
}

func (h AuthHandlers) externalInit(w http.ResponseWriter, r *http.Request) {
	r = copyProviderToContext(r)
	h.Log.Info("starting external authentication flow")

	beginUserAuth(w, r, external.NewDefaultExternalHandler())
}

func (h AuthHandlers) externalCallback(w http.ResponseWriter, r *http.Request) {
	r = copyProviderToContext(r)
	h.Log.Info("external authentication callback")

	if user, err := completeUserAuth(w, r, external.NewDefaultExternalHandler()); err != nil {
		h.Log.Error("failed to complete user auth", zap.Error(err))
		h.handleFailedExternalAuth(w, r, err)
	} else {
		h.handleSuccessfulExternalAuth(w, r, *user)
	}
}

// Handles authentication via external auth providers of
// unknown an user + appending authentication on external providers
// to a current user
func (h AuthHandlers) handleSuccessfulExternalAuth(w http.ResponseWriter, r *http.Request, cred types.ExternalAuthUser) {
	var (
		user *types.User
		err  error
		ctx  = r.Context()
	)

	h.Log.Info("login successful", zap.String("provider", cred.Provider))

	// Try to login/sign-up external user
	if user, err = h.AuthService.External(ctx, cred); err != nil {
		api.Send(w, r, err)
		return
	}

	h.handle(func(req *request.AuthReq) error {
		req.AuthUser = request.NewAuthUser(
			h.Settings,
			user,

			// external logins are never permanent!
			false,

			// set def. lifetime for this session
			h.Opt.SessionLifetime,
		)

		// auto-complete EmailOTP and TOTP when authenticating via external identity provider
		req.AuthUser.CompleteEmailOTP()
		req.AuthUser.CompleteTOTP()

		req.AuthUser.Save(req.Session)

		handleSuccessfulAuth(req)
		return nil
	})(w, r)
}

func (h AuthHandlers) handleFailedExternalAuth(w http.ResponseWriter, _ *http.Request, err error) {
	if strings.Contains(err.Error(), "Error processing your OAuth request: Invalid oauth_verifier parameter") {
		// Just take user through the same loop again
		w.Header().Set("Location", GetLinks().Profile)
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "SSO Error: %v", err.Error())
}

func beginUserAuth(w http.ResponseWriter, r *http.Request, eh external.ExternalAuthHandler) {
	eh.BeginUserAuth(w, r)
}

func completeUserAuth(w http.ResponseWriter, r *http.Request, eh external.ExternalAuthHandler) (u *types.ExternalAuthUser, err error) {
	return eh.CompleteUserAuth(w, r)
}
