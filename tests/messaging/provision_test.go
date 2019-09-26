package messaging

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/messaging/importer"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func TestProvisioning(t *testing.T) {
	h := newHelper(t)

	var (
		aux interface{}

		roles = sysTypes.RoleSet{
			&sysTypes.Role{ID: permissions.EveryoneRoleID, Handle: "everyone"},
			&sysTypes.Role{ID: permissions.AdminRoleID, Handle: "admins"},
		}

		ctx = h.secCtx()
		pi  = permissions.NewImporter(service.DefaultAccessControl.Whitelist())
		imp = importer.NewImporter(pi, importer.NewChannelImport(pi, nil))
	)

	h.allow(types.MessagingPermissionResource, "grant")

	var f, err = os.Open("../../provision/messaging/001_permission_rules.yaml")
	h.a.NoError(err)
	h.a.NoError(yaml.NewDecoder(f).Decode(&aux))

	h.a.NoError(imp.Cast(aux))

	h.a.NoError(imp.Store(
		ctx,
		service.DefaultChannel.With(ctx),
		service.DefaultAccessControl,
		roles,
	))
}
