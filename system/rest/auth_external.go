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
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	ExternalAuth struct {
		auth       service.AuthService
		jwtEncoder auth.TokenEncoder
	}
)

const (
	externalAuthBaseUrl = "/auth/external"
	redirCookieName     = "redir"
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
	// Copy provider from path (Chi URL param) to request context and return it
	copyProviderToContext := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), "provider", chi.URLParam(r, "provider")))
	}

	r.Route(externalAuthBaseUrl+"/{provider}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			r = copyProviderToContext(r)

			// Always set redir cookie, even if not requested.
			// If param is empty, cookie will be removed
			ctrl.setSessionCookie(w, r, redirCookieName, r.URL.Query().Get("redir"))

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

	var (
		u   *types.User
		err error
	)

	// Try to login/sign-up external user
	if u, err = svc.External(cred); err != nil {
		resputil.JSON(w, err)
		return
	}

	var (
		ctx      = r.Context()
		token    string
		redirUrl *url.URL
		c        *http.Cookie
	)

	if c, err = r.Cookie(redirCookieName); err != nil && err != http.ErrNoCookie {
		ctrl.log(ctx, zap.Error(err)).Warn("error reading cookies")
	}

	if c != nil {
		// Remove cookie
		ctrl.setSessionCookie(w, r, redirCookieName, "")
		redirUrl, err = url.Parse(c.Value)

		if redirUrl == nil {
			ctrl.log(ctx, zap.Error(err)).Warn("failed to parse URL from redir cookie")
		}
	}

	if redirUrl == nil {
		// Try with frontend redirect URL
		if fru := svc.FrontendRedirectURL(); fru != "" {
			redirUrl, err = url.Parse(fru)

			if redirUrl == nil {
				ctrl.log(ctx, zap.Error(err)).Warn("failed to parse URL from 'auth.frontend.url.redirect' settings")
			}
		} else {
			// No info about where should we be redirected.
			// Let's go directly to /auth and append the token
			redirUrl = r.URL
			redirUrl.RawQuery = ""
			search := "/auth/"
			p := strings.Index(redirUrl.Path, search)
			if p > -1 {
				redirUrl.Path = redirUrl.Path[0 : p+len(search)]
			} else {
				redirUrl.Path = "/"
			}
		}
	}

	if redirUrl != nil {
		q := redirUrl.Query()

		if u != nil {
			// Append auth request token to the URL
			// This token is used by the client and exchanged for JWT
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
