package handlers

import (
	"context"
	"encoding/gob"
	"github.com/cortezaproject/corteza-server/auth/external"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/session"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/types"
	oauth2server "github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sort"
	"strings"
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

	AuthHandlers struct {
		Log *zap.Logger

		Templates      TemplateExecutor
		OAuth2         *oauth2server.Server
		SessionManager *session.Manager
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

func init() {
	gob.Register(&types.User{})
	gob.Register(&types.AuthClient{})
	gob.Register([]request.Alert{})
	gob.Register(url.Values{})
}

// Stores user & roles
//
// We need to store roles separately because they do not get serialized alongside with user
// due to unexported field
func (h *AuthHandlers) storeUserToSession(req *request.AuthReq, u *types.User) {
	session.SetUser(req.Session, u)
	session.SetRoleMemberships(req.Session, u.Roles())
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

			req.Client = session.GetOauth2Client(req.Session)

			req.User = session.GetUser(req.Session)

			// make sure user (identity) is part of the context
			// so we can properly identify ourselves when interacting
			// with services
			if req.User != nil {
				req.Request = req.Request.Clone(auth.SetIdentityToContext(req.Context(), req.User))
			}

			// Alerts show for 1 session only!
			req.PrevAlerts = req.PopAlerts()
			if err = fn(req); err != nil {
				h.Log.Error("error in handler", zap.Error(err))
				return
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
	req.Data["user"] = req.User

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
	dSettings.Providers = nil
	d["settings"] = dSettings

	providers := h.AuthService.GetProviders()
	sort.Sort(providers)

	var pp = make([]provider, 0, len(providers))
	for i := range providers {
		if !providers[i].Enabled {
			continue
		}

		p := provider{
			Label:  providers[i].Label,
			Handle: providers[i].Handle,
			Icon:   providers[i].Handle,
		}

		if strings.HasPrefix(p.Icon, external.OIDC_PROVIDER_PREFIX) {
			p.Icon = "key"
		}

		pp = append(pp, p)
	}

	d["providers"] = pp

	return d
}

// redirects anonymous users to login
func authOnly(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if req.User == nil {
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
		if req.User != nil {
			req.RedirectTo = GetLinks().Profile
			return nil
		} else {
			return fn(req)
		}
	}
}
