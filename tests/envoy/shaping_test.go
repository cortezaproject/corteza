package envoy

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestShaping_simple(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)

		suites = []string{"shaping_csv_simple", "shaping_jsonl_simple"}
	)
	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
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

	ni := uint64(0)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite string) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := dd(ctx, suite)
		req.NoError(err)

		crs := resource.ComposeRecordShaper()
		nn, err = resource.Shape(nn, crs)
		req.NoError(err)

		return req, encode(ctx, s, nn)
	}
	// Prepare
	s, err = initStore(ctx)
	err = ce(
		err,

		s.TruncateRbacRules(ctx),
		s.TruncateRoles(ctx),
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

	for _, st := range suites {
		req, err := prepare(ctx, s, t, st)
		req.NoError(err)

		t.Run(fmt.Sprintf("record shaping; %s", st), func(t *testing.T) {
			ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
			req.NotNil(ns)
			ms, err := fullModLoad(ctx, s, req, ns.ID, "mod1")
			req.NotNil(ms)

			rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 2)

			r1 := rr[0]
			r2 := rr[1]

			req.Len(r1.Values, 2)
			req.Equal("c1.v1", r1.Values.FilterByName("f1")[0].Value)
			req.Equal("c2.v1", r1.Values.FilterByName("f2")[0].Value)

			req.Len(r2.Values, 2)
			req.Equal("c1.v2", r2.Values.FilterByName("f1")[0].Value)
			req.Equal("c2.v2", r2.Values.FilterByName("f2")[0].Value)

			s.TruncateComposeRecords(ctx, ms)
		})
	}
}
