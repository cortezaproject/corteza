package handlers

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

type (
	provider struct {
		Label, Handle, Icon string
	}
)

func (h *AuthHandlers) loginForm(req *request.AuthReq) error {
	req.Template = TmplLogin
	kv := req.PopKV()

	if kv == nil && h.Settings.SplitCredentialsCheck {
		// Force login form to show only email input
		//
		// KV is nil (means that this is first load of the form)
		// and credentials check is split into two parts (email first then credentials)
		kv = map[string]string{
			"splitCredentialsCheck": "split",
		}
	}

	req.Data["form"] = kv
	req.Data["enableRememberMe"] = h.Opt.SessionPermLifetime > 0
	return nil
}

func (h *AuthHandlers) loginProc(req *request.AuthReq) (err error) {
	// In most cases, we want to redirect back to login
	req.RedirectTo = GetLinks().Login
	req.SetKV(nil)

	var (
		user     *types.User
		email    = req.Request.PostFormValue("email")
		password = req.Request.PostFormValue("password")
	)

	err = func() (err error) {
		if len(email) > 0 && len(password) == 0 && h.Settings.SplitCredentialsCheck {
			// Email provided but no password and the
			// split credentials check enabled
			//
			// SetKV w/ email will prevent the login form to show only email input
			req.SetKV(map[string]string{"email": email})

			// If user does not have password set (or in case there is no user with such email)
			// and there is exactly one IdP we redirect user to that IdP
			if !h.AuthService.PasswordSet(req.Context(), email) && len(h.Settings.Providers) == 1 {
				// User w/o the password
				// In case there is one single IdP automatically redirect user there
				req.RedirectTo = fmt.Sprintf("%s/%s", GetLinks().External, h.Settings.Providers[0].Handle)

				// for the clarity of the flow
				// keeping this optional return here
				return
			}

			// This user has existing password credentials or with more than one IdP
			//
			// Take user back to the login form and ask for the password
			// or any other kind of login
			return
		}

		user, err = h.AuthService.InternalLogin(req.Context(), email, password)

		if err != nil {
			return
		}

		var (
			isPerm   = len(req.Request.PostFormValue("keep-session")) > 0
			lifetime = h.Opt.SessionLifetime
		)

		if isPerm {
			lifetime = h.Opt.SessionPermLifetime
		}

		req.AuthUser = request.NewAuthUser(h.Settings, user, isPerm)

		req.AuthUser.Save(req.Session)

		h.Log.Info(
			"login with password successful",
			zap.Any("mfa", req.AuthUser.MFAStatus),
			zap.Bool("perm-login", isPerm),
			zap.Duration("lifetime", lifetime),
		)

		t := translator(req, "auth")
		req.PushAlert(t("login.alerts.logged-in"))

		if req.AuthUser.PendingEmailOTP() {
			// Email OTP enforced (globally or by user sec. policy)
			//
			// @todo this should probably be part of login in auth service?
			if err = h.AuthService.SendEmailOTP(auth.SetIdentityToContext(req.Context(), req.AuthUser.User)); err != nil {
				return errors.Internal("could not send OTP via email, contact your administrator").Wrap(err)
			}
		}

		handleSuccessfulAuth(req)

		return
	}()

	if err == nil {
		return nil
	}

	switch {
	case service.AuthErrInternalLoginDisabledByConfig().Is(err):
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

	t := translator(req, "auth")

	req.PushDangerAlert(t("login.alert.local-disabled"))
}
