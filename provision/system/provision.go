package system

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/provision/util"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	settingsService interface {
		FindByPrefix(context.Context, ...string) (types.SettingValueSet, error)
		BulkSet(context.Context, types.SettingValueSet) error
	}
)

func provisionRoles(ctx context.Context, s store.Storer) error {
	if set, _, err := store.SearchRoles(ctx, s, types.RoleFilter{}); err != nil {
		return err
	} else if len(set) > 0 {
		return nil
	}

	rr := types.RoleSet{
		&types.Role{ID: rbac.AdminsRoleID, Name: "Administrators", Handle: "admins"},
		&types.Role{ID: rbac.EveryoneRoleID, Name: "Everyone", Handle: "everyone"},
	}

	err := rr.Walk(func(r *types.Role) error {
		r.CreatedAt = time.Now()
		return store.CreateRole(ctx, s, r)
	})

	return err
}

func Provision(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	log.Info("provisioning system")
	err = provisionRoles(ctx, s)
	if err != nil {
		return err
	}

	// Provision from YAML files
	// - access control
	// - settings
	// - applications
	if err = util.EncodeStatik(ctx, s, Asset, "/"); err != nil {
		return err
	}

	// These ones need some extra things, so we'll leave them there
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

	log.Info("provisioned system")
	return nil
}
