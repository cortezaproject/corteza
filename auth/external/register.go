package external

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/crusttech/go-oidc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func AddProvider(ctx context.Context, eap *types.ExternalAuthProvider, force bool) error {
	var (
		s   = service.CurrentSettings
		log = log.With(
			zap.Bool("force", force),
			zap.String("handle", eap.Handle),
			zap.String("key", eap.Key),
		)
	)

	if eap.IssuerUrl != "" {
		log = log.With(zap.String("issuer-url", eap.IssuerUrl))
	}

	log.Info("adding external authentication provider")

	if !force {
		if ex := s.Auth.External.Providers.FindByHandle(eap.Handle); ex != nil && ex.Key == eap.Key && ex.Secret == eap.Secret {
			return nil
		}
	}

	if vv, err := eap.EncodeKV(); err != nil {
		log.Error("could not prepare settings", zap.Error(err))
		return err
	} else if err = service.DefaultSettings.BulkSet(ctx, vv); err != nil {
		log.Error("could not store settings", zap.Error(err))
		return err
	}

	log.Info("external authentication provider added")
	return nil
}

// @todo remove dependency on github.com/crusttech/go-oidc (and github.com/coreos/go-oidc)
//       and move client registration to corteza codebase
func DiscoverOidcProvider(ctx context.Context, opt options.AuthOpt, name, url string) (eap *types.ExternalAuthProvider, err error) {
	var (
		provider    *oidc.Provider
		client      *oidc.Client
		redirectUrl = strings.Replace(opt.ExternalRedirectURL, "{provider}", OIDC_PROVIDER_PREFIX+name, 1)

		log = log.With(
			zap.String("name", name),
			zap.String("url", url),
		)
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

func RegisterOidcProvider(ctx context.Context, opt options.AuthOpt, name, providerUrl string, force, validate, enable bool) (eap *types.ExternalAuthProvider, err error) {
	var (
		s = service.CurrentSettings
	)

	if !force {
		if s.Auth.External.Providers.FindByHandle(OIDC_PROVIDER_PREFIX+name) != nil {
			return
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
		return nil, errors.New("refusing to register OIDC provider without redirect url")
	}

	p, err := parseExternalProviderUrl(providerUrl)
	if err != nil {
		return
	}

	eap, err = DiscoverOidcProvider(ctx, opt, name, p.String())
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
				return value.SetValue(true)
			}

			return nil
		})

		if err != nil {
			return
		}

		v := &types.SettingValue{Name: "auth.external.enabled"}
		err = v.SetValue(true)
		if err != nil {
			return
		}
		vv = append(vv, v)
	}

	err = service.DefaultSettings.BulkSet(ctx, vv)
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
		return errors.New("redirect URL is empty")
	}

	const (
		tpt = "test-provider-test"
	)
	p, err := url.Parse(strings.Replace(opt.ExternalRedirectURL, "{provider}", tpt, 1))
	if err != nil {
		return errors.Wrap(err, "invalid redirect URL")
	}

	if !strings.Contains(p.Path, tpt+"/callback") {
		return errors.Wrap(err, "could find injected provider in the URL, make sure you use '%s' as a placeholder")
	}

	if opt.ExternalCookieSecret == "" {
		return errors.New("AUTH_EXTERNAL_COOKIE_SECRET is empty")
	}

	if opt.SessionCookieSecure && p.Scheme != "https" {
		return errors.New("session store is secure, redirect URL should have HTTPS")
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
		return errors.Wrap(err, "could not get response from redirect URL")
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	if strings.Contains(string(body), tpt) {
		return nil
	}

	return errors.New("could not validate external auth redirection URL")
}
