package tests

import (
	"context"
	"fmt"
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

	cases := []*tc{
		{
			name:  "simple provision set",
			suite: "provision",
			pre: func() (err error) {
				s, err = initStore(ctx)
				return ce(
					err,
					storeRole(ctx, s, 1, "everyone"),
					storeRole(ctx, s, 2, "admins"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				// Checkup AC
				rrls, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
				req.NoError(err)
				req.NotNil(rrls)
				req.Len(rrls, 14)
				// -----------

				// Checkup namespaces
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
				req.NoError(err)
				req.NotNil(ns)
				// -----------

				// Checkup modules
				m1, err := fullModLoad(ctx, s, req, ns.ID, "mod1")
				m2, err := fullModLoad(ctx, s, req, ns.ID, "mod2")
				m3, err := fullModLoad(ctx, s, req, ns.ID, "mod3")

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
				req.Equal(m2.ID, uint64(m1.Fields[3].Options.Int64("module")))
				opt = getOptSlice(req, "queryFields", m1.Fields[3])
				req.Equal([]interface{}{"f1"}, opt)

				// mod2
				req.Len(m1.Fields, 4)

				req.Equal("Record", m2.Fields[1].Kind)
				req.Equal(m2.ID, uint64(m2.Fields[1].Options.Int64("module")))
				opt = getOptSlice(req, "queryFields", m2.Fields[1])
				req.Equal([]interface{}{"f1"}, opt)

				req.Equal("Record", m2.Fields[3].Kind)
				req.Equal(m3.ID, uint64(m2.Fields[3].Options.Int64("module")))
				opt = getOptSlice(req, "queryFields", m2.Fields[3])
				req.Equal([]interface{}{"f1"}, opt)
				// -----------

				// Checkup charts
				ch1, err := store.LookupComposeChartByNamespaceIDHandle(ctx, s, ns.ID, "chr1")
				req.NoError(err)
				req.NotNil(ch1)

				req.Len(ch1.Config.Reports, 1)
				r := ch1.Config.Reports[0]
				req.Len(r.Dimensions, 1)
				req.Len(r.Metrics, 1)
				req.Equal(m1.ID, r.ModuleID)
				// -----------

				// Checkup pages
				pg1, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "pg1")
				req.NoError(err)
				req.NotNil(pg1)
				pg2, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, ns.ID, "rpg2")
				req.NoError(err)
				req.NotNil(pg2)

				// pg1
				req.Len(pg1.Blocks, 3)

				req.Equal("Content", pg1.Blocks[0].Kind)
				req.Equal("pg1 b1 content body", pg1.Blocks[0].Options["body"])

				req.Equal("RecordList", pg1.Blocks[1].Kind)
				req.Equal(m1.ID, uint64((pg1.Blocks[1].Options["module"]).(float64)))

				req.Equal("Chart", pg1.Blocks[2].Kind)
				req.Equal(ch1.ID, uint64((pg1.Blocks[2].Options["chart"]).(float64)))

				// pg2
				req.Equal(m1.ID, pg2.ModuleID)
				req.Equal(pg1.ID, pg2.SelfID)

				req.Equal("Record", pg2.Blocks[0].Kind)
				// -----------

				// Checkup settings records
				ms, err := fullModLoad(ctx, s, req, ns.ID, "settings")
				rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
				req.NoError(err)
				req.Len(rr, 1)

				sr := rr[0]
				req.Equal("s1 value", sr.Values.Get("s1", 0).Value)
				req.Equal("1", sr.Values.Get("s2", 0).Value)
				// -----------
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s; %s", c.name, c.suite), func(t *testing.T) {
			err = c.pre()
			if err != nil {
				t.Fatal(err.Error())
			}
			req, err := prepare(ctx, s, t, c.suite)
			c.post(req, err)
			c.check(req)
		})

		ni = 0
	}
}
