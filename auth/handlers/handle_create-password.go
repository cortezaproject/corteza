package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

var (
	loginLink          = GetLinks().Login
	profileLink        = GetLinks().Profile
	createPasswordLink = GetLinks().CreatePassword
)

func (h *AuthHandlers) createPasswordForm(req *request.AuthReq) (err error) {
	h.Log.Debug("password create form")
	req.Template = TmplCreatePassword

	// user not set, expecting valid token in URL
	if token := req.Request.URL.Query().Get("token"); len(token) > 0 {
		var (
			user *types.User
			// If user has password or not
			passwordSet bool
		)

		user, err = h.AuthService.ValidatePasswordCreateToken(req.Context(), token)
		if err == nil {
			// If user does not have password (or in case there is no user with such email)
			passwordSet = h.AuthService.PasswordSet(req.Context(), user.Email)
			if !passwordSet {
				// login user
				req.AuthUser = request.NewAuthUser(h.Settings, user, false)

				// redirect back to self (but without token and with user in session)
				h.Log.Debug("valid password create token found, refreshing page with stored user")
				req.RedirectTo = createPasswordLink
				req.AuthUser.Save(req.Session)

				return nil
			}
		}
	}

	if req.AuthUser == nil || err != nil {
		h.Log.Warn("invalid password create token request", zap.Error(err))
		h.invalidPasswordCreateAlert(req, loginLink)
		return nil
	}

	req.Data["form"] = req.PopKV()
	return nil
}

func (h *AuthHandlers) createPasswordProc(req *request.AuthReq) (err error) {
	h.Log.Debug("password create proc")

	err = h.AuthService.SetPassword(
		auth.SetIdentityToContext(req.Context(), req.AuthUser.User),
		req.AuthUser.User.ID,
		req.Request.PostFormValue("password"),
	)

	if err == nil {
		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("create-password.alerts.password-create-success"),
		})

		req.RedirectTo = profileLink
		return nil
	}

	switch {
	case service.AuthErrPasswordCreateDisabledByConfig().Is(err):
		h.passwordCreateDisabledAlert(req)

		h.Log.Warn("handled error", zap.Error(err))
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

func (h *AuthHandlers) onlyIfPasswordCreateEnabled(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if !h.Settings.PasswordCreateEnabled || !h.Settings.LocalEnabled {
			h.passwordCreateDisabledAlert(req)
			return nil
		}

		return fn(req)
	}
}

func (h *AuthHandlers) passwordCreateDisabledAlert(req *request.AuthReq) {
	req.RedirectTo = loginLink
	t := translator(req, "auth")
	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "danger",
		Text: t("create-password.alert.password-create-disabled"),
	})
}

func (h *AuthHandlers) invalidPasswordCreateAlert(req *request.AuthReq, redirectTo string) {
	req.RedirectTo = redirectTo
	t := translator(req, "auth")
	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "warning",
		Text: t("create-password.alerts.invalid-expired-password-token"),
	})
}
