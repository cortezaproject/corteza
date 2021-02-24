package envoy

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestYamlStore_relations(t *testing.T) {
	type (
		tc struct {
			name string
			file string

			// Before the data gets processed
			pre func() error
			// After the data gets processed
			post func(req *require.Assertions, err error)
			// Data assertions
			check func(req *require.Assertions)
		}
	)

	var (
		ctx       = auth.SetSuperUserContext(context.Background())
		namespace = "relations"
		s         = initStore(ctx, t)
		err       error
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	cases := []*tc{
		{
			name: "module; self ref",
			file: "modules_self_ref",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 10, "ns1"),
				)
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
				req.Equal(strconv.FormatUint(mod.ID, 10), ff[0].Options.String("moduleID"))
			},
		},

		{
			name: "mods; ref to peer",
			file: "modules_peer_ref",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
				)
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
				req.Equal(strconv.FormatUint(pmod.ID, 10), ff[0].Options.String("moduleID"))
			},
		},

		{
			name: "mods; ref to store",
			file: "modules_store_ref",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
					storeComposeModule(ctx, s, 100, 200, "mod2"),
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
				req.Equal(strconv.FormatUint(smod.ID, 10), ff[0].Options.String("moduleID"))
			},
		},
	}

	for _, c := range cases {
		f := c.file + ".yaml"
		t.Run(fmt.Sprintf("%s; testdata/%s/%s", c.name, namespace, f), func(t *testing.T) {
			truncateStore(ctx, s, t)

			req := require.New(t)

			if c.pre != nil {
				err = c.pre()
				req.NoError(err)
			}

			nn, err := decodeYaml(ctx, namespace, f)
			req.NoError(err)

			err = encode(ctx, s, nn)
			if c.post != nil {
				c.post(req, err)
			} else {
				req.NoError(err)
			}

			if c.check != nil {
				c.check(req)
			}

			truncateStore(ctx, s, t)
		})
	}
}
