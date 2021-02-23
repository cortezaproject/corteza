package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"go.uber.org/zap"
)

func (h *AuthHandlers) changePasswordForm(req *request.AuthReq) error {
	h.Log.Debug("showing password change form")
	req.Template = TmplChangePassword
	req.Data["form"] = req.GetKV()
	return nil
}

func (h *AuthHandlers) changePasswordProc(req *request.AuthReq) (err error) {
	err = h.AuthService.ChangePassword(
		req.Context(),
		req.User.ID,
		req.Request.PostFormValue("oldPassword"),
		req.Request.PostFormValue("newPassword"),
	)

	if err == nil {
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "Password successfully changed.",
		})

		req.RedirectTo = GetLinks().Profile
		return nil
	}

	switch {
	case service.AuthErrInteralLoginDisabledByConfig().Is(err),
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
