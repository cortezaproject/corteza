package handlers

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func copyProviderToContext(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "provider", chi.URLParam(r, "provider")))
}

func (h AuthHandlers) externalInit(w http.ResponseWriter, r *http.Request) {
	r = copyProviderToContext(r)
	h.Log.Info("starting external authentication flow")

	gothic.BeginAuthHandler(w, r)
}

func (h AuthHandlers) externalCallback(w http.ResponseWriter, r *http.Request) {
	r = copyProviderToContext(r)
	h.Log.Info("external authentication callback")

	if user, err := gothic.CompleteUserAuth(w, r); err != nil {
		h.Log.Error("failed to complete user auth", zap.Error(err))
		h.handleFailedExternalAuth(w, r, err)
	} else {
		h.handleSuccessfulExternalAuth(w, r, user)
	}
}

// Handles authentication via external auth providers of
// unknown an user + appending authentication on external providers
// to a current user
func (h AuthHandlers) handleSuccessfulExternalAuth(w http.ResponseWriter, r *http.Request, cred goth.User) {
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

func (h AuthHandlers) handleFailedExternalAuth(w http.ResponseWriter, r *http.Request, err error) {
	//provider := chi.URLParam(r, "provider")

	if strings.Contains(err.Error(), "Error processing your OAuth request: Invalid oauth_verifier parameter") {
		// Just take user through the same loop again
		w.Header().Set("Location", GetLinks().Profile)
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, "SSO Error: %v", err.Error())
	w.WriteHeader(http.StatusOK)
}
