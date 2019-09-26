package system

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/system/importer"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func TestProvisioning(t *testing.T) {
	h := newHelper(t)

	var (
		aux interface{}

		roles = types.RoleSet{
			&types.Role{ID: permissions.EveryoneRoleID, Handle: "everyone"},
			&types.Role{ID: permissions.AdminRoleID, Handle: "admins"},
		}

		ctx = h.secCtx()
		pi  = permissions.NewImporter(service.DefaultAccessControl.Whitelist())
		imp = importer.NewImporter(pi, importer.NewRoleImport(pi, roles))
	)

	h.allow(types.SystemPermissionResource, "grant")
	h.allow(types.SystemPermissionResource, "role.create")
	h.allow(types.RolePermissionResource.AppendWildcard(), "update")

	var f, err = os.Open("../../provision/system/001_permission_rules.yaml")
	h.a.NoError(err)
	h.a.NoError(yaml.NewDecoder(f).Decode(&aux))

	h.a.NoError(imp.Cast(aux))

	h.a.NoError(imp.Store(
		ctx,
		service.DefaultRole.With(ctx),
		service.DefaultAccessControl,
		roles,
	))
}
