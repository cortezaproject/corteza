package envoy

import (
	"context"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestStoreYaml_moduleFieldRefs(t *testing.T) {
	type (
		tc struct {
			name string
			// Before the data gets processed
			pre func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter)
			// After the data gets processed
			postStoreDecode func(req *require.Assertions, err error)
			postYamlEncode  func(req *require.Assertions, err error)
			postStoreEncode func(req *require.Assertions, err error)
			// Data assertions
			check func(ctx context.Context, s store.Storer, req *require.Assertions)
		}
	)

	ctx := context.Background()
	s := initServices(ctx, t)
	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	cases := []*tc{
		{
			name: "user field; role filter",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				rl := sTestRole(ctx, t, s, "base")

				modID := su.NextID()
				mod := &types.Module{
					ID:          modID,
					Name:        "usr_rl",
					Handle:      "usr_rl",
					NamespaceID: ns.ID,
					Fields: types.ModuleFieldSet{
						&types.ModuleField{
							ID:       su.NextID(),
							ModuleID: modID,
							Kind:     "User",
							Place:    0,
							Name:     "usr_rl",
							Options: types.ModuleFieldOptions{
								"roles": []string{strconv.FormatUint(rl.ID, 10)},
							},
						},
					},
				}

				err := store.CreateComposeModule(ctx, s, mod)
				if err != nil {
					t.Fatal(err)
				}
				err = store.CreateComposeModuleField(ctx, s, mod.Fields...)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					Roles(&systemTypes.RoleFilter{}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "usr_rl")
				req.NoError(err)
				role, err := store.LookupRoleByHandle(ctx, s, "base_role")
				req.NoError(err)
				req.NotNil(role)

				mff, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{
					ModuleID: []uint64{mod.ID},
				})
				req.NoError(err)

				// Check module relations for Options.module variant
				f := mff.FindByName("usr_rl")
				req.NotNil(f)

				rr := f.Options["roles"].([]interface{})
				req.Len(rr, 1)
				req.Equal(strconv.FormatUint(role.ID, 10), rr[0])
			},
		},

		{
			name: "external module ref",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod1 := sTestComposeModule(ctx, t, s, ns.ID, "base")

				mod2ID := su.NextID()
				mod2 := &types.Module{
					ID:          mod2ID,
					Name:        "r_rel",
					Handle:      "r_rel",
					NamespaceID: ns.ID,
					Fields: types.ModuleFieldSet{
						&types.ModuleField{
							ID:       su.NextID(),
							ModuleID: mod2ID,
							Kind:     "Record",
							Place:    0,
							Name:     "r_rel_f1",
							Options: types.ModuleFieldOptions{
								"module": strconv.FormatUint(mod1.ID, 10),
							},
						},
						&types.ModuleField{
							ID:       su.NextID(),
							ModuleID: mod2ID,
							Kind:     "Record",
							Place:    0,
							Name:     "r_rel_f2",
							Options: types.ModuleFieldOptions{
								"moduleID": strconv.FormatUint(mod1.ID, 10),
							},
						},
					},
				}

				err := store.CreateComposeModule(ctx, s, mod2)
				if err != nil {
					t.Fatal(err)
				}
				err = store.CreateComposeModuleField(ctx, s, mod2.Fields...)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod1, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)

				mod2, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "r_rel")
				req.NoError(err)

				mff, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{
					ModuleID: []uint64{mod2.ID},
				})
				req.NoError(err)

				// Check module relations for Options.module variant
				f := mff.FindByName("r_rel_f1")
				req.NotNil(f)
				req.Equal(int64(mod1.ID), f.Options.Int64("moduleID"))

				// Check module relations for Options.moduleID variant
				f = mff.FindByName("r_rel_f2")
				req.NotNil(f)
				req.Equal(int64(mod1.ID), f.Options.Int64("moduleID"))
			},
		},

		{
			name: "self module ref",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod2ID := su.NextID()
				mod2 := &types.Module{
					ID:          mod2ID,
					Name:        "r_rel",
					Handle:      "r_rel",
					NamespaceID: ns.ID,
					Fields: types.ModuleFieldSet{
						&types.ModuleField{
							ID:       su.NextID(),
							ModuleID: mod2ID,
							Kind:     "Record",
							Place:    0,
							Name:     "r_rel_f1",
							Options: types.ModuleFieldOptions{
								"module": strconv.FormatUint(mod2ID, 10),
							},
						},
						&types.ModuleField{
							ID:       su.NextID(),
							ModuleID: mod2ID,
							Kind:     "Record",
							Place:    0,
							Name:     "r_rel_f2",
							Options: types.ModuleFieldOptions{
								"moduleID": strconv.FormatUint(mod2ID, 10),
							},
						},
					},
				}

				err := store.CreateComposeModule(ctx, s, mod2)
				if err != nil {
					t.Fatal(err)
				}
				err = store.CreateComposeModuleField(ctx, s, mod2.Fields...)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)

				mod2, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "r_rel")
				req.NoError(err)

				mff, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{
					ModuleID: []uint64{mod2.ID},
				})
				req.NoError(err)

				// Check module relations for Options.module variant
				f := mff.FindByName("r_rel_f1")
				req.NotNil(f)
				req.Equal(int64(mod2.ID), f.Options.Int64("moduleID"))

				// Check module relations for Options.moduleID variant
				f = mff.FindByName("r_rel_f2")
				req.NotNil(f)
				req.Equal(int64(mod2.ID), f.Options.Int64("moduleID"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := require.New(t)

			truncateStore(ctx, s, t)
			err, df := c.pre(ctx, s)
			if err != nil {
				t.Fatal(err.Error())
			}
			// Decode from store
			sd := su.Decoder()
			nn, err := sd.Decode(ctx, s, df)
			if c.postStoreDecode != nil {
				c.postStoreDecode(req, err)
			} else {
				req.NoError(err)
			}

			// Encode into YAML
			ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{})
			bld := envoy.NewBuilder(ye)
			g, err := bld.Build(ctx, nn...)
			req.NoError(err)
			err = envoy.Encode(ctx, g, ye)
			ss := ye.Stream()
			if c.postYamlEncode != nil {
				c.postYamlEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Cleanup the store
			truncateStore(ctx, s, t)

			// Encode back into store
			se := su.NewStoreEncoder(s, dal.Service(), &su.EncoderConfig{})
			yd := yaml.Decoder()
			nn = make([]resource.Interface, 0, len(nn))
			for _, s := range ss {
				mm, err := yd.Decode(ctx, s.Source, nil)
				req.NoError(err)
				nn = append(nn, mm...)
			}
			bld = envoy.NewBuilder(se)
			g, err = bld.Build(ctx, nn...)
			req.NoError(err)

			err = envoy.Encode(ctx, g, se)
			if c.postStoreEncode != nil {
				c.postStoreEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Assert
			c.check(ctx, s, req)

			// Cleanup the store
			truncateStore(ctx, s, t)
		})
		ni = 0
		truncateStore(ctx, s, t)
	}
}
