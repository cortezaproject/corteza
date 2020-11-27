package envoy

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestModuleRels(t *testing.T) {
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

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite, file string) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := yd(ctx, suite, file)
		req.NoError(err)

		return req, encode(ctx, s, nn)
	}

	cases := []*tc{
		{
			name:  "simple mods; self ref",
			suite: "mod_rel",
			file:  "modules_self_ref",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),

					storeNamespace(ctx, s, 10, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 10, "mod1")
				req.NoError(err)
				req.NotNil(mod)

				mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
				req.NoError(err)

				ff := mod.Fields
				req.Len(ff, 1)
				req.Equal("Record", ff[0].Kind)
				req.Equal(mod.ID, uint64(ff[0].Options.Int64("module")))
			},
		},

		{
			name:  "simple mods; ref to peer",
			suite: "mod_rel",
			file:  "modules_peer_ref",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeModuleFields(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				pmod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 100, "mod2")
				req.NoError(err)
				req.NotNil(pmod)

				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 100, "mod1")
				req.NoError(err)
				req.NotNil(mod)

				mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
				req.NoError(err)

				ff := mod.Fields
				req.Len(ff, 1)
				req.Equal("Record", ff[0].Kind)
				req.Equal(pmod.ID, uint64(ff[0].Options.Int64("module")))
			},
		},

		{
			name:  "simple mods; ref to store",
			suite: "mod_rel",
			file:  "modules_store_ref",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeModuleFields(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
					storeModule(ctx, s, 100, 200, "mod2"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				smod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 100, "mod2")
				req.NoError(err)
				req.NotNil(smod)

				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 100, "mod1")
				req.NoError(err)
				req.NotNil(mod)

				mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
				req.NoError(err)

				ff := mod.Fields
				req.Len(ff, 1)
				req.Equal("Record", ff[0].Kind)
				req.Equal(smod.ID, uint64(ff[0].Options.Int64("module")))
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s; %s/%s", c.name, c.suite, c.file), func(t *testing.T) {
			err = c.pre()
			if err != nil {
				t.Fatal(err.Error())
			}
			req, err := prepare(ctx, s, t, c.suite, c.file+".yaml")
			c.post(req, err)
			c.check(req)
		})

		ni = 0
	}
}
