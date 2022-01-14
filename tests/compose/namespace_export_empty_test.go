package compose

import (
	"testing"
)

func Test_namespace_export_empty(t *testing.T) {
	ctx, h, s := setup(t)
	loadScenario(ctx, defStore, t, h)
	grantImportExport(h)

	ns, _, _, _, _, err := fetchEntireNamespace(ctx, s, "ns1")
	h.a.NoError(err)

	arch := namespaceExportSafe(t, h, ns.ID)
	sessionID := namespaceImportInitSafe(t, h, arch)
	ns, mm, pp, cc := namespaceImportRun(ctx, s, t, h, sessionID, "imported", "imported")

	h.a.Equal("imported", ns.Slug)
	h.a.NotEqual(0, ns.ID)

	h.a.Len(mm, 0)
	h.a.Len(pp, 0)
	h.a.Len(cc, 0)

	cleanup(t)
}
