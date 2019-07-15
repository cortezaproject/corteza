package external

import (
	"context"
	"fmt"

	"github.com/crusttech/go-oidc"

	"github.com/cortezaproject/corteza-server/system/internal/service"
)

// @todo remove dependency on github.com/crusttech/go-oidc (and github.com/coreos/go-oidc)
//       and move client registration to corteza codebase
func RegisterNewOpenIdClient(ctx context.Context, eas service.AuthSettings, name, url string) (eap *service.AuthSettingsExternalAuthProvider, err error) {
	var (
		provider    *oidc.Provider
		client      *oidc.Client
		redirectUrl = fmt.Sprintf(eas.ExternalRedirectUrl, "openid-connect."+name)
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
		return
	}

	eap = &service.AuthSettingsExternalAuthProvider{
		RedirectUrl: redirectUrl,
		Key:         client.ID,
		Secret:      client.Secret,
		IssuerUrl:   url,
	}

	return
}
