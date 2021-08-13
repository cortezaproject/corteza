package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"go.uber.org/zap"
)

func (h *AuthHandlers) securityForm(req *request.AuthReq) error {
	req.Template = TmplSecurity

	// user's MFA security policy
	umsp := req.AuthUser.User.Meta.SecurityPolicy.MFA

	req.Data["emailOtpEnforced"] = umsp.EnforcedEmailOTP
	req.Data["totpEnforced"] = umsp.EnforcedTOTP

	return nil
}

func (h *AuthHandlers) securityProc(req *request.AuthReq) error {
	req.RedirectTo = GetLinks().Security

	action := req.Request.Form.Get("action")
	switch action {
	case "reconfigureTOTP", "configureTOTP":
		// make sure secret is regenerated
		delete(req.Session.Values, totpSecretKey)

		req.RedirectTo = GetLinks().MfaTotpNewSecret
	case "disableTOTP":
		req.RedirectTo = GetLinks().MfaTotpDisable

	case "disableEmailOTP", "enableEmailOTP":
		enable := action == "enableEmailOTP"
		if user, err := h.AuthService.ConfigureEmailOTP(req.Context(), req.AuthUser.User.ID, enable); err != nil {
			return err
		} else {
			t := translator(req, "auth")
			req.NewAlerts = append(req.NewAlerts, request.Alert{
				Type: "primary",
				Text: t("security.TFA-TOTP-disabled"),
			})

			// Make sure we update User's data in the session
			req.AuthUser.User = user
			req.AuthUser.Save(req.Session)

			h.Log.Info("email OTP configured", zap.Bool("enabled", enable))
		}
	}

	return nil
}
