package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/service"
	"go.uber.org/zap"
)

func (h *AuthHandlers) changePasswordForm(req *request.AuthReq) error {
	h.Log.Debug("showing password change form")
	req.Template = TmplChangePassword
	req.Data["form"] = req.PopKV()
	return nil
}

func (h *AuthHandlers) changePasswordProc(req *request.AuthReq) (err error) {
	err = h.AuthService.ChangePassword(
		auth.SetIdentityToContext(req.Context(), req.AuthUser.User),
		req.AuthUser.User.ID,
		req.Request.PostFormValue("oldPassword"),
		req.Request.PostFormValue("newPassword"),
	)

	t := translator(req, "auth")

	if err == nil {
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("change-password.alerts.text"),
		})

		req.RedirectTo = GetLinks().Profile
		return nil
	}

	switch {
	case service.AuthErrInternalLoginDisabledByConfig().Is(err),
		service.AuthErrPasswordNotSecure().Is(err),
		service.AuthErrPasswordChangeFailedForUnknownUser().Is(err),
		service.AuthErrPasswodResetFailedOldPasswordCheckFailed().Is(err):
		req.SetKV(map[string]string{
			"error": err.Error(),
		})
		req.RedirectTo = GetLinks().ChangePassword

		h.Log.Warn("handled error", zap.Error(err))
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}
