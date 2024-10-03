package compose

import (
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

func Test_namespace_export_simple(t *testing.T) {
	ctx, h, s := setup(t)
	loadScenario(ctx, defStore, t, h)
	grantImportExport(h)

	ns, _, _, _, _, err := fetchEntireNamespace(ctx, s, "ns1")
	h.a.NoError(err)

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "export")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.export")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "charts.export")

	arch := namespaceExportSafe(t, h, ns.ID)
	sessionID := namespaceImportInitSafe(t, h, arch)
	ns, mm, pp, cc := namespaceImportRun(ctx, s, t, h, sessionID, "imported", "imported")

	h.a.Equal("imported", ns.Slug)
	h.a.NotEqual(0, ns.ID)

	h.a.Len(mm, 3)
	exp := map[string]bool{"mod1": true, "mod2": true, "mod3": true}

	for i, m := range mm {
		h.a.True(exp[m.Handle])
		if i > 0 {
			h.a.NotEqual(m.Handle, mm[i-1].Handle)
		}
	}

	h.a.Len(pp, 3)
	exp = map[string]bool{
		"pg1":  true,
		"rpg2": true,
		"pg2":  true,
	}
	for i, p := range pp {
		h.a.True(exp[p.Handle])
		if i > 0 {
			h.a.NotEqual(p.Handle, mm[i-1].Handle)
		}
	}

	parent := pp.FindByHandle("rpg2")
	child := pp.FindByHandle("pg2")
	h.a.Equal(child.SelfID, parent.ID)

	h.a.Len(cc, 1)
	h.a.Equal("chr1", cc[0].Handle)

	cleanup(t)
}
