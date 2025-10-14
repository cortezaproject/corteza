package handlers

import (
	"github.com/cortezaproject/corteza/server/auth/request"
	"go.uber.org/zap"
)

func (h *AuthHandlers) securityForm(req *request.AuthReq) error {
	req.Template = TmplSecurity

	// user's MFA security policy
	umsp := req.AuthUser.User.Meta.SecurityPolicy.MFA

	// Check if email OTP or TOTP is enforced (either by user settings or globally)
	req.Data["emailOtpEnforced"] = umsp.EnforcedEmailOTP || h.Settings.MultiFactor.EmailOTP.Enforced
	req.Data["totpEnforced"] = umsp.EnforcedTOTP || h.Settings.MultiFactor.TOTP.Enforced

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

			if enable {
				req.NewAlerts = append(req.NewAlerts, request.Alert{
					Type: "primary",
					Text: t("security.template.mfa.email.enabled"),
				})
			} else {
				req.NewAlerts = append(req.NewAlerts, request.Alert{
					Type: "primary",
					Text: t("security.template.mfa.email.disabled"),
				})
			}

			// Make sure we update User's data in the session
			req.AuthUser.User = user
			req.AuthUser.Update(h.Settings, user)

			// If enabling, mark as complete in the current session; required for future logins
			if enable {
				req.AuthUser.CompleteEmailOTP()
			}

			req.AuthUser.Save(req.Session)

			h.Log.Info("email OTP configured", zap.Bool("enabled", enable))
		}
	}

	return nil
}
