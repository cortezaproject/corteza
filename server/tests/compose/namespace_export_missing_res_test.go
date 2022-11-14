package compose

import (
	"testing"

	"github.com/cortezaproject/corteza-server/store"
)

func Test_namespace_export_missing_res(t *testing.T) {
	ctx, h, s := setup(t)
	loadScenario(ctx, defStore, t, h)
	grantImportExport(h)

	ns, mm, _, _, _, err := fetchEntireNamespace(ctx, s, "ns1")
	h.a.NoError(err)

	// Removeing one of the resources
	h.a.NoError(store.DeleteComposeModuleByID(ctx, s, mm.FindByHandle("mod2").ID))

	_, err = namespaceExport(t, h, ns.ID)
	h.a.Error(err)

	cleanup(t)
}
