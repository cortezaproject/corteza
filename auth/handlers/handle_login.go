package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
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
		user  *types.User
		email = req.Request.PostFormValue("email")
	)

	user, err = h.AuthService.InternalLogin(
		req.Context(),
		email,
		req.Request.PostFormValue("password"),
	)

	if err == nil {
		var (
			isPerm   = len(req.Request.PostFormValue("keep-session")) > 0
			lifetime = h.Opt.SessionLifetime
		)

		if isPerm {
			lifetime = h.Opt.SessionPermLifetime
		}

		req.AuthUser = request.NewAuthUser(h.Settings, user, isPerm, lifetime)

		req.AuthUser.Save(req.Session)

		h.Log.Info(
			"login with password successful",
			zap.Any("mfa", req.AuthUser.MFAStatus),
			zap.Bool("perm-login", isPerm),
			zap.Duration("lifetime", lifetime),
		)
		req.PushAlert("You are now logged-in")

		if req.AuthUser.PendingEmailOTP() {
			// Email OTP enforced (globally or by user sec. policy)
			//
			// @todo this should probably be part of login in auth service?
			if err = h.AuthService.SendEmailOTP(auth.SetIdentityToContext(req.Context(), req.AuthUser.User)); err != nil {
				return errors.Internal("could not send OTP via email, contact your administrator").Wrap(err)
			}
		}

		handleSuccessfulAuth(req)

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
	if req.AuthUser.User != nil {
		req.RedirectTo = GetLinks().Profile
	} else {
		req.RedirectTo = GetLinks().Login
	}

	req.PushDangerAlert("Local accounts disabled")
}
