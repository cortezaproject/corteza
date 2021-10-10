package envoy

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestStoreYaml_pageRefs(t *testing.T) {
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
			name: "nested pages",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")

				parent := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "parent",
					Title:       "parent",
				}
				child := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "child",
					Title:       "child",
					SelfID:      parent.ID,
				}

				err := store.CreateComposePage(ctx, s, parent)
				if err != nil {
					t.Fatal(err)
				}
				err = store.CreateComposePage(ctx, s, child)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)

				parent, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "parent")
				req.NoError(err)

				child, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "child")
				req.NoError(err)

				req.Equal(parent.ID, child.SelfID)
			},
		},

		{
			name: "pageblock chart",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")
				chr := sTestComposeChart(ctx, t, s, ns.ID, mod.ID, "base")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "chart_1",
							Kind:  "Chart",
							Options: map[string]interface{}{
								"chart": strconv.FormatUint(chr.ID, 10),
							},
						},
						{
							Title: "chart_2",
							Kind:  "Chart",
							Options: map[string]interface{}{
								"chartID": strconv.FormatUint(chr.ID, 10),
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposeChart(&types.ChartFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				chr, err := store.LookupComposeChartByNamespaceIDHandle(ctx, s, n.ID, "base_chart")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)
				req.Len(pg.Blocks, 2)

				// provided as chart
				b := pg.Blocks[0]
				req.Equal(strconv.FormatUint(chr.ID, 10), b.Options["chartID"])

				// provided as chartID
				b = pg.Blocks[1]
				req.Equal(strconv.FormatUint(chr.ID, 10), b.Options["chartID"])
			},
		},

		{
			name: "pageblock comment",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "comment_1",
							Kind:  "Comment",
							Options: map[string]interface{}{
								"module": strconv.FormatUint(mod.ID, 10),
							},
						},
						{
							Title: "comment_2",
							Kind:  "Comment",
							Options: map[string]interface{}{
								"moduleID": strconv.FormatUint(mod.ID, 10),
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)
				req.Len(pg.Blocks, 2)

				// provided as module
				b := pg.Blocks[0]
				req.Equal(strconv.FormatUint(mod.ID, 10), b.Options["moduleID"])

				// provided as moduleID
				b = pg.Blocks[1]
				req.Equal(strconv.FormatUint(mod.ID, 10), b.Options["moduleID"])
			},
		},

		{
			name: "pageblock; automation",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				wf := sTestAutomationWorkflow(ctx, t, s, "base")
				sTestComposePageWithBlocks(ctx, t, s, ns.ID, "base",
					types.PageBlocks{
						{
							Title: "automation",
							Kind:  "Automation",
							Options: map[string]interface{}{
								"buttons": []map[string]interface{}{
									{
										"enabled":      true,
										"label":        "test",
										"resourceType": "compose",
										"variant":      "danger",
										"script":       "server-script/script1.js:default",
									},
									{
										"enabled":      true,
										"label":        "test",
										"resourceType": "compose",
										"stepID":       3,
										"variant":      "danger",
										"workflowID":   strconv.FormatUint(wf.ID, 10),
									},
								},
							},
						},
					},
				)

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					AutomationWorkflows(&atypes.WorkflowFilter{
						WorkflowID: []uint64{wf.ID},
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
						Handle:      "base_page",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				wf, err := store.LookupAutomationWorkflowByHandle(ctx, s, "base_handle")
				req.NoError(err)
				req.NotNil(wf)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "base_page")
				req.NoError(err)
				req.NotNil(pg)
				block := pg.Blocks[0]

				bb, _ := block.Options["buttons"].([]interface{})
				req.Len(bb, 2)

				b := (bb[0]).(map[string]interface{})
				req.Equal("server-script/script1.js:default", b["script"])

				b = (bb[1]).(map[string]interface{})
				req.Equal(strconv.FormatUint(wf.ID, 10), b["workflowID"])
				req.Equal(float64(3), b["stepID"])
			},
		},

		{
			name: "pageblock record",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					ModuleID:    mod.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "record_1",
							Kind:  "Record",
							Options: map[string]interface{}{
								"fields": []string{"f1", "f2", "f3"},
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)

				req.Equal(mod.ID, pg.ModuleID)

				req.Len(pg.Blocks, 1)
				b := pg.Blocks[0]
				req.Len(b.Options["fields"], 3)
			},
		},

		{
			name: "pageblock record list",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "rl_1",
							Kind:  "RecordList",
							Options: map[string]interface{}{
								"module": strconv.FormatUint(mod.ID, 10),
								"fields": []map[string]interface{}{
									{"name": "f1"},
									{"name": "f2"},
								},
							},
						},
						{
							Title: "rl_2",
							Kind:  "RecordList",
							Options: map[string]interface{}{
								"moduleID": strconv.FormatUint(mod.ID, 10),
								"fields": []map[string]interface{}{
									{"name": "f1"},
									{"name": "f2"},
								},
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)

				// provided as module
				b := pg.Blocks[0]
				req.Equal(strconv.FormatUint(mod.ID, 10), b.Options["moduleID"])
				casted := b.Options["fields"].([]interface{})
				for i, c := range casted {
					cc := c.(map[string]interface{})
					req.Equal(fmt.Sprintf("f%d", i+1), cc["name"].(string))
				}

				// provided as moduleID
				b = pg.Blocks[1]
				req.Equal(strconv.FormatUint(mod.ID, 10), b.Options["moduleID"])
				casted = b.Options["fields"].([]interface{})
				for i, c := range casted {
					cc := c.(map[string]interface{})
					req.Equal(fmt.Sprintf("f%d", i+1), cc["name"].(string))
				}
			},
		},

		{
			name: "pageblock record organizer",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "rl_1",
							Kind:  "RecordOrganizer",
							Options: map[string]interface{}{
								"module": strconv.FormatUint(mod.ID, 10),
							},
						},
						{
							Title: "rl_2",
							Kind:  "RecordOrganizer",
							Options: map[string]interface{}{
								"moduleID": strconv.FormatUint(mod.ID, 10),
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)

				// provided as module
				b := pg.Blocks[0]
				req.Equal(strconv.FormatUint(mod.ID, 10), b.Options["moduleID"])

				// provided as moduleID
				b = pg.Blocks[1]
				req.Equal(strconv.FormatUint(mod.ID, 10), b.Options["moduleID"])
			},
		},

		{
			name: "pageblock calendar",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod1 := sTestComposeModule(ctx, t, s, ns.ID, "base")
				mod2 := sTestComposeModule(ctx, t, s, ns.ID, "base_x")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "calendar_1",
							Kind:  "Calendar",
							Options: map[string]interface{}{
								"feeds": []map[string]interface{}{
									{
										"resource": "record",
										"options": map[string]interface{}{
											"module": strconv.FormatUint(mod1.ID, 10),
										},
									},
									{
										"resource": "record",
										"options": map[string]interface{}{
											"moduleID": strconv.FormatUint(mod2.ID, 10),
										},
									},
								},
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod1, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)
				mod2, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_x_module")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)
				b := pg.Blocks[0]

				ff, _ := b.Options["feeds"].([]interface{})
				req.Len(ff, 2)

				f := ff[0]
				feed, _ := f.(map[string]interface{})
				fOpts, _ := (feed["options"]).(map[string]interface{})
				req.Equal(strconv.FormatUint(mod1.ID, 10), fOpts["moduleID"])

				f = ff[1]
				feed, _ = f.(map[string]interface{})
				fOpts, _ = (feed["options"]).(map[string]interface{})
				req.Equal(strconv.FormatUint(mod2.ID, 10), fOpts["moduleID"])
			},
		},

		{
			name: "pageblock metric",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod1 := sTestComposeModule(ctx, t, s, ns.ID, "base")
				mod2 := sTestComposeModule(ctx, t, s, ns.ID, "base_x")

				pg := &types.Page{
					ID:          su.NextID(),
					NamespaceID: ns.ID,
					Handle:      "page",
					Title:       "page",
					Blocks: types.PageBlocks{
						{
							Title: "metric_1",
							Kind:  "Metric",
							Options: map[string]interface{}{
								"metrics": []map[string]interface{}{
									{
										"module": strconv.FormatUint(mod1.ID, 10),
									},
									{
										"moduleID": strconv.FormatUint(mod2.ID, 10),
									},
								},
							},
						},
					},
				}

				err := store.CreateComposePage(ctx, s, pg)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&types.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&types.ModuleFilter{
						NamespaceID: ns.ID,
					}).
					ComposePage(&types.PageFilter{
						NamespaceID: ns.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod1, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)
				mod2, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_x_module")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "page")
				req.NoError(err)
				b := pg.Blocks[0]

				mm, _ := b.Options["metrics"].([]interface{})
				req.Len(mm, 2)

				m := mm[0]
				mops, _ := m.(map[string]interface{})
				req.Equal(strconv.FormatUint(mod1.ID, 10), mops["moduleID"])

				m = mm[1]
				mops, _ = m.(map[string]interface{})
				req.Equal(strconv.FormatUint(mod2.ID, 10), mops["moduleID"])
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
			se := su.NewStoreEncoder(s, &su.EncoderConfig{})
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
