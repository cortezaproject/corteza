package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	intset "github.com/cortezaproject/corteza-server/internal/settings"
)

type (
	AuthSettings struct {
		// Password reset path (<frontend password reset url> "?token=" + <token>)
		FrontendUrlPasswordReset string

		// EmailAddress confirmation path (<frontend  email confirmation url> "?token=" + <token>)
		FrontendUrlEmailConfirmation string

		// Where to redirect user after external auth flow
		FrontendUrlRedirect string

		// Webapp Base URL
		FrontendUrlBase string

		MailFromAddress string
		MailFromName    string

		// Is internal authentication (username + password) enabled
		InternalEnabled bool

		// Can users register
		InternalSignUpEnabled bool

		// Users should confirm their emails when signing-up
		InternalSignUpEmailConfirmationRequired bool

		// Can users reset their passwords
		InternalPasswordResetEnabled bool

		// Is external authentication
		ExternalEnabled bool

		// Where to redirect (url used for registration)
		ExternalRedirectUrl string

		// session secret to use
		ExternalSessionStoreSecret string

		// session store should be secure
		ExternalSessionStoreSecure bool

		// all external providers we know
		ExternalProviders map[string]AuthSettingsExternalAuthProvider
	}

	AuthSettingsExternalAuthProvider struct {
		Enabled     bool
		Key         string
		Secret      string
		RedirectUrl string
		IssuerUrl   string
	}
)

// ParseAuthSettings maps from plain values to AuthSettings struct
//
// see settings.Initialize() func
func ParseAuthSettings(kv intset.KV) (as AuthSettings, err error) {
	as = AuthSettings{
		FrontendUrlPasswordReset:     kv.String("auth.frontend.url.password-reset"),
		FrontendUrlEmailConfirmation: kv.String("auth.frontend.url.email-confirmation"),
		FrontendUrlRedirect:          kv.String("auth.frontend.url.redirect"),
		FrontendUrlBase:              kv.String("auth.frontend.url.base"),

		MailFromAddress: kv.String("auth.mail.from-address"),
		MailFromName:    kv.String("auth.mail.from-name"),

		InternalEnabled: kv.Bool("auth.internal.enabled"),

		InternalSignUpEnabled:                   kv.Bool("auth.internal.signup.enabled"),
		InternalSignUpEmailConfirmationRequired: kv.Bool("auth.internal.signup-email-confirmation-required"),

		InternalPasswordResetEnabled: kv.Bool("auth.internal.password-reset.enabled"),

		ExternalEnabled: kv.Bool("auth.external.enabled"),

		ExternalRedirectUrl:        kv.String("auth.external.redirect-url"),
		ExternalSessionStoreSecret: kv.String("auth.external.session-store-secret"),
		ExternalSessionStoreSecure: kv.Bool("auth.external.session-store-secure"),
	}

	as.ExternalProviders, err = as.parseExternalProviders(kv)

	return
}

func (as *AuthSettings) parseExternalProviders(kv intset.KV) (map[string]AuthSettingsExternalAuthProvider, error) {
	// Standard providers:
	var (
		ep = map[string]AuthSettingsExternalAuthProvider{
			"github":   {},
			"facebook": {},
			"gplus":    {},
			"linkedin": {},
		}

		// Add all oidc providers we find
		extKeyBase  = "auth.external.providers."
		oidcKeyBase = extKeyBase + "openid-connect."
	)

	for k := range kv.Filter(oidcKeyBase) {
		if len(k) < len(oidcKeyBase)+2 {
			// skip invalid keys
			continue
		}

		// find next dot:
		name := k[len(oidcKeyBase):]
		dotPos := strings.Index(name, ".")
		if dotPos > 0 {
			name = name[:dotPos]
		}

		ep["openid-connect."+name] = AuthSettingsExternalAuthProvider{}
	}

	for provider := range ep {
		if p, err := as.parseExternalProvider(kv.Filter(extKeyBase + provider)); err != nil {
			return nil, err
		} else {
			if as.ExternalRedirectUrl != "" && p.Enabled {
				p.RedirectUrl = fmt.Sprintf(as.ExternalRedirectUrl, provider)
			}

			ep[provider] = *p
		}

	}

	return ep, nil
}

// Parses external provider out of KV set
//
// Function only looks at the end of key string (after last dot)
// so passing multiple providers will result in overriding values
func (as *AuthSettings) parseExternalProvider(kv intset.KV) (p *AuthSettingsExternalAuthProvider, err error) {
	p = &AuthSettingsExternalAuthProvider{}

	for k, v := range kv {
		ld := strings.LastIndex(k, ".")

		switch k[ld+1:] {
		case "enabled":
			err = v.Unmarshal(&p.Enabled)
		case "key":
			err = v.Unmarshal(&p.Key)
		case "secret":
			err = v.Unmarshal(&p.Secret)
		case "issuer":
			err = v.Unmarshal(&p.IssuerUrl)
		}

		if err != nil {
			return
		}
	}

	return
}

func (as AuthSettings) Format() map[string]interface{} {
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
		case "corteza-iam", "corteza", "corteza-one":
			label = "Corteza One"
		case "crust-iam", "crust", "crust-unify":
			label = "Crust Unify"
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
		"internalEnabled":                         as.InternalEnabled,
		"internalPasswordResetEnabled":            as.InternalPasswordResetEnabled,
		"internalSignUpEmailConfirmationRequired": as.InternalSignUpEmailConfirmationRequired,
		"internalSignUpEnabled":                   as.InternalSignUpEnabled,

		"externalEnabled":   as.ExternalEnabled,
		"externalProviders": providers,
	}
}

// StaticValidateExternal
//
// Simple checks of external auth settings
func (as AuthSettings) StaticValidateExternal() error {
	if as.ExternalRedirectUrl == "" {
		return errors.New("redirect URL is empty")
	}

	const (
		tpt = "test-provider-test"
	)
	p, err := url.Parse(fmt.Sprintf(as.ExternalRedirectUrl, tpt))
	if err != nil {
		return errors.Wrap(err, "invalid redirect URL")
	}

	if !strings.Contains(p.Path, tpt+"/callback") {
		return errors.Wrap(err, "could find injected provider in the URL, make sure you use '%s' as a placeholder")
	}

	if as.ExternalSessionStoreSecret == "" {
		return errors.New("session store secret is empty")
	}

	if as.ExternalSessionStoreSecure && p.Scheme != "https" {
		return errors.New("session store is secure, redirect URL should have HTTPS")
	}

	return nil
}

// ValidateExternalRedirectURL
//
// Validates external redirect URL
func (as AuthSettings) ValidateExternalRedirectURL() error {
	const tpt = "test-provider-test"
	const cb = "/callback"

	// Replace placeholders & remove /callback
	var url = fmt.Sprintf(as.ExternalRedirectUrl, tpt)
	url = url[0 : len(url)-len(cb)]

	rsp, err := http.DefaultClient.Get(url)
	if err != nil {
		return errors.Wrap(err, "could not get response from redirect URL")
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	if strings.Contains(string(body), tpt) {
		return nil
	}

	return errors.New("could not validate external auth redirection URL")
}

func (p AuthSettingsExternalAuthProvider) MakeValueSet(name string) (vv intset.ValueSet, err error) {
	set := func(name string, value interface{}) error {
		v := &intset.Value{Name: name}
		if v.Value, err = json.Marshal(value); err != nil {
			return err
		}

		vv = append(vv, v)
		return nil
	}

	prefix := "auth.external.providers." + name

	if err = set(prefix+".enabled", p.Enabled); err != nil {
		return nil, err
	}

	if err = set(prefix+".key", p.Key); err != nil {
		return nil, err
	}

	if err = set(prefix+".secret", p.Secret); err != nil {
		return nil, err
	}

	if err = set(prefix+".issuer", p.IssuerUrl); err != nil {
		return nil, err
	}

	return vv, err
}
