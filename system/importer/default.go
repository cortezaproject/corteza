package importer

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"gopkg.in/yaml.v2"
	"io"
)

// Import performs standard import procedure with default services
func Import(ctx context.Context, ff ...io.Reader) (err error) {
	var (
		roles types.RoleSet
		aux   interface{}
	)

	roles, _, err = service.DefaultRole.With(ctx).Find(types.RoleFilter{})

	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	if err != nil {
		return err
	}

	pi := rbac.NewImporter(service.DefaultAccessControl.Whitelist())
	imp := NewImporter(
		pi,
		settings.NewImporter(),
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
		service.DefaultSettings,
		roles,
	)
}
