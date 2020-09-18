package importer

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v2"
	"io"
)

// Import performs standard import procedure with default services
func Import(ctx context.Context, ff ...io.Reader) (err error) {
	var (
		cc  types.ChannelSet
		aux interface{}
	)

	cc, _, err = service.DefaultChannel.With(ctx).Find(types.ChannelFilter{})
	if err != nil {
		return err
	}

	var (
		p   = rbac.NewImporter(service.DefaultAccessControl.Whitelist())
		imp = NewImporter(
			p,
			settings.NewImporter(),
			NewChannelImport(p, cc),
		)

		// At the moment, we can not load roles from system service
		// so we'll just use static set of known roles
		//
		// Roles are use for resolving access control
		roles = sysTypes.RoleSet{
			&sysTypes.Role{ID: rbac.EveryoneRoleID, Handle: "everyone"},
			&sysTypes.Role{ID: rbac.AdminsRoleID, Handle: "admins"},
		}
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
	// Store all imported
	return imp.Store(
		ctx,
		service.DefaultChannel.With(ctx),
		service.DefaultAccessControl,
		sysService.DefaultSettings,
		roles,
	)
}
