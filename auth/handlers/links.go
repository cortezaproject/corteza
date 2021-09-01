package handlers

import "strings"

type (
	Links struct {
		Profile,
		Signup,
		ConfirmEmail,
		PendingEmailConfirmation,
		Login,

		Security,
		ChangePassword,
		CreatePassword,

		RequestPasswordReset,
		PasswordResetRequested,
		ResetPassword,
		Sessions,
		AuthorizedClients,
		Logout,

		OAuth2Authorize,
		OAuth2AuthorizeClient,
		OAuth2Token,
		OAuth2Info,
		OAuth2DefaultClient,

		Mfa,

		MfaTotpNewSecret,
		MfaTotpQRImage,
		MfaTotpDisable,

		External,

		SamlInit,
		SamlCallback,
		SamlMetadata,
		SamlLogout,

		Base,

		Assets string
	}
)

var BasePath string = "/"

func GetLinks() Links {
	var b = strings.TrimSuffix(BasePath, "/") + "/"

	return Links{
		Profile:                  b + "auth",
		Signup:                   b + "auth/signup",
		ConfirmEmail:             b + "auth/confirm-email",
		PendingEmailConfirmation: b + "auth/pending-email-confirmation",
		Login:                    b + "auth/login",
		Security:                 b + "auth/security",
		ChangePassword:           b + "auth/change-password",
		CreatePassword:           b + "auth/create-password",
		RequestPasswordReset:     b + "auth/request-password-reset",
		PasswordResetRequested:   b + "auth/password-reset-requested",
		ResetPassword:            b + "auth/reset-password",
		Sessions:                 b + "auth/sessions",
		AuthorizedClients:        b + "auth/authorized-clients",
		Logout:                   b + "auth/logout",

		OAuth2Authorize:       b + "auth/oauth2/authorize",
		OAuth2AuthorizeClient: b + "auth/oauth2/authorize-client",
		OAuth2Token:           b + "auth/oauth2/token",
		OAuth2Info:            b + "auth/oauth2/info",
		OAuth2DefaultClient:   b + "auth/oauth2/default-client",

		Mfa:              b + "auth/mfa",
		MfaTotpNewSecret: b + "auth/mfa/totp/setup",
		MfaTotpQRImage:   b + "auth/mfa/totp/qr.png",
		MfaTotpDisable:   b + "auth/mfa/totp/disable",

		External: b + "auth/external",

		SamlInit:     b + "auth/external/saml/init",
		SamlCallback: b + "auth/external/saml/callback",
		SamlMetadata: b + "auth/external/saml/metadata",
		SamlLogout:   b + "auth/external/saml/slo",

		Assets: b + "auth/assets/public",
		Base: b,
	}
}

// trim base path
func tbp(s string) string {
	s = strings.TrimPrefix(s, BasePath)
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}

	return s
}
