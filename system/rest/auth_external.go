package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/titpetric/factory/resputil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/internal/service"
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

func NewExternalAuth() *ExternalAuth {
	return &ExternalAuth{
		auth:       service.DefaultAuth,
		jwtEncoder: auth.DefaultJwtHandler,
	}
}

func (ctrl ExternalAuth) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.ContextValue(ctx).Named("external-auth").With(fields...)
}

func (ctrl *ExternalAuth) ApiServerRoutes(r chi.Router) {

	// Make sure we're backwards compatible and redirect /oidc to /auth/external/openid-connect.corteza-iam
	r.Get("/oidc", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, externalAuthBaseUrl+"/openid-connect.corteza-iam", http.StatusMovedPermanently)
	})

	// Copy provider from path (Chi URL param) to request context and return it
	copyProviderToContext := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), "provider", chi.URLParam(r, "provider")))
	}

	r.Route(externalAuthBaseUrl+"/{provider}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			r = copyProviderToContext(r)

			// Always set redir cookie, even if not requested.
			// If param is empty, cookie will be removed
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
				ctrl.log(r.Context(), zap.Error(err)).Error("failed to complete user auth")
				ctrl.handleFailedCallback(w, r, err)
			} else {
				ctrl.handleSuccessfulAuth(w, r, user)
			}
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			if err := gothic.Logout(w, r); err != nil {
				ctrl.log(r.Context(), zap.Error(err)).Error("failed to execute external logout")
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
//
//
// Redirection rules:
// 1) use cookie (set from query-string param on first step
// 2) use `auth.frontend.url.redirect` setting
// 3) use current url
func (ctrl *ExternalAuth) handleSuccessfulAuth(w http.ResponseWriter, r *http.Request, cred goth.User) {
	ctrl.log(r.Context(), zap.String("provider", cred.Provider)).Info("external login successful")

	svc := ctrl.auth.With(r.Context())

	if u, err := svc.External(cred); err != nil {
		resputil.JSON(w, err)
	} else {
		var (
			token    string
			redirUrl *url.URL
			c        *http.Cookie
		)

		if c, err = r.Cookie("redir"); c != nil && err == nil {
			if redirUrl, err = url.Parse(c.Value); err == nil {
				// @todo validate origin/redir-domain
				ctrl.setSessionCookie(w, r, "redir", "")
			}
		} else if fru := svc.FrontendRedirectURL(); fru != "" {
			redirUrl, err = url.Parse(fru)
		} else {
			redirUrl = r.URL
		}

		if err != nil {
			resputil.JSON(w, err)
			return
		}

		if redirUrl != nil {

			q := redirUrl.Query()

			if u != nil {
				if token, err = svc.IssueAuthRequestToken(u); err == nil {
					q.Set("token", token)
				}
			}

			if err != nil {
				q.Set("err", err.Error())
			}

			redirUrl.RawQuery = q.Encode()

			w.Header().Set("Location", redirUrl.String())
			w.WriteHeader(http.StatusSeeOther)
		}

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
