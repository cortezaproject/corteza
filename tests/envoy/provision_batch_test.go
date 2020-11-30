package envoy

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

// TestProvision_batch simulates the worst possible case, where two namespaces
// have exact same items.
func TestProvision_batch(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)
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
		s.TruncateComposeRecords(ctx, nil),

		storeRole(ctx, s, 1, "everyone"),
		storeRole(ctx, s, 2, "admins"),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := prepare(ctx, s, t, "provision_batch")
	req.NoError(err)

	checkBatchProvision(ctx, t, req, s, "ns1")
	checkBatchProvision(ctx, t, req, s, "ns2")
}

func checkBatchProvision(ctx context.Context, t *testing.T, req *require.Assertions, s store.Storer, slug string) {
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
	mod1, err = fullModLoad(ctx, s, req, ns.ID, "mod1")
	req.NoError(err)
	req.NotNil(mod1)

	mod2, err = fullModLoad(ctx, s, req, ns.ID, "mod2")
	req.NoError(err)
	req.NotNil(mod2)

	mods, err = fullModLoad(ctx, s, req, ns.ID, "settings")
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

func fullModLoad(ctx context.Context, s store.Storer, req *require.Assertions, nsID uint64, handle string) (*types.Module, error) {
	mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, nsID, handle)
	req.NoError(err)
	req.NotNil(mod)

	mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	req.NoError(err)
	req.NotNil(mod.Fields)
	return mod, err
}
