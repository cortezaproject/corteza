package external

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/crusttech/go-oidc"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/system/internal/service"
)

func AddProvider(name string, eap *service.AuthSettingsExternalAuthProvider, force bool) error {
	var (
		as  = service.DefaultAuthSettings
		log = log().With(
			zap.Bool("force", force),
			zap.String("name", name),
			zap.String("key", eap.Key),
		)
	)

	if eap.IssuerUrl != "" {
		log = log.With(zap.String("issuer-url", eap.IssuerUrl))
	}

	log.Info("adding external auth provider")

	if !force {
		if e, exists := as.ExternalProviders[name]; exists && e.Key == eap.Key && e.Secret == eap.Secret {
			return nil
		}
	}

	if vv, err := eap.MakeValueSet(name); err != nil {
		log.Error("could not prepare settings", zap.Error(err))
		return err
	} else if err = service.DefaultIntSettings.BulkSet(vv); err != nil {
		log.Error("could not store settings", zap.Error(err))
		return err
	}

	log.Info("external provider added")
	return nil
}

// @todo remove dependency on github.com/crusttech/go-oidc (and github.com/coreos/go-oidc)
//       and move client registration to corteza codebase
func DiscoverOidcProvider(ctx context.Context, eas service.AuthSettings, name, url string) (eap *service.AuthSettingsExternalAuthProvider, err error) {
	var (
		provider    *oidc.Provider
		client      *oidc.Client
		redirectUrl = fmt.Sprintf(eas.ExternalRedirectUrl, OIDC_PROVIDER_PREFIX+name)

		log = log().With(
			zap.String("redirect-url", redirectUrl),
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

	eap = &service.AuthSettingsExternalAuthProvider{
		RedirectUrl: redirectUrl,
		Key:         client.ID,
		Secret:      client.Secret,
		IssuerUrl:   url,
	}

	log.Info("oidc provider registered", zap.String("key", client.ID))

	return
}

func RegisterOidcProvider(ctx context.Context, name, providerUrl string, force, validate, enable bool) (eap *service.AuthSettingsExternalAuthProvider, err error) {
	var (
		as = service.DefaultAuthSettings
	)

	if !force {
		if _, exists := as.ExternalProviders[OIDC_PROVIDER_PREFIX+name]; exists {
			return
		}
	}

	if validate {
		// Do basic validation of external auth settings
		// will fail if secret or url are not set
		if err = as.StaticValidateExternal(); err != nil {
			return
		}

		// Do full rediredct-URL check
		if err = as.ValidateExternalRedirectURL(); err != nil {
			return
		}
	}

	if as.ExternalRedirectUrl == "" {
		return nil, errors.New("refusing to register oidc privider withoit redirect url")
	}

	p, err := parseExternalProviderUrl(providerUrl)
	if err != nil {
		return
	}

	eap, err = DiscoverOidcProvider(ctx, as, name, p.String())
	if err != nil {
		return
	}

	vv, err := eap.MakeValueSet(OIDC_PROVIDER_PREFIX + name)
	if err != nil {
		return
	}

	if enable {
		err = vv.Walk(func(value *settings.Value) error {
			if strings.HasSuffix(value.Name, ".enabled") {
				return value.SetValue(true)
			}

			return nil
		})

		if err != nil {
			return
		}

		v := &settings.Value{Name: "auth.external.enabled"}
		err = v.SetValue(true)
		if err != nil {
			return
		}
		vv = append(vv, v)
	}

	err = service.DefaultIntSettings.BulkSet(vv)
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
