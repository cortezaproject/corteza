package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/system/internal/service"
)

type (
	ExternalAuth struct {
		auth       service.AuthService
		jwtEncoder auth.TokenEncoder
	}
)

const (
	externalAuthBaseUrl = "/auth/external"
)

func NewSocial(jwtEncoder auth.TokenEncoder) *ExternalAuth {
	return &ExternalAuth{
		auth:       service.DefaultAuth,
		jwtEncoder: jwtEncoder,
	}
}

func (ctrl *ExternalAuth) MountRoutes(r chi.Router) {

	// Make sure we're backwards compatible and redirect /oidc to /auth/external/openid-connect-didmos2
	r.Get("/oidc", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, externalAuthBaseUrl+"/openid-connect-didmos2", http.StatusMovedPermanently)
	})

	// Copy provider from path (Chi URL param) to request context and return it
	copyProviderToContext := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), "provider", chi.URLParam(r, "provider")))
	}

	r.Route(externalAuthBaseUrl+"/{provider}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			r = copyProviderToContext(r)

			// Always set redir cookie, even if not requested. If param is empty, cookie is removed
			ctrl.setSessionCookie(w, r, "redir", r.URL.Query().Get("redir"))

			// try to get the user without re-authenticating
			if user, err := gothic.CompleteUserAuth(w, r); err != nil {
				gothic.BeginAuthHandler(w, r)
			} else {
				// We've successfully singed-in through 3rd party auth
				ctrl.handleSuccessfulAuth(w, r, user)
			}
		})

		r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
			r = copyProviderToContext(r)

			if user, err := gothic.CompleteUserAuth(w, r); err != nil {
				log.Printf("Failed to complete user auth: %v", err)
				ctrl.handleFailedCallback(w, r, err)
			} else {
				ctrl.handleSuccessfulAuth(w, r, user)
			}
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			if err := gothic.Logout(w, r); err != nil {
				log.Printf("Failed to external logout: %v", err)
			}

			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		})
	})
}

func (ctrl *ExternalAuth) handleFailedCallback(w http.ResponseWriter, r *http.Request, err error) {
	provider := chi.URLParam(r, "provider")

	if strings.Contains(err.Error(), "Error processing your OAuth request: Invalid oauth_verifier parameter") {
		// Just take user through the same loop again
		w.Header().Set("Location", externalAuthBaseUrl+"/"+provider)
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, "SSO Error: %v", err.Error())
	w.WriteHeader(http.StatusOK)
}

// Handles authentication via external auth providers of
// unknown an user + appending authentication on external providers
// to a current user
func (ctrl *ExternalAuth) handleSuccessfulAuth(w http.ResponseWriter, r *http.Request, cred goth.User) {
	log.Printf("Successful external login: %v", cred)

	if u, err := ctrl.auth.With(r.Context()).External(cred); err != nil {
		resputil.JSON(w, err)
	} else {
		ctrl.jwtEncoder.SetCookie(w, r, u)

		if c, err := r.Cookie("redir"); c != nil && err == nil {
			ctrl.setSessionCookie(w, r, "redir", "")
			w.Header().Set("Location", c.Value)
			w.WriteHeader(http.StatusSeeOther)

		}

		resputil.JSON(w, u, err)
	}
}

// Extracts and authenticates JWT from context, validates claims
func (ctrl *ExternalAuth) setSessionCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := &http.Cookie{
		Name:   name,
		Secure: r.URL.Scheme == "https",
		Domain: r.URL.Hostname(),
		Path:   externalAuthBaseUrl,
	}

	if value == "" {
		cookie.Expires = time.Unix(0, 0)
	} else {
		cookie.Value = value
	}

	http.SetCookie(w, cookie)
}
