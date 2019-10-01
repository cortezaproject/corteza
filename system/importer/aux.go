package importer

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Import performs standard import procedure with default services
func Import(ctx context.Context, ff ...io.Reader) (err error) {
	var (
		roles types.RoleSet
		aux   interface{}
	)

	roles, err = service.DefaultRole.With(ctx).Find(&types.RoleFilter{})
	if err != nil {
		return err
	}

	pi := permissions.NewImporter(service.DefaultAccessControl.Whitelist())
	imp := NewImporter(
		pi,
		NewRoleImport(pi, roles),
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

	// Get roles across the system
	// roles, err := service.DefaultSystemRole.Find(ctx)
	// if err != nil {
	// 	return
	// }

	// Store all imported
	return imp.Store(
		ctx,
		service.DefaultRole.With(ctx),
		service.DefaultAccessControl,
		roles,
	)
}
