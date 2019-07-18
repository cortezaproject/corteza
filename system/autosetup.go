package system

import (
	"context"
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/system/internal/auth/external"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func accessControlSetup(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	c.InitServices(ctx, c)

	// Calling grant directly on internal permissions service to avoid AC check for "grant"
	var p = service.DefaultPermissions
	var ac = service.DefaultAccessControl
	return p.Grant(ctx, ac.Whitelist(), ac.DefaultRules()...)
}

func makeDefaultApplications(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	db, err := factory.Database.Get(system)
	if err != nil {
		return err
	}

	repo := repository.Application(ctx, db)

	aa, err := repo.Find()
	if err != nil {
		return err
	}

	// List of apps to create.
	//
	// We use Unify.Url field for matching,
	// so make sure it's always present!
	defApps := types.ApplicationSet{
		&types.Application{
			Name:    "CRM",
			Enabled: true,
			Unify: &types.ApplicationUnify{
				Name:   "CRM",
				Listed: true,
				Icon:   "/applications/crust_favicon.png",
				Logo:   "/applications/crust.jpg",
				Url:    "/compose/ns/crm/pages",
			},
		},
	}

	return defApps.Walk(func(defApp *types.Application) error {
		for _, a := range aa {
			if a.Unify != nil && a.Unify.Url == defApp.Unify.Url {
				// App already added.
				return nil
			}
		}

		defApp, err = repo.Create(defApp)
		c.Log.Info(
			"creating default application",
			zap.String("name", defApp.Name),
			zap.Uint64("name", defApp.ID),
			zap.Error(err),
		)
		return err

		return nil
	})
}

func discoverSettings(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	return service.DefaultSettings.With(ctx).AutoDiscovery()
}

func oidcAutoDiscovery(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
	var provider = strings.TrimSpace(options.EnvString("", "PROVISION_OIDC_PROVIDER", ""))

	c.Log.Debug("OIDC auto discovery provision", zap.String("providers", provider))

	if len(provider) == 0 {
		return
	}

	var (
		providers  = strings.Split(provider, " ")
		plen       = len(providers)
		name, purl string
		eap        *service.AuthSettingsExternalAuthProvider
	)

	if plen%2 == 1 {
		return errors.New("expecting even number of providers")
	}

	for p := 0; p < plen; p = p + 2 {
		name, purl = providers[p], providers[p+1]

		// force:    false
		// because we do not want to override the provider every time the system restarts
		//
		// validate: false
		// because at the initial (empty db) start, we can not actually validate (server is not yet up)
		//
		// enable:   true
		// we want provider & the entire external auth to be validated
		eap, err = external.RegisterOidcProvider(ctx, name, purl, false, false, true)

		if err != nil {
			c.Log.Error(
				"could not register oidc provider",
				zap.String("url", purl),
				zap.String("name", name),
				zap.Error(err))
			return
		} else if eap == nil {
			c.Log.Info("provider already exists",
				zap.String("name", name))
		} else {
			c.Log.Info("provider successfuly registered",
				zap.String("url", purl),
				zap.String("key", eap.Key),
				zap.String("name", name))
		}
	}

	return
}

func authAddExternals(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
	var (
		kinds = []string{
			"github",
			"facebook",
			"gplus",
			"linkedin",
			"oidc",
		}

		env, p, name string

		pp []string

		eap service.AuthSettingsExternalAuthProvider
	)

	for _, kind := range kinds {
		env = "PROVISION_SETTINGS_AUTH_EXTERNAL_" + strings.ToUpper(kind)

		p = strings.TrimSpace(options.EnvString("", env, ""))
		if len(p) == 0 {
			continue
		}

		eap = service.AuthSettingsExternalAuthProvider{Enabled: true}

		if kind == "oidc" {
			pp = strings.SplitN(p, " ", 4)

			// Spread name, issuer-url, key and secret from provision string for OIDC provider
			name, eap.IssuerUrl, eap.Key, eap.Secret = pp[0], pp[1], pp[2], pp[3]

			name = external.OIDC_PROVIDER_PREFIX + name
		} else {
			pp = strings.SplitN(p, " ", 2)

			// Spread key and secret from provision string
			eap.Key, eap.Secret = pp[0], pp[1]
			name = kind
		}

		_ = external.AddProvider(name, &eap, false)
	}

	return
}
