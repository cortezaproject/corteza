package envoy

import (
	"context"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestProvision(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)

		rrls = rbac.RuleSet{}
		ns   = &types.Namespace{}
		m1   = &types.Module{}
		m2   = &types.Module{}
		m3   = &types.Module{}
		ch1  = &types.Chart{}
		pg1  = &types.Page{}
		pg2  = &types.Page{}
	)
	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	ni := uint64(0)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite string) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := dd(ctx, suite)
		req.NoError(err)

		return req, encode(ctx, s, nn)
	}

	fullModLoad := func(ctx context.Context, s store.Storer, req *require.Assertions, nsID uint64, handle string) (*types.Module, error) {
		mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, nsID, handle)
		req.NoError(err)
		req.NotNil(mod)

		mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
		req.NoError(err)
		req.NotNil(mod.Fields)
		return mod, err
	}

	getOptSlice := func(req *require.Assertions, key string, f *types.ModuleField) []interface{} {
		opt, _ := (f.Options[key]).([]interface{})
		req.NotNil(opt)
		return opt
	}

	// Prepare
	s, err = initStore(ctx)
	err = ce(
		err,

		s.TruncateActionlogs(ctx),
		s.TruncateApplications(ctx),
		s.TruncateAttachments(ctx),
		s.TruncateComposeAttachments(ctx),
		s.TruncateComposeCharts(ctx),
		s.TruncateComposeNamespaces(ctx),
		s.TruncateComposeModules(ctx),
		s.TruncateComposeModuleFields(ctx),
		s.TruncateComposePages(ctx),

		storeRole(ctx, s, 1, "everyone"),
		storeRole(ctx, s, 2, "admins"),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := prepare(ctx, s, t, "provision")
	req.NoError(err)

	t.Run("ebac rules", func(t *testing.T) {
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
		m1, err = fullModLoad(ctx, s, req, ns.ID, "mod1")
		m2, err = fullModLoad(ctx, s, req, ns.ID, "mod2")
		m3, err = fullModLoad(ctx, s, req, ns.ID, "mod3")

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
	})

	t.Run("setting records", func(t *testing.T) {
		ms, err := fullModLoad(ctx, s, req, ns.ID, "settings")
		rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
		req.NoError(err)
		req.Len(rr, 1)

		sr := rr[0]
		req.Equal("s1 value", sr.Values.Get("s1", 0).Value)
		req.Equal("1", sr.Values.Get("s2", 0).Value)
	})

}
