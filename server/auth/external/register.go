package external

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/crusttech/go-oidc"
	"go.uber.org/zap"
)

func AddProvider(ctx context.Context, log *zap.Logger, s store.SettingValues, eap *types.ExternalAuthProvider, force bool) error {
	var (
		prefix = "auth.external.providers." + eap.Key + "."
	)

	log = log.With(
		zap.Bool("force", force),
		zap.String("handle", eap.Handle),
		zap.String("key", eap.Key),
	)

	if eap.IssuerUrl != "" {
		log = log.With(zap.String("issuer-url", eap.IssuerUrl))
	}

	log.Info("adding external authentication provider")

	ss, _, err := store.SearchSettingValues(ctx, s, types.SettingsFilter{
		Prefix: prefix,
	})

	if err != nil {
		return err
	}

	ex := ss.KV().CutPrefix(prefix)

	if !force {
		// check if exists before storing it
		if len(ex) > 0 && ex.String("key") == eap.Key && ex.String("secret") == eap.Secret {
			return nil
		}
	}
	if vv, err := eap.EncodeKV(); err != nil {
		return fmt.Errorf("could not encode auth provider values: %w", err)
	} else if err = store.UpsertSettingValue(ctx, s, vv...); err != nil {
		return fmt.Errorf("could not store auth provider values: %w", err)
	}

	log.Info("external authentication provider added")
	return nil
}

// @todo remove dependency on github.com/crusttech/go-oidc (and github.com/coreos/go-oidc)
//
//	and move client registration to corteza codebase
func DiscoverOidcProvider(ctx context.Context, log *zap.Logger, opt options.AuthOpt, name, url string) (eap *types.ExternalAuthProvider, err error) {
	var (
		provider    *oidc.Provider
		client      *oidc.Client
		redirectUrl = strings.Replace(opt.ExternalRedirectURL, "{provider}", OIDC_PROVIDER_PREFIX+name, 1)
	)

	log = log.With(
		zap.String("name", name),
		zap.String("url", url),
	)

	if provider, err = oidc.NewProvider(ctx, url); err != nil {
		return
	}

	client, err = provider.RegisterClient(ctx, &oidc.ClientRegistration{
		Name:          "Corteza",
		RedirectURIs:  []string{redirectUrl},
		ResponseTypes: []string{"token id_token", "code"},
	})

	if err != nil {
		log.Error("could not register oidc provider", zap.Error(err))
		return
	}

	eap = &types.ExternalAuthProvider{
		Handle:      OIDC_PROVIDER_PREFIX + name,
		RedirectUrl: redirectUrl,
		Key:         client.ID,
		Secret:      client.Secret,
		IssuerUrl:   url,
	}

	log.Info("oidc provider registered", zap.String("key", client.ID))

	return
}

func RegisterOidcProvider(ctx context.Context, log *zap.Logger, s store.SettingValues, opt options.AuthOpt, name, providerUrl string, force, validate, enable bool) (eap *types.ExternalAuthProvider, err error) {
	if !force {
		prefix := "auth.external.providers." + eap.Key + "."

		vv, _, err := store.SearchSettingValues(ctx, s, types.SettingsFilter{
			Prefix: prefix,
		})

		if err != nil || len(vv) > 0 {
			return nil, err
		}
	}

	if validate {
		// Do basic validation of external auth settings
		// will fail if secret or url are not set
		if err = staticValidateExternal(opt); err != nil {
			return
		}

		// Do full rediredct-URL check
		if err = validateExternalRedirectURL(opt); err != nil {
			return
		}
	}

	if opt.ExternalRedirectURL == "" {
		return nil, fmt.Errorf("refusing to register OIDC provider without redirect url")
	}

	p, err := parseExternalProviderUrl(providerUrl)
	if err != nil {
		return
	}

	eap, err = DiscoverOidcProvider(ctx, log, opt, name, p.String())
	if err != nil {
		return
	}

	vv, err := eap.EncodeKV()
	if err != nil {
		return
	}

	if enable {
		err = vv.Walk(func(value *types.SettingValue) error {
			if strings.HasSuffix(value.Name, ".enabled") {
				return value.SetSetting(true)
			}

			return nil
		})

		if err != nil {
			return
		}

		v := &types.SettingValue{Name: "auth.external.enabled"}
		err = v.SetSetting(true)
		if err != nil {
			return
		}
		vv = append(vv, v)
	}

	err = store.UpsertSettingValue(ctx, s, vv...)
	if err != nil {
		return
	}

	return
}

func parseExternalProviderUrl(in string) (p *url.URL, err error) {
	if i := strings.Index(in, "://"); i == -1 {
		// Add schema if missing
		in = "https://" + in
	}

	if p, err = url.Parse(in); err != nil {
		// Try to parse it
		return
	} else if i := strings.Index(p.Path, WellKnown); i > -1 {
		// Cut off well-known-path
		p.Path = p.Path[:i]
	}

	return
}

// StaticValidateExternal
//
// Simple checks of external auth settings
func staticValidateExternal(opt options.AuthOpt) error {
	if opt.ExternalRedirectURL == "" {
		return fmt.Errorf("redirect URL is empty")
	}

	const (
		tpt = "test-provider-test"
	)
	p, err := url.Parse(strings.Replace(opt.ExternalRedirectURL, "{provider}", tpt, 1))
	if err != nil {
		return fmt.Errorf("invalid redirect URL: %w", err)
	}

	if !strings.Contains(p.Path, tpt+"/callback") {
		return fmt.Errorf("could find injected provider in the URL, make sure you use '%%s' as a placeholder: %w", err)
	}

	if opt.ExternalCookieSecret == "" {
		return fmt.Errorf("AUTH_EXTERNAL_COOKIE_SECRET is empty")
	}

	if opt.SessionCookieSecure && p.Scheme != "https" {
		return fmt.Errorf("session store is secure, redirect URL should have HTTPS")
	}

	return nil
}

// ValidateExternalRedirectURL
//
// Validates external redirect URL
func validateExternalRedirectURL(opt options.AuthOpt) error {
	const tpt = "test-provider-test"
	const cb = "/callback"

	// Replace placeholders & remove /callback
	var url = strings.Replace(opt.ExternalRedirectURL, "{provider}", tpt, 1)
	url = url[0 : len(url)-len(cb)]

	rsp, err := http.DefaultClient.Get(url)
	if err != nil {
		return fmt.Errorf("could not get response from redirect URL: %w", err)
	}

	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	if strings.Contains(string(body), tpt) {
		return nil
	}

	return fmt.Errorf("could not validate external auth redirection URL")
}
