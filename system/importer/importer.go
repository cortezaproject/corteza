package importer

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

type (
	Importer struct {
		roles       *Role
		permissions importer.PermissionImporter
		settings    importer.SettingImporter
	}

	roleFinder interface {
		Find(context.Context) (types.RoleSet, error)
	}
)

func NewImporter(p importer.PermissionImporter, s importer.SettingImporter, ri *Role) *Importer {
	return &Importer{
		roles:       ri,
		permissions: p,
		settings:    s,
	}
}

func (imp *Importer) Cast(in interface{}) (err error) {
	return deinterfacer.Each(in, func(index int, key string, val interface{}) (err error) {
		switch key {
		case "roles":
			return imp.roles.CastSet(val)
		case "role":
			return imp.roles.CastSet([]interface{}{val})

		case "allow", "deny":
			return imp.permissions.CastResourcesSet(key, val)

		case "settings":
			return imp.settings.CastSet(val)

		default:
			err = fmt.Errorf("unexpected key %q", key)
		}

		return err
	})
}

func (imp *Importer) Store(
	ctx context.Context,
	rk roleKeeper,
	pk permissions.ImportKeeper,
	sk settings.ImportKeeper,
	roles types.RoleSet,
) (err error) {
	err = imp.roles.Store(ctx, rk)
	if err != nil {
		return
	}

	// Make sure we properly replace role handles with IDs
	roles.Walk(func(role *types.Role) error {
		imp.permissions.UpdateRoles(role.Handle, role.ID)
		return nil
	})

	err = imp.permissions.Store(ctx, pk)
	if err != nil {
		return
	}

	err = imp.settings.Store(ctx, sk)
	if err != nil {
		return errors.Wrap(err, "could not import settings")
	}

	return nil
}
