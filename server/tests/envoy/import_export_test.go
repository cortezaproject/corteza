package envoy

import (
	"context"
	"fmt"
	"os"
	"testing"

	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestImportExport(t *testing.T) {
	var (
		ctx   = context.Background()
		req   = require.New(t)
		nodes envoyx.NodeSet
		gg    *envoyx.DepGraph
		err   error
	)

	cleanup(t)

	// The test
	//
	// * imports some YAML files
	// * checks the DB state
	// * exports the DB into a YAML
	// * clears the DB
	// * imports the YAML
	// * checks the DB state the same way as before
	//
	// The above outlined flow allows us to trivially check if the data is both
	// imported and exported correctly.
	//
	// The initial step could also manually populate the DB but the YAML import
	// is more convenient.

	t.Run("initial import", func(t *testing.T) {
		t.Run("parse configs", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": "file://testdata/full",
				},
			})
			req.NoError(err)
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, nil, nodes...)
			req.NoError(err)
		})

		t.Run("import into DB", func(t *testing.T) {
			err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, gg)
			req.NoError(err)
		})

		assertFullState(ctx, t, defaultStore, req)
	})

	// Prepare a temp file where we'll dump the YAML into
	auxFile, err := os.CreateTemp(os.TempDir(), "*.yaml")
	req.NoError(err)
	spew.Dump(auxFile.Name())
	// defer os.Remove(auxFile.Name())
	defer auxFile.Close()

	t.Run("export", func(t *testing.T) {
		t.Run("export from DB", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
				Filter: map[string]envoyx.ResourceFilter{
					types.ChartResourceType:      {},
					types.ModuleResourceType:     {},
					types.NamespaceResourceType:  {},
					types.PageResourceType:       {},
					types.PageLayoutResourceType: {},

					systemTypes.ApplicationResourceType:         {},
					systemTypes.ApigwRouteResourceType:          {},
					systemTypes.AuthClientResourceType:          {},
					systemTypes.QueueResourceType:               {},
					systemTypes.ReportResourceType:              {},
					systemTypes.RoleResourceType:                {},
					systemTypes.TemplateResourceType:            {},
					systemTypes.UserResourceType:                {},
					systemTypes.DalConnectionResourceType:       {},
					systemTypes.DalSensitivityLevelResourceType: {},

					automationTypes.WorkflowResourceType: {},
				},
			})
			req.NoError(err)
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, nil, nodes...)
			req.NoError(err)
		})

		t.Run("write file", func(t *testing.T) {
			err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeIo,
				Params: map[string]any{
					"writer": auxFile,
				},
			}, gg)
			req.NoError(err)
		})
	})

	cleanup(t)

	t.Run("second import", func(t *testing.T) {
		t.Run("yaml parse", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": fmt.Sprintf("file://%s", auxFile.Name()),
				},
			})
			req.NoError(err)
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, nil, nodes...)
			req.NoError(err)
		})

		t.Run("run import", func(t *testing.T) {
			err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, gg)
			req.NoError(err)
		})

		assertFullState(ctx, t, defaultStore, req)
	})
}

func assertFullState(ctx context.Context, t *testing.T, s store.Storer, req *require.Assertions) {
	t.Run("check state", func(t *testing.T) {
		t.Run("corteza::compose", func(t *testing.T) {
			// Namespaces
			ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "test_ns_1")
			req.NoError(err)
			req.NotNil(ns)

			// Modules
			modules, _, err := store.SearchComposeModules(ctx, s, types.ModuleFilter{
				NamespaceID: ns.ID,
				Sorting: filter.Sorting{
					Sort: filter.SortExprSet{{Column: "handle", Descending: false}},
				},
			})
			req.NoError(err)

			mod1 := modules.FindByHandle("test_ns_1_mod_1")
			req.NotNil(mod1)

			mod2 := modules.FindByHandle("test_ns_1_mod_2")
			req.NotNil(mod2)

			// Fields
			fields, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{
				ModuleID: []uint64{mod1.ID, mod2.ID},
			})
			req.NoError(err)

			assignModuleFields(modules, fields)

			req.Len(mod1.Fields, 2)
			req.Len(mod2.Fields, 2)

			req.Equal(mod1.ID, mod2.Fields.FindByName("test_ns_1_mod_2_f1").Options.UInt64("moduleID"))
			req.Equal(mod2.ID, mod2.Fields.FindByName("test_ns_1_mod_2_f2").Options.UInt64("moduleID"))

			// Charts
			charts, _, err := store.SearchComposeCharts(ctx, s, types.ChartFilter{
				NamespaceID: ns.ID,
			})
			req.NoError(err)
			ch1 := charts.FindByHandle("test_ns_1_chart_1")
			req.NotNil(ch1)

			req.Equal(mod1.ID, ch1.Config.Reports[0].ModuleID)
			req.Equal(mod2.ID, ch1.Config.Reports[1].ModuleID)

			// Pages
			pages, _, err := store.SearchComposePages(ctx, s, types.PageFilter{
				NamespaceID: ns.ID,
			})
			req.NoError(err)

			pg1 := pages.FindByHandle("test_ns_1_page_1")
			req.NotNil(pg1)

			req.Len(pg1.Blocks, 8)
			// @todo test page block references
			// spew.Dump(pg1.Blocks)

			rpg1 := pages.FindByHandle("test_ns_1_record_page_1")
			req.NotNil(rpg1)
			req.Equal(mod1.ID, rpg1.ModuleID)

			// Page layouts
			layouts, _, err := store.SearchComposePageLayouts(ctx, s, types.PageLayoutFilter{
				NamespaceID: ns.ID,
			})
			req.NoError(err)

			ly1 := layouts.FindByHandle("test_ns_1_page_1_layout_1")
			req.NotNil(ly1)

			ly2 := layouts.FindByHandle("test_ns_1_page_1_layout_2")
			req.NotNil(ly2)

			// Record page layouts
			ly3 := layouts.FindByHandle("test_ns_1_record_page_1_layout_1")
			req.NotNil(ly3)

			ly4 := layouts.FindByHandle("test_ns_1_record_page_1_layout_2")
			req.NotNil(ly4)
		})
		t.Run("corteza::system", func(t *testing.T) {
			// Users
			users, _, err := store.SearchUsers(ctx, s, systemTypes.UserFilter{})
			req.NoError(err)
			req.Len(users, 1)

			u1 := users[0]
			req.Equal("test_user_1", u1.Handle)

			// ApigwRoutes
			apigwroutes, _, err := store.SearchApigwRoutes(ctx, s, systemTypes.ApigwRouteFilter{})
			req.NoError(err)
			req.Len(apigwroutes, 1)

			r1 := apigwroutes[0]
			req.Equal("/test/endpoint/1", r1.Endpoint)
			req.Equal(u1.ID, r1.CreatedBy)
			req.Equal(u1.ID, r1.UpdatedBy)

			// ApigwFilters
			apigwfilter, _, err := store.SearchApigwFilters(ctx, s, systemTypes.ApigwFilterFilter{})
			req.NoError(err)
			req.Len(apigwfilter, 1)

			f1 := apigwfilter[0]
			req.Equal(r1.ID, f1.Route)

			// AuthClients
			clients, _, err := store.SearchAuthClients(ctx, s, systemTypes.AuthClientFilter{})
			req.NoError(err)
			req.Len(clients, 1)

			c1 := clients[0]
			req.Equal("test_authclient_1", c1.Handle)
			req.Equal(u1.ID, c1.OwnedBy)
			req.Equal(u1.ID, c1.CreatedBy)
			req.Equal(u1.ID, c1.UpdatedBy)

			// Queues
			queues, _, err := store.SearchQueues(ctx, s, systemTypes.QueueFilter{})
			req.NoError(err)
			req.Len(queues, 1)

			q1 := queues[0]
			req.Equal("test_queue_1", q1.Queue)
			req.Equal(u1.ID, q1.CreatedBy)
			req.Equal(u1.ID, q1.UpdatedBy)

			// Reports
			reports, _, err := store.SearchReports(ctx, s, systemTypes.ReportFilter{})
			req.NoError(err)
			req.Len(reports, 1)

			rp1 := reports[0]
			req.Equal("test_report_1", rp1.Handle)
			req.Equal(u1.ID, rp1.OwnedBy)
			req.Equal(u1.ID, rp1.CreatedBy)
			req.Equal(u1.ID, rp1.UpdatedBy)

			// Roles
			roles, _, err := store.SearchRoles(ctx, s, systemTypes.RoleFilter{})
			req.NoError(err)

			rl1 := roles[0]
			req.Equal("test_role_1", rl1.Handle)

			// Templates
			templates, _, err := store.SearchTemplates(ctx, s, systemTypes.TemplateFilter{})
			req.NoError(err)
			req.Len(templates, 1)

			tpl1 := templates[0]
			req.Equal("test_template_1", tpl1.Handle)
			req.Equal(u1.ID, tpl1.OwnerID)

			// DalConnections
			connections, _, err := store.SearchDalConnections(ctx, s, systemTypes.DalConnectionFilter{})
			req.NoError(err)
			req.Len(connections, 1)
			conn1 := connections[0]
			req.Equal("test_connection_1", conn1.Handle)

			req.Equal(u1.ID, conn1.CreatedBy)
			req.Equal(u1.ID, conn1.UpdatedBy)

			// DalSensitivityLevels
			dalsensitivitylevels, _, err := store.SearchDalSensitivityLevels(ctx, s, systemTypes.DalSensitivityLevelFilter{})
			req.NoError(err)
			req.Len(dalsensitivitylevels, 1)
			sens1 := dalsensitivitylevels[0]

			req.Equal("test_sensitivity_level_1", sens1.Handle)
			req.Equal(u1.ID, sens1.CreatedBy)
			req.Equal(u1.ID, sens1.UpdatedBy)
		})
		t.Run("corteza::automation", func(t *testing.T) {
			// Workflows
			workflows, _, err := store.SearchAutomationWorkflows(ctx, s, automationTypes.WorkflowFilter{})
			req.NoError(err)
			req.Len(workflows, 1)

			wf1 := workflows[0]
			req.Equal("test_workflow_1", wf1.Handle)

			// Triggers
			triggers, _, err := store.SearchAutomationTriggers(ctx, s, automationTypes.TriggerFilter{})
			req.NoError(err)
			req.Len(triggers, 1)

			t1 := triggers[0]
			req.Equal(wf1.ID, t1.WorkflowID)
		})
	})
}

func assignModuleFields(modules types.ModuleSet, fields types.ModuleFieldSet) {
	for _, m := range modules {
		for _, f := range fields {
			if m.ID == f.ModuleID {
				m.Fields = append(m.Fields, f)
			}
		}
	}
}
