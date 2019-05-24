package service

import (
	"strings"

	"github.com/markbates/goth"
)

type (
	authSettings struct {
		// Password reset path (<frontend password reset url> "?token=" + <token>)
		frontendUrlPasswordReset string

		// EmailAddress confirmation path (<frontend  email confirmation url> "?token=" + <token>)
		frontendUrlEmailConfirmation string

		// Where to redirect user after external auth flow
		frontendUrlRedirect string

		// Webapp Base URL
		frontendUrlBase string

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
		frontendUrlRedirect:          kv.String("auth.frontend.url.redirect"),
		frontendUrlBase:              kv.String("auth.frontend.url.base"),

		mailFromAddress: kv.String("auth.mail.from-address"),
		mailFromName:    kv.String("auth.mail.from-name"),

		externalEnabled: kv.Bool("auth.external.enabled"),
		internalEnabled: kv.Bool("auth.internal.enabled"),

		internalSignUpEnabled:                   kv.Bool("auth.internal.signup.enabled"),
		internalSignUpEmailConfirmationRequired: kv.Bool("auth.internal.signup-email-confirmation-required"),

		internalPasswordResetEnabled: kv.Bool("auth.internal.password-reset.enabled"),
	}
}

func (s authSettings) Format() map[string]interface{} {
	type (
		externalProvider struct {
			Label  string `json:"label"`
			Handle string `json:"handle"`
		}
	)

	var (
		providers = []externalProvider{}
	)

	for p := range goth.GetProviders() {
		label := p
		if strings.Index(p, "openid-connect.") == 0 {
			label = strings.SplitN(p, ".", 2)[1]
		}

		switch label {
		case "corteza-iam":
			label = "Corteza Unify"
		case "facebook":
			label = "Facebook"
		case "gplus":
			label = "Google"
		case "linkedin":
			label = "LinkedIn"
		case "github":
			label = "GitHub"
		}

		providers = append(providers, externalProvider{
			Label:  label,
			Handle: p,
		})
	}

	return map[string]interface{}{
		"internalEnabled":                         s.internalEnabled,
		"internalPasswordResetEnabled":            s.internalPasswordResetEnabled,
		"internalSignUpEmailConfirmationRequired": s.internalSignUpEmailConfirmationRequired,
		"internalSignUpEnabled":                   s.internalSignUpEnabled,

		"externalEnabled":   s.externalEnabled,
		"externalProviders": providers,
	}
}
