package provision

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

var (
	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)

func Run(ctx context.Context, log *zap.Logger, s store.Storer, provisionOpt options.ProvisionOpt, authOpt options.AuthOpt) error {
	log = log.Named("provision")

	ffn := []func() error{
		// Migrations:
		func() error { return migrateApplications(ctx, s) },
		func() error { return migrateEmailTemplates(ctx, log.Named("email-templates"), s) },
		func() error { return migratePre202109Roles(ctx, log.Named("pre-202109-roles"), s) },
		func() error { return migratePre202109RbacRules(ctx, log.Named("pre-202109-rbac-rules"), s) },
		func() error { return cleanupPre202109Settings(ctx, log.Named("pre-202109-settings"), s) },
		func() error { return migrateResourceTranslations(ctx, log.Named("resource-translations"), s) },
		func() error { return migrateReportIdentifiers(ctx, log.Named("report-identifiers"), s) },
		func() error {
			return migratePost202203ResourceTranslations(ctx, log.Named("post-202203-resource-translations"), s)
		},

		// Config (full & partial)
		func() error { return importConfig(ctx, log.Named("config"), s, provisionOpt.Path) },

		// Auto-discoveries and other parts that cannot be imported from static files
		func() error { return emailSettings(ctx, s) },
		func() error { return authAddExternals(ctx, log.Named("auth.externals"), s) },
		func() error { return oidcAutoDiscovery(ctx, log.Named("auth.oidc-auto-discovery"), s, authOpt) },
		func() error { return defaultAuthClient(ctx, log.Named("auth.clients"), s, authOpt) },
	}

	for _, fn := range ffn {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// defaultAuthClient checks if default client exists (handle = AUTH_DEFAULT_CLIENT) and adds it
func defaultAuthClient(ctx context.Context, log *zap.Logger, s store.AuthClients, authOpt options.AuthOpt) error {
	if authOpt.DefaultClient == "" {
		// Default client not set
		return nil
	}

	c := &types.AuthClient{
		ID:     id.Next(),
		Handle: authOpt.DefaultClient,
		Meta: &types.AuthClientMeta{
			Name: "Corteza Web Applications",
		},
		ValidGrant: "authorization_code",
		RedirectURI: func() string {
			// Disabling protection by redirection URL for now, it caused too much confusion on simple setups
			//baseURL, _ := url.Parse(authOpt.BaseURL)
			//return fmt.Sprintf("%s://%s", baseURL.Scheme, baseURL.Hostname())
			return ""
		}(),

		Secret:    string(rand.Bytes(64)),
		Scope:     "profile api",
		Enabled:   true,
		Trusted:   true,
		Security:  &types.AuthClientSecurity{},
		Labels:    nil,
		CreatedAt: *now(),
	}

	_, err := store.LookupAuthClientByHandle(ctx, s, c.Handle)
	if err == nil {
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	if err = store.CreateAuthClient(ctx, s, c); err != nil {
		return err
	}

	log.Info(
		"Added OAuth2 client",
		zap.String("name", c.Meta.Name),
		zap.String("redirectURI", c.RedirectURI),
		zap.Uint64("clientId", c.ID),
	)

	return nil
}
