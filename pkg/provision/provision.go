package provision

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	url "net/url"
	"time"
)

func Run(ctx context.Context, log *zap.Logger, s store.Storer, provisionOpt options.ProvisionOpt, authOpt options.AuthOpt) error {
	ffn := []func() error{
		// Old relics
		func() error { return roles(ctx, s) },

		// Migrations:
		func() error { return apps(ctx, s) },
		func() error { return migrateEmailTemplates(ctx, log, s) },

		// Config (full & partial)
		func() error { return importConfig(ctx, log, s, provisionOpt.Path) },

		// Auto-discoveries and other parts that can not be imported from static files
		func() error { return authSettingsAutoDiscovery(ctx, log, service.DefaultSettings) },
		func() error { return authAddExternals(ctx, log) },
		func() error { return service.DefaultSettings.UpdateCurrent(ctx) },
		func() error { return oidcAutoDiscovery(ctx, log, authOpt) },
		func() error { return defaultAuthClient(ctx, log, s, authOpt) },
	}

	for _, fn := range ffn {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func defaultAuthClient(ctx context.Context, log *zap.Logger, s store.AuthClients, authOpt options.AuthOpt) error {
	clients := types.AuthClientSet{
		&types.AuthClient{
			ID:     id.Next(),
			Handle: "corteza-webapp",
			Meta: &types.AuthClientMeta{
				Name: "Corteza Web Applications",
			},
			ValidGrant: "authorization_code",
			RedirectURI: func() string {
				baseURL, _ := url.Parse(authOpt.BaseURL)
				return fmt.Sprintf("%s://%s", baseURL.Scheme, baseURL.Hostname())
			}(),
			Secret:    string(rand.Bytes(64)),
			Scope:     "profile api",
			Enabled:   true,
			Trusted:   true,
			Security:  &types.AuthClientSecurity{},
			Labels:    nil,
			CreatedAt: time.Now(),
		},
	}

	for _, c := range clients {
		_, err := store.LookupAuthClientByHandle(ctx, s, c.Handle)
		if err == nil {
			continue
		}

		if !errors.IsNotFound(err) {
			return err
		}

		if err = store.CreateAuthClient(ctx, s, c); err != nil {
			return err
		}

		log.Info("Added OAuth2 client", zap.String("name", c.Meta.Name), zap.Uint64("clientId", c.ID))
	}

	return nil
}
