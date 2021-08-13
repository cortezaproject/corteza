package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

// Handles MFA TOTP configuration form
//
// Where the TOTP QR & code are displayed and where
func (h AuthHandlers) mfaForm(req *request.AuthReq) (err error) {
	req.Template = TmplMfa

	if !req.AuthUser.PendingMFA() {
		req.RedirectTo = GetLinks().Profile
		return nil
	}

	if req.Request.URL.Query().Get("action") == "resendEmailOtp" {
		err = h.AuthService.SendEmailOTP(
			auth.SetIdentityToContext(req.Context(), req.AuthUser.User),
		)

		if err != nil {
			req.SetKV(map[string]string{"emailOtpError": err.Error()})
			return nil
		}

		t := translator(req, "auth")

		req.PushAlert(t("mfa.handle.email-resent"))
		req.RedirectTo = GetLinks().Mfa
	}

	req.Data["form"] = req.PopKV()
	req.Data["emailOtpDisabled"] = req.AuthUser.DisabledEmailOTP()
	req.Data["emailOtpPending"] = req.AuthUser.PendingEmailOTP()
	req.Data["totpDisabled"] = req.AuthUser.DisabledTOTP()
	req.Data["totpPending"] = req.AuthUser.PendingTOTP()
	return nil
}

// Handles MFA OTP form processing
func (h AuthHandlers) mfaProc(req *request.AuthReq) (err error) {
	req.RedirectTo = GetLinks().Mfa
	req.SetKV(nil)

	switch req.Request.Form.Get("action") {
	case "verifyEmailOtp":
		err = h.AuthService.ValidateEmailOTP(
			auth.SetIdentityToContext(req.Context(), req.AuthUser.User),
			req.Request.PostFormValue("code"),
		)

		if err != nil {
			req.SetKV(map[string]string{"emailOtpError": err.Error()})
			return nil
		}

		t := translator(req, "auth")

		req.PushAlert(t("mfa.handle.email-resent"))
		req.AuthUser.CompleteEmailOTP()

	case "verifyTotp":
		err = h.AuthService.ValidateTOTP(
			auth.SetIdentityToContext(req.Context(), req.AuthUser.User),
			req.Request.PostFormValue("code"),
		)

		if err != nil {
			req.SetKV(map[string]string{"totpError": err.Error()})
			return nil
		}

		t := translator(req, "auth")

		req.PushAlert(t("mfa.handle.topt-valid"))
		req.AuthUser.CompleteTOTP()
	}

	// All required MFA's confirmed, proceed to profile
	handleSuccessfulAuth(req)

	return nil
}
