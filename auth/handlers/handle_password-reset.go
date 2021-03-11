package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/service"
	"go.uber.org/zap"
)

func (h *AuthHandlers) requestPasswordResetForm(req *request.AuthReq) error {
	h.Log.Debug("showing request password reset form")
	req.Template = TmplRequestPasswordReset
	req.Data["form"] = req.GetKV()
	return nil
}

func (h *AuthHandlers) requestPasswordResetProc(req *request.AuthReq) (err error) {
	h.Log.Debug("processing password change request")

	email := req.Request.PostFormValue("email")
	err = h.AuthService.SendPasswordResetToken(req.Context(), email)

	if err == nil || errors.IsNotFound(err) {
		req.RedirectTo = GetLinks().PasswordResetRequested
		return nil
	}

	switch {
	case service.AuthErrPasswordResetDisabledByConfig().Is(err):
		h.passwordResetDisabledAlert(req)
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

func (h *AuthHandlers) passwordResetRequested(req *request.AuthReq) error {
	req.Template = TmplPasswordResetRequested
	return nil
}

func (h *AuthHandlers) resetPasswordForm(req *request.AuthReq) (err error) {
	h.Log.Debug("password reset form")

	req.Template = TmplResetPassword

	if req.AuthUser.User == nil {
		// user not set, expecting valid token in URL
		if token := req.Request.URL.Query().Get("token"); len(token) > 0 {
			req.AuthUser.User, err = h.AuthService.ValidatePasswordResetToken(req.Context(), token)
			if err == nil {
				// redirect back to self (but without token and with user in session
				h.Log.Debug("valid password reset token found, refreshing page with stored user")
				req.RedirectTo = GetLinks().ResetPassword
				req.AuthUser.Save(req.Session)
				return nil
			}
		}

		h.Log.Warn("invalid password reset token used", zap.Error(err))
		req.RedirectTo = GetLinks().RequestPasswordReset
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "warning",
			Text: "Invalid or expired password reset token, please repeat password reset request.",
		})
	}

	req.Data["form"] = req.GetKV()
	return nil
}

func (h *AuthHandlers) resetPasswordProc(req *request.AuthReq) (err error) {
	h.Log.Debug("password reset proc")

	err = h.AuthService.SetPassword(req.Context(), req.AuthUser.User.ID, req.Request.PostFormValue("password"))

	if err == nil {
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "Password successfully reset.",
		})

		req.RedirectTo = GetLinks().Profile
		return nil
	}

	switch {
	case service.AuthErrPasswordResetDisabledByConfig().Is(err):
		h.passwordResetDisabledAlert(req)
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

func (h *AuthHandlers) onlyIfPasswordResetEnabled(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if !h.Settings.PasswordResetEnabled || !h.Settings.LocalEnabled {
			h.passwordResetDisabledAlert(req)
			return nil
		}

		return fn(req)
	}
}

func (h *AuthHandlers) passwordResetDisabledAlert(req *request.AuthReq) {
	req.RedirectTo = GetLinks().Login
	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "danger",
		Text: "Password reset disabled",
	})
}
