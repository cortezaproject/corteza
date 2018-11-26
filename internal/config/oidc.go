package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	OIDC struct {
		Enabled      bool
		Issuer       string
		ClientID     string
		ClientSecret string

		RedirectURL string
		AppURL      string

		StateCookieExpiry int64
	}
)

var oidc *OIDC

func (c *OIDC) Validate() error {
	if c == nil {
		return nil
	}
	if c.Enabled == false {
		return nil
	}
	if c.Issuer == "" {
		return errors.New("OIDC Issuer not set for AUTH")
	}
	if c.RedirectURL == "" {
		return errors.New("OIDC RedirectURL not set for AUTH")
	}

	return nil
}

func (*OIDC) Init(prefix ...string) *OIDC {
	if oidc != nil {
		return oidc
	}

	oidc = new(OIDC)
	flag.BoolVar(&oidc.Enabled, "auth-oidc-enabled", true, "OIDC enabled")
	flag.StringVar(&oidc.Issuer, "auth-oidc-issuer", "", "OIDC Issuer")
	flag.StringVar(&oidc.ClientID, "auth-oidc-client-id", "", "OIDC Client ID")
	flag.StringVar(&oidc.ClientSecret, "auth-oidc-client-secret", "", "OIDC Client Secret")
	flag.StringVar(&oidc.RedirectURL, "auth-oidc-redirect-url", "", "OIDC RedirectURL")
	flag.StringVar(&oidc.AppURL, "auth-oidc-app-url", "", "OIDC AppURL")
	flag.Int64Var(&oidc.StateCookieExpiry, "auth-oidc-state-cookie-expiry", 15, "OIDC State cookie expiry in minutes")
	return oidc
}
