package handlers

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
	"github.com/gorilla/csrf"
)

func (h *AuthHandlers) MountHttpRoutes(r chi.Router) {
	var (
		l = GetLinks()
	)

	if h.Opt.DevelopmentMode {
		r.Get("/auth/dev", h.handle(h.devView))
		r.Get("/auth/dev/scenarios", h.devSceneView)
	}

	r.Handle("/auth/", http.RedirectHandler("/auth", http.StatusSeeOther))
	r.Group(func(r chi.Router) {
		r.Use(locale.DetectLanguage(locale.Global()))

		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := actionlog.RequestOriginToContext(r.Context(), actionlog.RequestOrigin_Auth)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		if h.Opt.RequestRateLimit > 0 {
			r.Use(httprate.LimitByIP(h.Opt.RequestRateLimit, h.Opt.RequestRateWindowLength))
		}

		r.Use(request.ExtraReqInfoMiddleware)

		r.Group(func(r chi.Router) {
			// all routes protected with CSRF:
			if h.Opt.CsrfEnabled {
				r.Use(csrf.Protect(
					[]byte(h.Opt.CsrfSecret),
					csrf.SameSite(csrf.SameSiteStrictMode),
					csrf.Secure(h.Opt.SessionCookieSecure),
					csrf.FieldName(h.Opt.CsrfFieldName),
					csrf.CookieName(h.Opt.CsrfCookieName),
				))
			}

			r.Get(tbp(l.Profile), h.handle(authOnly(h.profileForm)))
			r.Post(tbp(l.Profile), h.handle(authOnly(h.profileProc)))

			r.HandleFunc(tbp(l.Logout), h.handle(h.logoutProc))

			r.Get(tbp(l.Sessions), h.handle(authOnly(h.sessionsView)))
			r.Post(tbp(l.Sessions), h.handle(authOnly(h.sessionsProc)))

			r.Get(tbp(l.AuthorizedClients), h.handle(authOnly(h.clientsView)))
			r.Post(tbp(l.AuthorizedClients), h.handle(authOnly(h.clientsProc)))

			r.Get(tbp(l.Signup), h.handle(h.onlyIfSignupEnabled(anonyOnly(h.signupForm))))
			r.Post(tbp(l.Signup), h.handle(h.onlyIfSignupEnabled(anonyOnly(h.signupProc))))
			r.Get(tbp(l.PendingEmailConfirmation), h.handle(h.pendingEmailConfirmation))
			r.Get(tbp(l.ConfirmEmail), h.handle(h.confirmEmail))

			r.Get(tbp(l.Login), h.handle(anonyOnly(h.loginForm)))
			r.Post(tbp(l.Login), h.handle(h.onlyIfLocalEnabled(anonyOnly(h.loginProc))))

			r.Get(tbp(l.Mfa), h.handle(h.mfaForm))
			r.Post(tbp(l.Mfa), h.handle(h.mfaProc))

			r.Get(tbp(l.RequestPasswordReset), h.handle(h.onlyIfPasswordResetEnabled(anonyOnly(h.requestPasswordResetForm))))
			r.Post(tbp(l.RequestPasswordReset), h.handle(h.onlyIfPasswordResetEnabled(anonyOnly(h.requestPasswordResetProc))))
			r.Get(tbp(l.PasswordResetRequested), h.handle(h.onlyIfPasswordResetEnabled(anonyOnly(h.passwordResetRequested))))
			r.Get(tbp(l.ResetPassword), h.handle(h.onlyIfPasswordResetEnabled(h.resetPasswordForm)))
			r.Post(tbp(l.ResetPassword), h.handle(h.onlyIfPasswordResetEnabled(authOnly(h.resetPasswordProc))))

			r.Get(tbp(l.Security), h.handle(authOnly(h.securityForm)))
			r.Post(tbp(l.Security), h.handle(authOnly(h.securityProc)))
			r.Get(tbp(l.ChangePassword), h.handle(h.onlyIfLocalEnabled(authOnly(h.changePasswordForm))))
			r.Post(tbp(l.ChangePassword), h.handle(h.onlyIfLocalEnabled(authOnly(h.changePasswordProc))))
			r.Get(tbp(l.CreatePassword), h.handle(h.onlyIfPasswordCreateEnabled(h.createPasswordForm)))
			r.Post(tbp(l.CreatePassword), h.handle(h.onlyIfPasswordCreateEnabled(h.createPasswordProc)))

			r.Get(tbp(l.MfaTotpNewSecret), h.handle(partAuthOnly(h.mfaTotpConfigForm)))
			r.Post(tbp(l.MfaTotpNewSecret), h.handle(partAuthOnly(h.mfaTotpConfigProc)))
			r.Get(tbp(l.MfaTotpQRImage), h.handle(partAuthOnly(h.mfaTotpConfigQR)))
			r.Get(tbp(l.MfaTotpDisable), h.handle(authOnly(h.mfaTotpDisableForm)))
			r.Post(tbp(l.MfaTotpDisable), h.handle(authOnly(h.mfaTotpDisableProc)))

		})

		r.Group(func(r chi.Router) {
			// OAuth2 routes
			r.HandleFunc(tbp(l.OAuth2Authorize), h.handle(h.oauth2Authorize))
			r.Get(tbp(l.OAuth2AuthorizeClient), h.handle(authOnly(h.oauth2AuthorizeClient)))
			r.Post(tbp(l.OAuth2AuthorizeClient), h.handle(authOnly(h.oauth2AuthorizeClientProc)))
			r.Get(tbp(l.OAuth2DefaultClient), h.handle(h.oauth2authorizeDefaultClient))
			r.Post(tbp(l.OAuth2DefaultClient), h.handle(h.oauth2authorizeDefaultClientProc))
		})

		// Wrapping SAML structs so we assure that fresh ones are always used in case
		// of settings changes.
		//
		// @todo refactor this with a wrapping struct to handle serve HTTPS and pass
		// calls to internal SAML service.
		r.Group(func(r chi.Router) {
			r.HandleFunc(tbp(l.SamlMetadata), func(rw http.ResponseWriter, r *http.Request) {
				if h.SamlSPService == nil || !h.SamlSPService.Enabled {
					rw.WriteHeader(http.StatusServiceUnavailable)
					return
				}

				h.SamlSPService.ServeHTTP(rw, r)
			})
			r.HandleFunc(tbp(l.SamlCallback), func(rw http.ResponseWriter, r *http.Request) {
				if h.SamlSPService == nil || !h.SamlSPService.Enabled {
					rw.WriteHeader(http.StatusServiceUnavailable)
					return
				}

				h.SamlSPService.ServeHTTP(rw, r)
			})
			r.HandleFunc(tbp(l.SamlInit), func(rw http.ResponseWriter, r *http.Request) {
				h.samlInit(rw, r)
			})
		})

		r.Route(tbp(l.External)+"/{provider}", func(r chi.Router) {
			// External provider
			r.Get("/", h.externalInit)
			r.Get("/callback", h.externalCallback)
		})

		r.HandleFunc("/auth/oauth2/token", h.handle(h.oauth2Token))
		r.HandleFunc("/auth/oauth2/info", h.oauth2Info)
		r.HandleFunc("/auth/oauth2/public-keys", h.oauth2PublicKeys)
	})
}
