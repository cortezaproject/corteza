package compose

import (
	"testing"
)

func Test_namespace_export_automation(t *testing.T) {
	ctx, h, s := setup(t)
	loadScenario(ctx, defStore, t, h)
	grantImportExport(h)

	ns, _, _, _, _, err := fetchEntireNamespace(ctx, s, "ns1")
	h.a.NoError(err)

	arch := namespaceExportSafe(t, h, ns.ID)
	sessionID := namespaceImportInitSafe(t, h, arch)
	ns, _, pp, _ := namespaceImportRun(ctx, s, t, h, sessionID, "imported", "imported")

	h.a.Equal("imported", ns.Slug)
	h.a.NotEqual(0, ns.ID)

	p := pp.FindByHandle("pg1")
	h.a.Len(p.Blocks, 2)
	for _, b := range p.Blocks {
		h.a.NotEqual("Automation", b.Kind)
	}

	cleanup(t)
}
