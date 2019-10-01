package importer

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

// Import performs standard import procedure with default services
func Import(ctx context.Context, ns *types.Namespace, ff ...io.Reader) (err error) {
	var (
		aux interface{}
		imp = NewImporter(
			service.DefaultNamespace.With(ctx),
			service.DefaultModule.With(ctx),
			service.DefaultChart.With(ctx),
			service.DefaultPage.With(ctx),
			service.DefaultInternalAutomationManager,
			permissions.NewImporter(service.DefaultAccessControl.Whitelist()),
		)

		// At the moment, we can not load roles from system service
		// so we'll just use static set of known roles
		//
		// Roles are use for resolving access control
		roles = sysTypes.RoleSet{
			&sysTypes.Role{ID: permissions.EveryoneRoleID, Handle: "everyone"},
			&sysTypes.Role{ID: permissions.AdminsRoleID, Handle: "admins"},
		}
	)

	for _, f := range ff {
		if err = yaml.NewDecoder(f).Decode(&aux); err != nil {
			return
		}

		if ns != nil {
			// If we're importing with --namespace switch,
			// we're going to import all into one NS

			err = imp.GetNamespaceImporter().Cast(ns.Slug, aux)
		} else {
			// importing one or more namespaces
			err = imp.Cast(aux)
		}

		if err != nil {
			return
		}
	}

	// Get roles across the system
	// roles, err := service.DefaultSystemRole.Find(ctx)
	// if err != nil {
	// 	return
	// }

	// Store all imported
	return imp.Store(
		ctx,
		service.DefaultNamespace.With(ctx),
		service.DefaultModule.With(ctx),
		service.DefaultChart.With(ctx),
		service.DefaultPage.With(ctx),
		service.DefaultRecord.With(ctx),
		service.DefaultInternalAutomationManager,
		service.DefaultAccessControl,
		roles,
	)
}
