package envoy

import (
	"context"
	"strconv"
	"testing"
	"time"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	ctypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	stypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

// TestStoreYaml_base takes data from s1, encodes it into yaml files, decodes
// created yaml files, encodes into s2 and compares the data from s2.
//
// Yaml marshling can be flaky (because of map structs) so this is the "best & simplest" approach
func TestStoreYaml_base(t *testing.T) {
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
			name: "base automation workflow",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				wf := sTestAutomationWorkflow(ctx, t, s, "base")
				sTestAutomationTrigger(ctx, t, s, wf.ID, "base")

				df := su.NewDecodeFilter().AutomationWorkflows(&atypes.WorkflowFilter{})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				wf, err := store.LookupAutomationWorkflowByHandle(ctx, s, "base_handle")
				req.NoError(err)
				req.NotNil(wf)

				req.Equal("base_handle", wf.Handle)
				req.Equal("base_name", wf.Meta.Name)
				req.Equal("base_description", wf.Meta.Description)
				req.True(wf.Enabled)
				req.True(wf.Trace)
				req.Equal(10, wf.KeepSessions)
				req.Len(wf.Steps, 2)
				req.Len(wf.Paths, 1)
				req.Equal(createdAt.Format(time.RFC3339), wf.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), wf.UpdatedAt.Format(time.RFC3339))

				tt, _, err := store.SearchAutomationTriggers(ctx, s, atypes.TriggerFilter{
					WorkflowID: []uint64{wf.ID},
				})
				req.NoError(err)
				req.NotNil(tt)
				req.Len(tt, 1)
			},
		},

		{
			name: "base namespace",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				sTestComposeNamespace(ctx, t, s, "base")
				df := su.NewDecodeFilter().ComposeNamespace(&ctypes.NamespaceFilter{
					Slug: "base_namespace",
				})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)

				req.Equal("base namespace", n.Name)
				req.Equal("base_namespace", n.Slug)
				req.True(n.Enabled)
				req.Equal("subtitle", n.Meta.Subtitle)
				req.Equal("description", n.Meta.Description)
				req.Equal(createdAt.Format(time.RFC3339), n.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), n.UpdatedAt.Format(time.RFC3339))
			},
		},

		{
			name: "base module",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				sTestComposeModule(ctx, t, s, ns.ID, "base")

				df := su.NewDecodeFilter().
					ComposeNamespace(&ctypes.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&ctypes.ModuleFilter{
						NamespaceID: ns.ID,
						Handle:      "base_module",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)

				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, n.ID, "base_module")
				req.NoError(err)
				mff, _, err := store.SearchComposeModuleFields(ctx, s, ctypes.ModuleFieldFilter{
					ModuleID: []uint64{mod.ID},
				})
				req.NoError(err)

				// Module stuff
				req.Equal("base module", mod.Name)
				req.Equal("base_module", mod.Handle)
				req.Equal(n.ID, mod.NamespaceID)
				req.Equal(createdAt.Format(time.RFC3339), n.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), n.UpdatedAt.Format(time.RFC3339))

				// Module fields
				f := mff.FindByName("module_field_string")
				req.Equal(mod.ID, f.ModuleID)
				req.Equal(0, f.Place)
				req.Equal("String", f.Kind)
				req.Equal("module_field_string", f.Name)
				req.Equal("module field string", f.Label)
				req.Equal(true, f.Private)
				req.Equal(true, f.Required)
				req.Equal(true, f.Visible)
				req.Equal(true, f.Multi)
				req.Equal("opt_value_1", f.Options["opt1"])

				f = mff.FindByName("module_field_number")
				req.Equal(mod.ID, f.ModuleID)
				req.Equal(1, f.Place)
				req.Equal("Number", f.Kind)
				req.Equal("module_field_number", f.Name)
				req.Equal("module field number", f.Label)
				req.Equal(false, f.Private)
				req.Equal(false, f.Required)
				req.Equal(false, f.Visible)
				req.Equal(false, f.Multi)
				req.Equal("opt_value_1", f.Options["opt1"])
			},
		},

		{
			name: "base page",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				sTestComposePage(ctx, t, s, ns.ID, "base")

				df := su.NewDecodeFilter().
					ComposeNamespace(&ctypes.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposePage(&ctypes.PageFilter{
						NamespaceID: ns.ID,
						Handle:      "base_page",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				n, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)

				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, n.ID, "base_page")
				req.NoError(err)

				// Page
				req.Equal(n.ID, pg.NamespaceID)
				req.Equal("base_page", pg.Handle)
				req.Equal("base page", pg.Title)
				req.Equal("description", pg.Description)
				req.Equal(true, pg.Visible)
				req.Equal(0, pg.Weight)
				req.Len(pg.Blocks, 2)
				req.Equal(createdAt.Format(time.RFC3339), pg.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), pg.UpdatedAt.Format(time.RFC3339))

				// Blocks
				b := pg.Blocks[0]
				req.Equal("page block content", b.Title)
				req.Equal("description", b.Description)
				req.Equal("Content", b.Kind)

				b = pg.Blocks[1]
				req.Equal("page block qwerty", b.Title)
				req.Equal("description", b.Description)
				req.Equal("Qwerty", b.Kind)
			},
		},

		{
			name: "base chart",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")
				sTestComposeChart(ctx, t, s, ns.ID, mod.ID, "base")

				df := su.NewDecodeFilter().
					ComposeNamespace(&ctypes.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&ctypes.ModuleFilter{
						NamespaceID: ns.ID,
						Handle:      "base_module",
					}).
					ComposeChart(&ctypes.ChartFilter{
						NamespaceID: ns.ID,
						Handle:      "base_chart",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)

				chr, err := store.LookupComposeChartByNamespaceIDHandle(ctx, s, ns.ID, "base_chart")
				req.NoError(err)

				req.Equal(ns.ID, chr.NamespaceID)
				req.Equal("base_chart", chr.Handle)
				req.Equal("base chart", chr.Name)
				req.Equal(createdAt.Format(time.RFC3339), chr.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), chr.UpdatedAt.Format(time.RFC3339))

				req.Equal("colorscheme", chr.Config.ColorScheme)
				req.Len(chr.Config.Reports, 1)

				r := chr.Config.Reports[0]
				req.Equal("filter", r.Filter)
				req.Equal(mod.ID, r.ModuleID)

				req.Equal("Euro", r.YAxis["label"])
				req.Equal(true, r.YAxis["beginAtZero"])
			},
		},

		{
			name: "base record",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")
				usr := sTestUser(ctx, t, s, "base")
				sTestComposeRecord(ctx, t, s, ns.ID, mod.ID, usr.ID)

				df := su.NewDecodeFilter().
					ComposeNamespace(&ctypes.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&ctypes.ModuleFilter{
						NamespaceID: ns.ID,
						Handle:      "base_module",
					}).
					Users(&stypes.UserFilter{
						Email: "base_user@test.tld",
					}).
					ComposeRecord(&ctypes.RecordFilter{
						NamespaceID: ns.ID,
						ModuleID:    mod.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				rr, _, err := store.SearchComposeRecords(ctx, s, mod, ctypes.RecordFilter{
					ModuleID:    mod.ID,
					NamespaceID: ns.ID,
				})
				req.NoError(err)
				req.Len(rr, 1)
				rec := rr[0]

				req.Equal(ns.ID, rec.NamespaceID)
				req.Equal(mod.ID, rec.ModuleID)

				req.Equal(createdAt.Format(time.RFC3339), rec.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), rec.UpdatedAt.Format(time.RFC3339))
				req.Equal(usr.ID, rec.OwnedBy)
				req.Equal(usr.ID, rec.CreatedBy)
				req.Equal(usr.ID, rec.UpdatedBy)

				req.Len(rec.Values, 2)
				vv := rec.Values.FilterByName("module_field_string")
				req.Len(vv, 1)
				req.Equal("string value", vv[0].Value)

				vv = rec.Values.FilterByName("module_field_number")
				req.Len(vv, 1)
				req.Equal("10", vv[0].Value)
			},
		},

		{
			name: "full value record",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModuleFull(ctx, s, t, ns.ID, "base")
				usr := sTestUser(ctx, t, s, "base")

				recID := su.NextID()
				rec := &ctypes.Record{
					ID:          recID,
					NamespaceID: ns.ID,
					ModuleID:    mod.ID,

					Values: ctypes.RecordValueSet{
						{
							RecordID: recID,
							Name:     "BoolTrue",
							Value:    "1",
						},
						{
							RecordID: recID,
							Name:     "BoolFalse",
							Value:    "0",
						},
						{
							RecordID: recID,
							Name:     "DateTime",
							Value:    "2021-01-01T11:10:09Z",
						},
						{
							RecordID: recID,
							Name:     "Email",
							Value:    "test@mail.tld",
						},
						{
							RecordID: recID,
							Name:     "Select",
							Value:    "v1",
						},
						{
							RecordID: recID,
							Name:     "Number",
							Value:    "10.01",
						},
						{
							RecordID: recID,
							Name:     "String",
							Value:    "testing",
						},
						{
							RecordID: recID,
							Name:     "Url",
							Value:    "htts://www.testing.tld",
						},
						{
							RecordID: recID,
							Name:     "User",
							Value:    strconv.FormatUint(usr.ID, 10),
							Ref:      usr.ID,
						},
					},
				}
				err := store.CreateComposeRecord(ctx, s, mod, rec)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeNamespace(&ctypes.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&ctypes.ModuleFilter{
						NamespaceID: ns.ID,
						Handle:      "base_module",
					}).
					Users(&stypes.UserFilter{
						Email: "base_user@test.tld",
					}).
					ComposeRecord(&ctypes.RecordFilter{
						NamespaceID: ns.ID,
						ModuleID:    mod.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				rr, _, err := store.SearchComposeRecords(ctx, s, mod, ctypes.RecordFilter{
					ModuleID:    mod.ID,
					NamespaceID: ns.ID,
				})
				req.NoError(err)
				req.Len(rr, 1)
				rec := rr[0]

				req.Equal("1", rec.Values.FilterByName("BoolTrue")[0].Value)
				req.Equal("", rec.Values.FilterByName("BoolFalse")[0].Value)
				req.Equal("2021-01-01T11:10:09Z", rec.Values.FilterByName("DateTime")[0].Value)
				req.Equal("test@mail.tld", rec.Values.FilterByName("Email")[0].Value)
				req.Equal("v1", rec.Values.FilterByName("Select")[0].Value)
				req.Equal("10.01", rec.Values.FilterByName("Number")[0].Value)
				req.Equal("testing", rec.Values.FilterByName("String")[0].Value)
				req.Equal("htts://www.testing.tld", rec.Values.FilterByName("Url")[0].Value)
				req.Equal(strconv.FormatUint(usr.ID, 10), rec.Values.FilterByName("User")[0].Value)
				req.Equal(usr.ID, rec.Values.FilterByName("User")[0].Ref)
			},
		},

		{
			name: "base roles",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				sTestRole(ctx, t, s, "base")

				df := su.NewDecodeFilter().
					Roles(&stypes.RoleFilter{
						Handle: "base_role",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				rl, err := store.LookupRoleByHandle(ctx, s, "base_role")
				req.NoError(err)

				req.Equal("base role", rl.Name)
				req.Equal("base_role", rl.Handle)
				req.Equal(createdAt.Format(time.RFC3339), rl.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), rl.UpdatedAt.Format(time.RFC3339))
			},
		},

		{
			name: "base users",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				sTestUser(ctx, t, s, "base")

				df := su.NewDecodeFilter().
					Users(&stypes.UserFilter{
						Handle: "base_user",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				req.Equal("base_user_u", usr.Username)
				req.Equal("base_user@test.tld", usr.Email)
				req.Equal("base user", usr.Name)
				req.Equal("base_user", usr.Handle)
				req.Equal(stypes.NormalUser, usr.Kind)
				req.Equal("avatar", usr.Meta.Avatar)
				req.Equal(true, usr.EmailConfirmed)
				req.Equal(createdAt.Format(time.RFC3339), usr.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), usr.UpdatedAt.Format(time.RFC3339))
			},
		},

		{
			name: "base templates",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				sTestTemplate(ctx, t, s, "base")

				df := su.NewDecodeFilter().
					Templates(&stypes.TemplateFilter{
						Handle: "base_template",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				tpl, err := store.LookupTemplateByHandle(ctx, s, "base_template")
				req.NoError(err)

				req.Equal("base_template", tpl.Handle)
				req.Equal(stypes.DocumentTypeHTML, tpl.Type)
				req.Equal(true, tpl.Partial)
				req.Equal("base_short", tpl.Meta.Short)
				req.Equal("base_description", tpl.Meta.Description)
				req.Equal("base_template content", tpl.Template)
				req.Equal(createdAt.Format(time.RFC3339), tpl.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), tpl.UpdatedAt.Format(time.RFC3339))
			},
		},

		{
			name: "base applications",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				usr := sTestUser(ctx, t, s, "base")
				sTestApplication(ctx, t, s, usr.ID, "base")

				df := su.NewDecodeFilter().
					Users(&stypes.UserFilter{
						Handle: "base_user",
					}).
					Applications(&stypes.ApplicationFilter{
						Name: "base application",
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				aa, _, err := store.SearchApplications(ctx, s, stypes.ApplicationFilter{
					Name: "base application",
				})
				req.NoError(err)
				req.Len(aa, 1)
				app := aa[0]

				req.Equal("base application", app.Name)
				req.Equal(usr.ID, app.OwnerID)
				req.Equal(true, app.Enabled)
				req.Equal("name", app.Unify.Name)
				req.Equal(true, app.Unify.Listed)
				req.Equal("icon", app.Unify.Icon)
				req.Equal("logo", app.Unify.Logo)
				req.Equal("url", app.Unify.Url)
				req.Equal("{\"config\": \"config\"}", app.Unify.Config)
				req.Equal(createdAt.Format(time.RFC3339), app.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), app.UpdatedAt.Format(time.RFC3339))
			},
		},

		{
			name: "base settings",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				usr := sTestUser(ctx, t, s, "base")
				sTestSettings(ctx, t, s, usr.ID, "base")

				df := su.NewDecodeFilter().
					Users(&stypes.UserFilter{
						Handle: "base_user",
					}).
					Settings(&stypes.SettingsFilter{})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				ss, _, err := store.SearchSettings(ctx, s, stypes.SettingsFilter{})
				req.NoError(err)
				req.Len(ss, 2)

				sv := ss[0]
				req.Equal("base_setting_1", sv.Name)
				req.Equal("10", sv.Value.String())
				req.Equal(updatedAt.Format(time.RFC3339), sv.UpdatedAt.Format(time.RFC3339))
				req.Equal(usr.ID, sv.UpdatedBy)

				sv = ss[1]
				req.Equal("base_setting_2", sv.Name)
				req.Equal("20", sv.Value.String())
				req.Equal(updatedAt.Format(time.RFC3339), sv.UpdatedAt.Format(time.RFC3339))
				req.Equal(usr.ID, sv.UpdatedBy)
			},
		},

		{
			name: "base rbac",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				rl := sTestRole(ctx, t, s, "base")
				sTestRbac(ctx, t, s, rl.ID)

				df := su.NewDecodeFilter().
					Roles(&stypes.RoleFilter{
						Handle: "base_role",
					}).
					Rbac(&rbac.RuleFilter{})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				rl, err := store.LookupRoleByHandle(ctx, s, "base_role")
				req.NoError(err)

				rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
				req.NoError(err)
				req.Len(rr, 4)

				for _, r := range rr {
					switch r.Resource {
					case ctypes.ComponentRbacResource():
						req.Equal(rl.ID, r.RoleID)
						req.Equal("read", r.Operation)
						req.Equal(rbac.Allow, r.Access)
					case stypes.ComponentRbacResource():
						req.Equal(rl.ID, r.RoleID)
						req.Equal("read", r.Operation)
						req.Equal(rbac.Deny, r.Access)
					case stypes.RoleRbacResource(0):
						req.Equal(rl.ID, r.RoleID)
						req.Equal("read", r.Operation)
						req.Equal(rbac.Allow, r.Access)
					case stypes.UserRbacResource(rl.ID):
						req.Equal(rl.ID, r.RoleID)
						req.Equal(stypes.UserRbacResource(rl.ID), r.Resource)
						req.Equal("read", r.Operation)
						req.Equal(rbac.Deny, r.Access)
					}
				}
			},
		},

		{
			name: "ComposeRecord rbac",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				rl := sTestRole(ctx, t, s, "base")
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")
				rec := sTestComposeRecord(ctx, t, s, ns.ID, mod.ID, 0)

				rr := rbac.RuleSet{
					{
						RoleID:    rl.ID,
						Resource:  ctypes.RecordRbacResource(0, 0, 0),
						Operation: "allow.op1",
						Access:    rbac.Allow,
					},
					{
						RoleID:    rl.ID,
						Resource:  ctypes.RecordRbacResource(ns.ID, 0, 0),
						Operation: "allow.op2",
						Access:    rbac.Deny,
					},
					{
						RoleID:    rl.ID,
						Resource:  ctypes.RecordRbacResource(ns.ID, mod.ID, 0),
						Operation: "allow.op3",
						Access:    rbac.Allow,
					},
					{
						RoleID:    rl.ID,
						Resource:  ctypes.RecordRbacResource(ns.ID, mod.ID, rec.ID),
						Operation: "allow.op4",
						Access:    rbac.Allow,
					},
				}
				if err := store.CreateRbacRule(ctx, s, rr...); err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					Roles(&stypes.RoleFilter{
						Handle: "base_role",
					}).
					ComposeNamespace(&ctypes.NamespaceFilter{
						Slug: "base_namespace",
					}).
					ComposeModule(&ctypes.ModuleFilter{
						Handle: "base_module",
					}).
					ComposeRecord(&ctypes.RecordFilter{
						NamespaceID: ns.ID,
						ModuleID:    mod.ID,
					}).
					Rbac(&rbac.RuleFilter{})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				rl, err := store.LookupRoleByHandle(ctx, s, "base_role")
				req.NoError(err)

				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)

				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)

				rcs, _, err := store.SearchComposeRecords(ctx, s, mod, ctypes.RecordFilter{})
				req.NoError(err)
				rec := rcs[0]

				rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
				req.NoError(err)
				req.Len(rr, 4)

				for _, r := range rr {
					switch r.Operation {
					case "allow.op1":
						req.Equal(rl.ID, r.RoleID)
						req.Equal(rbac.Allow, r.Access)
						req.Equal(ctypes.RecordRbacResource(0, 0, 0), r.Resource)

					case "allow.op2":
						req.Equal(rl.ID, r.RoleID)
						req.Equal(rbac.Deny, r.Access)
						req.Equal(ctypes.RecordRbacResource(ns.ID, 0, 0), r.Resource)

					case "allow.op3":
						req.Equal(rl.ID, r.RoleID)
						req.Equal(rbac.Allow, r.Access)
						req.Equal(ctypes.RecordRbacResource(ns.ID, mod.ID, 0), r.Resource)

					case "allow.op4":
						req.Equal(rl.ID, r.RoleID)
						req.Equal(rbac.Allow, r.Access)
						req.Equal(ctypes.RecordRbacResource(ns.ID, mod.ID, rec.ID), r.Resource)

					default:
						req.FailNow("unexpected rbac operation for test cases: ", r.Operation)
					}
				}
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

		})
		ni = 0
	}

	// Cleanup the store
	truncateStore(ctx, s, t)
}
