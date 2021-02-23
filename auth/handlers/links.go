package handlers

type (
	Links struct {
		Profile,
		Signup,
		ConfirmEmail,
		PendingEmailConfirmation,
		Login,
		ChangePassword,
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

		External,

		Assets string
	}
)

func GetLinks() Links {
	return Links{
		Profile:                  "/auth",
		Signup:                   "/auth/signup",
		ConfirmEmail:             "/auth/confirm-email",
		PendingEmailConfirmation: "/auth/pending-email-confirmation",
		Login:                    "/auth/login",
		ChangePassword:           "/auth/change-password",
		RequestPasswordReset:     "/auth/request-password-reset",
		PasswordResetRequested:   "/auth/password-reset-requested",
		ResetPassword:            "/auth/reset-password",
		Sessions:                 "/auth/sessions",
		AuthorizedClients:        "/auth/authorized-clients",
		Logout:                   "/auth/logout",

		OAuth2Authorize:       "/auth/oauth2/authorize",
		OAuth2AuthorizeClient: "/auth/oauth2/authorize-client",
		OAuth2Token:           "/auth/oauth2/token",
		OAuth2Info:            "/auth/oauth2/info",
		OAuth2DefaultClient:   "/auth/oauth2/default-client",

		External: "/auth/external",

		Assets: "/auth/assets/public",
	}
}
