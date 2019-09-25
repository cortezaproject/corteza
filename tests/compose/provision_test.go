package compose

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func TestProvisioning(t *testing.T) {
	h := newHelper(t)

	var (
		aux interface{}

		ctx = h.secCtx()
		imp = importer.NewImporter(
			service.DefaultNamespace.With(ctx),
			service.DefaultModule.With(ctx),
			service.DefaultChart.With(ctx),
			service.DefaultPage.With(ctx),
			service.DefaultInternalAutomationManager,
			permissions.NewImporter(service.DefaultAccessControl.Whitelist()),
		)
	)

	h.allow(types.ComposePermissionResource, "grant")

	var f, err = os.Open("../../provision/compose/001_permission_rules.yaml")
	h.a.NoError(err)
	h.a.NoError(yaml.NewDecoder(f).Decode(&aux))

	h.a.NoError(imp.Cast(aux))

	h.a.NoError(imp.Store(
		ctx,
		service.DefaultNamespace.With(ctx),
		service.DefaultModule.With(ctx),
		service.DefaultChart.With(ctx),
		service.DefaultPage.With(ctx),
		service.DefaultInternalAutomationManager,
		service.DefaultAccessControl,
		sysTypes.RoleSet{
			&sysTypes.Role{ID: 1, Handle: "everyone"},
			&sysTypes.Role{ID: 2, Handle: "admins"},
		},
	))
}
