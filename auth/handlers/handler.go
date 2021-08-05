package handlers

import (
	"context"
	"encoding/gob"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza-server/auth/external"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"go.uber.org/zap"
)

type (
	authService interface {
		External(ctx context.Context, profile goth.User) (u *types.User, err error)
		InternalSignUp(ctx context.Context, input *types.User, password string) (u *types.User, err error)
		InternalLogin(ctx context.Context, email string, password string) (u *types.User, err error)
		SetPassword(ctx context.Context, userID uint64, password string) (err error)
		ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error)
		ValidateEmailConfirmationToken(ctx context.Context, token string) (user *types.User, err error)
		ValidatePasswordResetToken(ctx context.Context, token string) (user *types.User, err error)
		SendEmailAddressConfirmationToken(ctx context.Context, u *types.User) (err error)
		SendPasswordResetToken(ctx context.Context, email string) (err error)
		GetProviders() types.ExternalAuthProviderSet
		PasswordSet(ctx context.Context, email string) (is bool)
		ValidateTOTP(ctx context.Context, code string) (err error)
		ConfigureTOTP(ctx context.Context, secret string, code string) (u *types.User, err error)
		RemoveTOTP(ctx context.Context, userID uint64, code string) (u *types.User, err error)

		SendEmailOTP(ctx context.Context) (err error)
		ConfigureEmailOTP(ctx context.Context, userID uint64, enable bool) (u *types.User, err error)
		ValidateEmailOTP(ctx context.Context, code string) (err error)
	}

	userService interface {
		Update(context.Context, *types.User) (*types.User, error)
	}

	clientService interface {
		LookupByID(context.Context, uint64) (*types.AuthClient, error)
		Confirmed(context.Context, uint64) (types.AuthConfirmedClientSet, error)
		Revoke(ctx context.Context, userID, clientID uint64) error
	}

	// @todo this should probably be a little more decoupled from the store and nicely named
	tokenService interface {
		SearchByUserID(ctx context.Context, userID uint64) (types.AuthOa2tokenSet, error)
		DeleteByID(ctx context.Context, ID uint64) error
		DeleteByUserID(ctx context.Context, userID uint64) error
	}

	templateExecutor interface {
		ExecuteTemplate(io.Writer, string, interface{}) error
	}

	oauth2Service interface {
		GetRedirectURI(req *server.AuthorizeRequest, data map[string]interface{}) (string, error)
		CheckResponseType(rt oauth2.ResponseType) bool
		CheckCodeChallengeMethod(ccm oauth2.CodeChallengeMethod) bool
		ValidationAuthorizeRequest(r *http.Request) (*server.AuthorizeRequest, error)
		GetAuthorizeToken(ctx context.Context, req *server.AuthorizeRequest) (oauth2.TokenInfo, error)
		GetAuthorizeData(rt oauth2.ResponseType, ti oauth2.TokenInfo) map[string]interface{}
		HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error
		ValidationTokenRequest(r *http.Request) (oauth2.GrantType, *oauth2.TokenGenerateRequest, error)
		CheckGrantType(gt oauth2.GrantType) bool
		GetAccessToken(ctx context.Context, gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (oauth2.TokenInfo, error)
		GetTokenData(ti oauth2.TokenInfo) map[string]interface{}
		HandleTokenRequest(w http.ResponseWriter, r *http.Request) error
		GetErrorData(err error) (map[string]interface{}, int, http.Header)
		BearerAuth(r *http.Request) (string, bool)
		ValidationBearerToken(r *http.Request) (oauth2.TokenInfo, error)
	}

	AuthHandlers struct {
		Log *zap.Logger

		Templates      templateExecutor
		OAuth2         oauth2Service
		SessionManager *request.SessionManager
		AuthService    authService
		UserService    userService
		ClientService  clientService
		TokenService   tokenService
		DefaultClient  *types.AuthClient
		Opt            options.AuthOpt
		Settings       *settings.Settings
	}

	handlerFn func(req *request.AuthReq) error
)

const (
	TmplAuthorizedClients        = "authorized-clients.html.tpl"
	TmplChangePassword           = "change-password.html.tpl"
	TmplLogin                    = "login.html.tpl"
	TmplLogout                   = "logout.html.tpl"
	TmplOAuth2AuthorizeClient    = "oauth2-authorize-client.html.tpl"
	TmplRequestPasswordReset     = "request-password-reset.html.tpl"
	TmplPasswordResetRequested   = "password-reset-requested.html.tpl"
	TmplResetPassword            = "reset-password.html.tpl"
	TmplSecurity                 = "security.html.tpl"
	TmplProfile                  = "profile.html.tpl"
	TmplSessions                 = "sessions.html.tpl"
	TmplSignup                   = "signup.html.tpl"
	TmplPendingEmailConfirmation = "pending-email-confirmation.html.tpl"
	TmplMfa                      = "mfa.html.tpl"
	TmplMfaTotp                  = "mfa-totp.html.tpl"
	TmplMfaTotpDisable           = "mfa-totp-disable.html.tpl"
	TmplInternalError            = "error-internal.html.tpl"
)

func init() {
	gob.Register(&types.User{})
	gob.Register(&types.AuthClient{})
	gob.Register([]request.Alert{})
	gob.Register(url.Values{})
}

// handles auth request and prepares request struct with request, session and response helper
func (h *AuthHandlers) handle(fn handlerFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req = &request.AuthReq{
				Response:   w,
				Request:    r,
				Data:       make(map[string]interface{}),
				NewAlerts:  make([]request.Alert, 0),
				PrevAlerts: make([]request.Alert, 0),
				Session:    h.SessionManager.Get(r),
			}
		)

		h.Log.Debug(
			"handling request",
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method),
		)

		err := func() (err error) {
			if err = r.ParseForm(); err != nil {
				return
			}

			req.Client = request.GetOauth2Client(req.Session)

			req.AuthUser = request.GetAuthUser(req.Session)

			// make sure user (identity) is part of the context
			// so we can properly identify ourselves when interacting
			// with services
			if req.AuthUser != nil && !req.AuthUser.PendingMFA() {
				req.Request = req.Request.Clone(auth.SetIdentityToContext(req.Context(), req.AuthUser.User))
			}

			// Alerts show for 1 session only!
			req.PrevAlerts = req.PopAlerts()
			if err = fn(req); err != nil {
				h.Log.Error("error in handler", zap.Error(err))
				return
			}

			if req.RedirectTo != "" && len(req.PrevAlerts) > 0 {
				// redirect happened, so probably none noticed alerts
				// lets push them at the end of new alerts
				req.NewAlerts = append(req.NewAlerts, req.PrevAlerts...)
			}

			if len(req.NewAlerts) > 0 {
				req.SetAlerts(req.NewAlerts...)
			}

			if err = sessions.Save(r, w); err != nil {
				h.Log.Error("could not save session", zap.Error(err))
			}

			if req.Status == 0 {
				switch {
				case req.RedirectTo != "":
					req.Status = http.StatusSeeOther
					req.Template = ""
				case req.Template != "":
					req.Status = http.StatusOK
				default:
					req.Status = http.StatusInternalServerError
					req.Template = TmplInternalError
				}
			}

			return nil
		}()

		if err == nil {
			if req.Status >= 300 && req.Status < 400 {
				// redirect, nothing special to handle
				http.Redirect(w, r, req.RedirectTo, req.Status)
				return
			}

			if req.Status > 0 {
				// in cases when something else already wrote the status
				w.WriteHeader(req.Status)
			}
		}

		// Handling just text/html response types from here on
		//
		// If handler does not wish to use the template leave/set it to "" (empty string)

		if err == nil && req.Template != "" {
			err = h.Templates.ExecuteTemplate(w, req.Template, h.enrichTmplData(req))
			h.Log.Debug("template executed", zap.String("name", req.Template), zap.Error(err))

		}

		if err != nil {
			err = h.Templates.ExecuteTemplate(w, TmplInternalError, map[string]interface{}{
				"error": err,
			})

			if err == nil {
				return
			}
		}

		if err != nil {
			h.Log.Error("unhandled error", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Add alerts, settings, providers, csrf token
func (h *AuthHandlers) enrichTmplData(req *request.AuthReq) interface{} {
	d := req.Data
	if req.AuthUser != nil {
		req.Data["user"] = req.AuthUser.User
	}

	if req.Client != nil {
		c := authClient{
			ID:   req.Client.ID,
			Name: req.Client.Handle,
		}

		if req.Client.Meta != nil {
			c.Name = req.Client.Meta.Name
			c.Description = req.Client.Meta.Description
		}

		req.Data["client"] = c
	}

	d[csrf.TemplateTag] = csrf.TemplateField(req.Request)

	// In case we did not redirect, join previous alerts with new ones
	d["alerts"] = append(req.PrevAlerts, req.NewAlerts...)

	dSettings := *h.Settings

	var pp = make([]provider, 0, len(dSettings.Providers))
	for _, p := range dSettings.Providers {
		if _, err := goth.GetProvider(p.Handle); err != nil {
			continue
		}

		out := provider{
			Label:  p.Label,
			Handle: p.Handle,
			Icon:   p.Handle,
		}

		if strings.HasPrefix(out.Icon, external.OIDC_PROVIDER_PREFIX) {
			out.Icon = "key"
		}

		pp = append(pp, out)
	}

	d["providers"] = pp

	dSettings.Providers = nil
	d["settings"] = dSettings

	return d
}

// Handle successful auth (on any factor)
func handleSuccessfulAuth(req *request.AuthReq) {
	switch {
	case req.AuthUser.PendingMFA():
		req.RedirectTo = GetLinks().Mfa

	case request.GetOAuth2AuthParams(req.Session) != nil:
		// client authorization flow was paused, continue.
		req.RedirectTo = GetLinks().OAuth2AuthorizeClient

	default:
		// Always go to profile
		req.RedirectTo = GetLinks().Profile
	}
}

// redirects anonymous users to login
func authOnly(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		// these next few lines keep users away from the pages they should not see
		// and redirect them to where they need to be
		switch {
		case req.AuthUser == nil || req.AuthUser.User == nil:
			// not authenticated at all, move to login
			req.RedirectTo = GetLinks().Login

		case req.AuthUser.UnconfiguredTOTP():
			// authenticated but need to configure MFA
			req.RedirectTo = GetLinks().MfaTotpNewSecret

		case req.AuthUser.PendingMFA():
			// authenticated but MFA pending
			req.RedirectTo = GetLinks().Mfa

		default:
			return fn(req)

		}

		return nil
	}
}

func partAuthOnly(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if req.AuthUser == nil || req.AuthUser.User == nil {
			req.RedirectTo = GetLinks().Login
			return nil
		} else {
			return fn(req)
		}
	}
}

// redirects authenticated users to profile
func anonyOnly(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if req.AuthUser != nil && req.AuthUser.User != nil {
			req.RedirectTo = GetLinks().Profile
			return nil
		} else {
			return fn(req)
		}
	}
}
