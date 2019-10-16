package importer

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/settings"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	Importer struct {
		channels    *Channel
		permissions importer.PermissionImporter
		settings    importer.SettingImporter
	}

	channelFinder interface {
		Find(context.Context) (types.ChannelSet, error)
	}
)

func NewImporter(p importer.PermissionImporter, s importer.SettingImporter, ci *Channel) *Importer {
	return &Importer{
		channels:    ci,
		permissions: p,
		settings:    s,
	}
}

func (imp *Importer) Cast(in interface{}) (err error) {
	return deinterfacer.Each(in, func(index int, key string, val interface{}) (err error) {
		switch key {
		case "channels":
			return imp.channels.CastSet(val)
		case "channel":
			return imp.channels.CastSet([]interface{}{val})

		case "settings":
			return imp.settings.CastSet(val)

		case "allow", "deny":
			return imp.permissions.CastResourcesSet(key, val)

		default:
			err = fmt.Errorf("unexpected key %q", key)
		}

		return err
	})
}

func (imp *Importer) Store(ctx context.Context, rk channelKeeper, pk permissions.ImportKeeper, sk settings.ImportKeeper, roles sysTypes.RoleSet) (err error) {
	err = imp.channels.Store(ctx, rk)
	if err != nil {
		return
	}

	// Make sure we properly replace channel handles with IDs
	roles.Walk(func(r *sysTypes.Role) error {
		imp.permissions.UpdateRoles(r.Handle, r.ID)
		return nil
	})

	err = imp.permissions.Store(ctx, pk)
	if err != nil {
		return
	}

	err = imp.settings.Store(ctx, sk)
	if err != nil {
		return
	}

	return nil
}
