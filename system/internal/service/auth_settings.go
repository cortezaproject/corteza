package service

type (
	authSettings struct {
		// Password reset path (<frontend password reset url> "?token=" + <token>)
		frontendUrlPasswordReset string

		// EmailAddress confirmation path (<frontend  email confirmation url> "?token=" + <token>)
		frontendUrlEmailConfirmation string

		mailFromAddress string
		mailFromName    string

		// Is external authentication
		externalEnabled bool

		// Is internal authentication (username + password) enabled
		internalEnabled bool

		// Can users register
		internalSignUpEnabled bool

		// Users should confirm their emails when signing-up
		internalSignUpEmailConfirmationRequired bool

		// Can users reset their passwords
		internalPasswordResetEnabled bool
	}

	authSettingsStore interface {
		Bool(string) bool
		String(string) string
	}
)

// AuthSettings maps from plain values to authSettings struct
//
// see settings.Initialize() func
func AuthSettings(kv authSettingsStore) authSettings {
	return authSettings{
		frontendUrlPasswordReset:     kv.String("auth.frontend.url.password-reset"),
		frontendUrlEmailConfirmation: kv.String("auth.frontend.url.email-confirmation"),

		mailFromAddress: kv.String("auth.mail.from-address"),
		mailFromName:    kv.String("auth.mail.from-name"),

		externalEnabled: kv.Bool("auth.external.enabled"),
		internalEnabled: kv.Bool("auth.internal.enabled"),

		internalSignUpEnabled:                   kv.Bool("auth.internal.sign-up.enabled"),
		internalSignUpEmailConfirmationRequired: kv.Bool("auth.internal.sign-up-email-confirmation-required.enabled"),

		internalPasswordResetEnabled: kv.Bool("auth.internal.password-reset.enabled"),
	}
}
