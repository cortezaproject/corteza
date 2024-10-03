package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/store"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/spf13/cast"
)

func fetchEntireNamespace(ctx context.Context, s store.Storer, slug string) (ns *types.Namespace, mm types.ModuleSet, pp types.PageSet, cc types.ChartSet, slg string, err error) {
	slg = slug

	ns, err = store.LookupComposeNamespaceBySlug(ctx, s, slug)
	if err != nil {
		return
	}

	mm, _, err = store.SearchComposeModules(ctx, s, types.ModuleFilter{NamespaceID: ns.ID})
	if err != nil {
		return
	}

	for i := 0; i < len(mm); i++ {
		mm[i].Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mm[i].ID}})
		if err != nil {
			return
		}
	}

	pp, _, err = store.SearchComposePages(ctx, s, types.PageFilter{NamespaceID: ns.ID})
	if err != nil {
		return
	}

	cc, _, err = store.SearchComposeCharts(ctx, s, types.ChartFilter{NamespaceID: ns.ID})
	if err != nil {
		return
	}

	return
}

func findModuleByHandle(mm types.ModuleSet, h string) *types.Module {
	for _, m := range mm {
		if m.Handle == h {
			return m
		}
	}
	return nil
}

func findChartByHandle(cc types.ChartSet, h string) *types.Chart {
	for _, m := range cc {
		if m.Handle == h {
			return m
		}
	}
	return nil
}

func findPageByHandle(pp types.PageSet, h string) *types.Page {
	for _, m := range pp {
		if m.Handle == h {
			return m
		}
	}
	return nil
}

func Test_namespace_duplicate(t *testing.T) {
	ctx, h, s := setup(t)
	loadScenario(ctx, defStore, t, h)
	ns, _, _, _, _, err := fetchEntireNamespace(ctx, s, "ns1")
	h.a.NoError(err)

	helpers.AllowMe(h, types.ComponentRbacResource(), "namespace.create")
	helpers.AllowMe(h, types.ComponentRbacResource(), "module.create")
	helpers.AllowMe(h, types.ComponentRbacResource(), "page.create")
	helpers.AllowMe(h, types.ComponentRbacResource(), "chart.create")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "charts.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "pages.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "export")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.export")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "charts.export")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ChartRbacResource(0, 0), "read")
	helpers.AllowMe(h, sysTypes.ComponentRbacResource(), "roles.search")
	helpers.AllowMe(h, sysTypes.RoleRbacResource(0), "read")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/clone", ns.ID)).
		JSON(`{ "slug": "cloned", "name": "cloned name" }`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	usedIDs := make(map[uint64]bool)
	checkID := func(h helper, id uint64) {
		if usedIDs[id] {
			h.a.FailNow(fmt.Sprintf("the ID is not unique across cloned resources: %d", id))
		}
		usedIDs[id] = true
	}

	checker := func(ns *types.Namespace, mm types.ModuleSet, pp types.PageSet, cc types.ChartSet, slug string, err error) {
		h.a.NoError(err)

		// NS
		h.a.Equal(slug, ns.Slug)

		// Modules
		mod1 := findModuleByHandle(mm, "mod1")
		h.a.NotNil(mod1)
		checkID(h, mod1.ID)
		mod2 := findModuleByHandle(mm, "mod2")
		h.a.NotNil(mod2)
		checkID(h, mod2.ID)
		mod3 := findModuleByHandle(mm, "mod3")
		h.a.NotNil(mod3)
		checkID(h, mod3.ID)

		h.a.Len(mod1.Fields, 4)
		h.a.Equal("Record", mod1.Fields.FindByName("f4").Kind)
		h.a.Equal(mod2.ID, mod1.Fields.FindByName("f4").Options.UInt64("moduleID"))

		h.a.Len(mod2.Fields, 4)
		h.a.Equal("Record", mod2.Fields.FindByName("f_ref_self").Kind)
		h.a.Equal(mod2.ID, mod2.Fields.FindByName("f_ref_self").Options.UInt64("moduleID"))
		h.a.Equal("Record", mod2.Fields.FindByName("f2").Kind)
		h.a.Equal(mod3.ID, mod2.Fields.FindByName("f2").Options.UInt64("moduleID"))

		h.a.Len(mod3.Fields, 1)

		// Charts
		chr1 := findChartByHandle(cc, "chr1")
		checkID(h, chr1.ID)
		h.a.NotNil(chr1)
		h.a.Len(chr1.Config.Reports, 1)
		h.a.Equal(mod1.ID, chr1.Config.Reports[0].ModuleID)

		// Pages
		pg1 := findPageByHandle(pp, "pg1")
		checkID(h, pg1.ID)
		h.a.NotNil(pg1)
		rpg2 := findPageByHandle(pp, "rpg2")
		checkID(h, rpg2.ID)
		h.a.NotNil(rpg2)

		h.a.Len(pg1.Blocks, 3)

		h.a.Equal("RecordList", pg1.Blocks[1].Kind)
		h.a.Equal(mod1.ID, cast.ToUint64(pg1.Blocks[1].Options["moduleID"]))
		h.a.Equal("Chart", pg1.Blocks[2].Kind)
		h.a.Equal(chr1.ID, cast.ToUint64(pg1.Blocks[2].Options["chartID"]))

		h.a.Equal(rpg2.ModuleID, mod1.ID)
	}

	checker(fetchEntireNamespace(ctx, s, "ns1"))
	checker(fetchEntireNamespace(ctx, s, "cloned"))
}
