package envoy

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestYamlStore_provision(t *testing.T) {
	var (
		ctx = context.Background()
		s   = initStore(ctx, t)
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	prep := func(ctx context.Context, t *testing.T, s store.Storer) {
		truncateStore(ctx, s, t)
		ni = uint64(10)

		err := collect(
			storeRole(ctx, s, 1, "everyone"),
			storeRole(ctx, s, 2, "admins"),
		)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	encode := func(ctx context.Context, s store.Storer, req *require.Assertions, dir string) error {
		nn, err := decodeDirectory(ctx, dir)
		req.NoError(err)

		return encode(ctx, s, nn)
	}

	getOptSlice := func(req *require.Assertions, key string, f *types.ModuleField) []interface{} {
		opt, _ := (f.Options[key]).([]interface{})
		req.NotNil(opt)
		return opt
	}

	t.Run("testdata/simple", func(t *testing.T) {
		var (
			req = require.New(t)
			err error

			rrls = rbac.RuleSet{}
			ns   = &types.Namespace{}
			m1   = &types.Module{}
			m2   = &types.Module{}
			m3   = &types.Module{}
			wf   = &atypes.Workflow{}
			tt   = atypes.TriggerSet{}
			ch1  = &types.Chart{}
			pg1  = &types.Page{}
			pg2  = &types.Page{}
			pgA  = &types.Page{}
		)

		prep(ctx, t, s)
		encode(ctx, s, req, "provision/simple")

		t.Run("rbac rules", func(t *testing.T) {
			rrls, _, err = store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
			req.NoError(err)
			req.NotNil(rrls)
			req.Len(rrls, 14)
		})

		t.Run("namespaces", func(t *testing.T) {
			ns, err = store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
			req.NoError(err)
			req.NotNil(ns)
		})

		t.Run("modules", func(t *testing.T) {
			m1, err = loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			m2, err = loadComposeModuleFull(ctx, s, req, ns.ID, "mod2")
			m3, err = loadComposeModuleFull(ctx, s, req, ns.ID, "mod3")

			// mod1
			req.Len(m1.Fields, 4)
			req.Equal("String", m1.Fields[0].Kind)

			req.Equal("Select", m1.Fields[1].Kind)
			opt := getOptSlice(req, "options", m1.Fields[1])
			req.Equal([]interface{}{"f2 opt 1", "f2 opt 2", "f2 opt 3"}, opt)

			req.Equal("Select", m1.Fields[2].Kind)
			opt = getOptSlice(req, "options", m1.Fields[2])
			req.Equal([]interface{}{"☆☆☆☆☆", "★☆☆☆☆"}, opt)

			req.Equal("Record", m1.Fields[3].Kind)
			req.Equal(strconv.FormatUint(m2.ID, 10), m1.Fields[3].Options.String("moduleID"))
			opt = getOptSlice(req, "queryFields", m1.Fields[3])
			req.Equal([]interface{}{"f1"}, opt)

			// mod2
			req.Len(m1.Fields, 4)

			req.Equal("Record", m2.Fields[1].Kind)
			req.Equal(strconv.FormatUint(m2.ID, 10), m2.Fields[1].Options.String("moduleID"))
			opt = getOptSlice(req, "queryFields", m2.Fields[1])
			req.Equal([]interface{}{"f1"}, opt)

			req.Equal("Record", m2.Fields[3].Kind)
			req.Equal(strconv.FormatUint(m3.ID, 10), m2.Fields[3].Options.String("moduleID"))
			opt = getOptSlice(req, "queryFields", m2.Fields[3])
			req.Equal([]interface{}{"f1"}, opt)
		})

		t.Run("workflows", func(t *testing.T) {
			wf, err = store.LookupAutomationWorkflowByHandle(ctx, s, "test")
			req.NoError(err)
			req.NotNil(wf)
			tt, _, err = store.SearchAutomationTriggers(ctx, s, atypes.TriggerFilter{
				WorkflowID: []uint64{wf.ID},
			})
			req.NoError(err)
			req.NotNil(tt)
			req.Len(tt, 1)

			req.Len(wf.Steps, 1)
			s := wf.Steps[0]
			req.Equal(uint64(3), s.ID)

			tr := tt[0]
			req.Equal(uint64(3), tr.StepID)
		})

		t.Run("charts", func(t *testing.T) {
			ch1, err = store.LookupComposeChartByNamespaceIDHandle(ctx, s, ns.ID, "chr1")
			req.NoError(err)
			req.NotNil(ch1)

			req.Len(ch1.Config.Reports, 1)
			r := ch1.Config.Reports[0]
			req.Len(r.Dimensions, 1)
			req.Len(r.Metrics, 1)
			req.Equal(m1.ID, r.ModuleID)
		})

		t.Run("pages", func(t *testing.T) {
			pg1, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "pg1")
			req.NoError(err)
			req.NotNil(pg1)
			pg2, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "rpg2")
			req.NoError(err)
			req.NotNil(pg2)
			pgA, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "pgA")
			req.NoError(err)
			req.NotNil(pgA)

			// pg1
			req.Len(pg1.Blocks, 3)

			req.Equal("Content", pg1.Blocks[0].Kind)
			req.Equal("pg1 b1 content body", pg1.Blocks[0].Options["body"])

			req.Equal("RecordList", pg1.Blocks[1].Kind)
			req.Equal(strconv.FormatUint(m1.ID, 10), (pg1.Blocks[1].Options["moduleID"].(string)))

			req.Equal("Chart", pg1.Blocks[2].Kind)
			req.Equal(strconv.FormatUint(ch1.ID, 10), (pg1.Blocks[2].Options["chartID"].(string)))

			// pg2
			req.Equal(m1.ID, pg2.ModuleID)
			req.Equal(pg1.ID, pg2.SelfID)

			req.Equal("Record", pg2.Blocks[0].Kind)

			// pgA
			b := pgA.Blocks[0].Options["buttons"].([]interface{})[0]
			opts := b.(map[string]interface{})

			req.Equal(strconv.FormatUint(wf.ID, 10), opts["workflowID"])
			req.Equal(float64(3), opts["stepID"])
		})

		t.Run("setting records", func(t *testing.T) {
			ms, err := loadComposeModuleFull(ctx, s, req, ns.ID, "settings")
			rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 1)

			sr := rr[0]
			req.Equal("s1 value", sr.Values.Get("s1", 0).Value)
			req.Equal("1", sr.Values.Get("s2", 0).Value)
		})

		t.Run("templates", func(t *testing.T) {
			tpls, _, err := store.SearchTemplates(ctx, s, stypes.TemplateFilter{})
			req.NoError(err)
			req.Len(tpls, 4)
		})

		truncateStore(ctx, s, t)
	})

	// This one checks the absolute worst condition, where we are importing 2 namespaces
	// where all of the underlying resources (modules, pages) are the exact same.
	//
	// Concrete case; CRM and SS where they both have Contacts and Accounts -- this caused
	// some conflicts in the logic.
	t.Run("testdata/multiple_apps", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		check := func(ctx context.Context, slug string) {
			var (
				ns   = &types.Namespace{}
				mod1 = &types.Module{}
				mod2 = &types.Module{}
				mods = &types.Module{}
				pg1  = &types.Page{}
				rpg2 = &types.Page{}
				chr1 = &types.Chart{}
				rr   = types.RecordSet{}
			)

			// Preload things
			// * NS
			ns, err := store.LookupComposeNamespaceBySlug(ctx, s, slug)
			req.NoError(err)
			req.NotNil(ns)

			// * Mods
			mod1, err = loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NoError(err)
			req.NotNil(mod1)

			mod2, err = loadComposeModuleFull(ctx, s, req, ns.ID, "mod2")
			req.NoError(err)
			req.NotNil(mod2)

			mods, err = loadComposeModuleFull(ctx, s, req, ns.ID, "settings")
			req.NoError(err)
			req.NotNil(mods)

			// * Pages
			pg1, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "pg1")
			req.NoError(err)
			req.NotNil(pg1)

			rpg2, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "rpg2")
			req.NoError(err)
			req.NotNil(rpg2)

			// * Charts
			chr1, err = store.LookupComposeChartByNamespaceIDHandle(ctx, s, ns.ID, "chr1")
			req.NoError(err)
			req.NotNil(chr1)

			// * Records
			rr, _, err = store.SearchComposeRecords(ctx, s, mods, types.RecordFilter{
				ModuleID:    mods.ID,
				NamespaceID: ns.ID,
			})
			req.NoError(err)
			req.NotNil(rr)
			req.Len(rr, 1)

			// Check things
			t.Run("NS", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s name", slug), ns.Name)
			})

			t.Run("mod1", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s mod1 name", slug), mod1.Name)
				req.Len(mod1.Fields, 1)

				f1 := mod1.Fields[0]
				req.Equal(strconv.FormatUint(mod2.ID, 10), f1.Options.String("moduleID"))
			})

			t.Run("mod2", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s mod2 name", slug), mod2.Name)
				req.Len(mod2.Fields, 3)

				f1 := mod2.Fields[1]
				req.Equal(strconv.FormatUint(mod2.ID, 10), f1.Options.String("moduleID"))
			})

			t.Run("mods", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s settings", slug), mods.Name)
				req.Len(mods.Fields, 1)
			})

			t.Run("pg1", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s pg1 title", slug), pg1.Title)
				req.Len(pg1.Blocks, 4)

				req.Equal("pg1 RecordList", pg1.Blocks[0].Title)
				req.Equal(strconv.FormatUint(mod1.ID, 10), pg1.Blocks[0].Options["moduleID"])

				req.Equal("pg1 Chart", pg1.Blocks[1].Title)
				req.Equal(strconv.FormatUint(chr1.ID, 10), pg1.Blocks[1].Options["chartID"])

				req.Equal("pg1 Calendar", pg1.Blocks[2].Title)
				f := pg1.Blocks[2].Options["feeds"].([]interface{})[0]
				feed, _ := f.(map[string]interface{})
				fOpts, _ := (feed["options"]).(map[string]interface{})
				id, _ := fOpts["moduleID"].(string)
				req.Equal(strconv.FormatUint(mod1.ID, 10), id)

				req.Equal("pg1 Metric", pg1.Blocks[3].Title)
				m := pg1.Blocks[3].Options["metrics"].([]interface{})[0]
				mops, _ := m.(map[string]interface{})
				id = mops["moduleID"].(string)
				req.Equal(strconv.FormatUint(mod1.ID, 10), id)
			})

			t.Run("rpg2", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s Record page for module \"mod1\"", slug), rpg2.Title)
			})

			t.Run("chr1", func(t *testing.T) {
				req.Equal(fmt.Sprintf("%s chr1 name", slug), chr1.Name)
				req.Equal(mod1.ID, chr1.Config.Reports[0].ModuleID)
			})

			t.Run("rr", func(t *testing.T) {
				req.Len(rr, 1)
				v := rr[0].Values.FilterByName("s1")[0]
				req.Equal(fmt.Sprintf("%s value", slug), v.Value)
			})
		}

		prep(ctx, t, s)
		encode(ctx, s, req, "provision/multiple_apps")

		check(ctx, "ns1")
		check(ctx, "ns2")

		truncateStore(ctx, s, t)
	})
}
