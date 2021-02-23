package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/session"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	provider struct {
		Label, Handle, Icon string
	}
)

func (h *AuthHandlers) loginForm(req *request.AuthReq) error {
	req.Template = TmplLogin
	req.Data["form"] = req.GetKV()
	return nil
}

func (h *AuthHandlers) loginProc(req *request.AuthReq) (err error) {
	req.RedirectTo = GetLinks().Login
	req.SetKV(nil)

	var (
		authUser *types.User
		email    = req.Request.PostFormValue("email")
	)

	authUser, err = h.AuthService.InternalLogin(
		req.Context(),
		email,
		req.Request.PostFormValue("password"),
	)

	if err == nil {
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "You are now logged-in",
		})

		if session.GetOAuth2AuthParams(req.Session) == nil {
			// Not in the OAuth2 flow, go to profile
			req.RedirectTo = GetLinks().Profile
		} else {
			h.Log.Info("oauth2 params found, continuing with authorization flow")
			req.RedirectTo = GetLinks().OAuth2AuthorizeClient
		}

		h.Log.Info("login successful")
		h.storeUserToSession(req, authUser)

		if len(req.Request.PostFormValue("keep-session")) > 0 {
			session.SetPerm(req.Session, h.Opt.SessionPermLifetime)
		}

		return nil
	}

	switch {
	case service.AuthErrInteralLoginDisabledByConfig().Is(err):
		h.localDisabledAlert(req)
		return nil
	case service.AuthErrInvalidEmailFormat().Is(err),
		service.AuthErrInvalidCredentials().Is(err),
		service.AuthErrCredentialsLinkedToInvalidUser().Is(err):
		req.SetKV(map[string]string{
			"error": err.Error(),
			"email": email,
		})

		h.Log.Warn("handled error", zap.Error(err))
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

func (h *AuthHandlers) onlyIfLocalEnabled(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if !h.Settings.LocalEnabled {
			h.localDisabledAlert(req)
			return nil
		}

		return fn(req)
	}
}

func (h *AuthHandlers) localDisabledAlert(req *request.AuthReq) {
	if req.User != nil {
		req.RedirectTo = GetLinks().Profile
	} else {
		req.RedirectTo = GetLinks().Login
	}

	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "danger",
		Text: "Local accounts disabled",
	})
}
