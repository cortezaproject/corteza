package system

import (
	"context"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/importer"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io"
)

type (
	settingsService interface {
		FindByPrefix(context.Context, ...string) (types.SettingValueSet, error)
		BulkSet(context.Context, types.SettingValueSet) error
	}
)

func Provision(ctx context.Context, log *zap.Logger) (err error) {
	var provisioned bool
	var readers []io.Reader

	if provisioned, err = notProvisioned(ctx); err != nil {
		return err
	} else if !provisioned {
		log.Info("provisioning system")
		readers, err = impAux.ReadStatic(Asset)
		if err != nil {
			return err
		}

		if err = importer.Import(ctx, readers...); err != nil {
			return err
		}

	} else {
		log.Info("provisioning system")
		// When already provisioned, make sure settings are re-provisioned
		readers, err = impAux.ReadStatic(Asset)
		if err != nil {
			return err
		}

		if err = partialImportSettings(ctx, service.DefaultSettings, readers...); err != nil {
			return err
		}
	}

	if err = makeDefaultApplications(ctx, log); err != nil {
		return
	}
	if err = authSettingsAutoDiscovery(ctx, log, service.DefaultSettings); err != nil {
		return
	}
	if err = authAddExternals(ctx, log); err != nil {
		return
	}
	if err = service.DefaultSettings.UpdateCurrent(ctx); err != nil {
		return
	}
	if err = oidcAutoDiscovery(ctx, log); err != nil {
		return
	}

	return nil
}

// provision ONLY when there are no rules for role admins
func notProvisioned(ctx context.Context) (bool, error) {
	return len(service.DefaultPermissions.FindRulesByRoleID(permissions.AdminsRoleID)) == 0, nil
}

func makeDefaultApplications(ctx context.Context, log *zap.Logger) error {
	db := factory.Database.MustGet()

	repo := repository.Application(ctx, db)

	aa, _, err := repo.Find(types.ApplicationFilter{})
	if err != nil {
		return err
	}

	// Update icon & logo on all apps
	const (
		oldIconUrl = "/applications/crust_favicon.png"
		oldLogoUrl = "/applications/crust.jpg"

		newIconUrl = "/applications/default_icon.png"
		newLogoUrl = "/applications/default_logo.jpg"
	)

	for a := 0; a < len(aa); a++ {
		var dirty bool

		if aa[a].Unify == nil {
			continue
		}

		if aa[a].Unify.Icon == oldIconUrl {
			aa[a].Unify.Icon = newIconUrl
			dirty = true
		}

		if aa[a].Unify.Logo == oldLogoUrl {
			aa[a].Unify.Logo = newLogoUrl
			dirty = true
		}

		if !dirty {
			continue
		}

		if aa[a], err = repo.Update(aa[a]); err != nil {
			return err
		}
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
				Icon:   newIconUrl,
				Logo:   newLogoUrl,
				Url:    "/compose/ns/crm/pages",
			},
		},

		&types.Application{
			Name:    "Service Cloud",
			Enabled: true,
			Unify: &types.ApplicationUnify{
				Name:   "Service Cloud",
				Listed: true,
				Icon:   newIconUrl,
				Logo:   newLogoUrl,
				Url:    "/compose/ns/service-cloud/pages",
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
		log.Info(
			"creating default application",
			zap.String("name", defApp.Name),
			zap.Uint64("name", defApp.ID),
			zap.Error(err),
		)
		return err
	})
}

// Partial import of settings from provision files
func partialImportSettings(ctx context.Context, ss settingsService, ff ...io.Reader) (err error) {
	var (
		// decoded content from YAML files
		aux interface{}

		si = settings.NewImporter()

		// importer w/o permissions & roles
		// we need only settings
		imp = importer.NewImporter(nil, si, nil)

		// current value
		current types.SettingValueSet

		// unexisting values
		unex types.SettingValueSet
	)

	for _, f := range ff {
		if err = yaml.NewDecoder(f).Decode(&aux); err != nil {
			return
		}

		err = imp.Cast(aux)
		if err != nil {
			return
		}
	}

	// Get all "current" settings storage
	current, err = ss.FindByPrefix(ctx)
	if err != nil {
		return
	}

	// Compare current settings with imported, get all that do not exist yet
	if unex = si.GetValues(); len(unex) > 0 {
		// Store non existing
		err = ss.BulkSet(ctx, current.New(unex))
		if err != nil {
			return
		}
	}

	return nil
}
