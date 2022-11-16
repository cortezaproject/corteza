package envoy

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"testing"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	ctypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza/server/pkg/envoy/store"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	systypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestYamlStore_base(t *testing.T) {
	type (
		tc struct {
			name  string
			file  string
			asDir bool

			// Before the data gets processed
			pre func() error
			// After the data gets processed
			post func(req *require.Assertions, err error)
			// Data assertions
			check func(req *require.Assertions)
		}
	)

	var (
		ctx       = context.Background()
		namespace = "base"
		s         = initServices(ctx, t)
		err       error
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	// Allow rec. field values to be updated for this user.
	rbac.Global().Grant(
		ctx,
		rbac.AllowRule(
			auth.ServiceUser().Roles()[0],
			ctypes.ModuleFieldRbacResource(0, 0, 0),
			"record.value.update",
		),
	)

	cases := []*tc{
		{
			name: "namespaces",
			file: "namespaces",
			check: func(req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
				req.NoError(err)
				req.NotNil(ns)
				req.Equal("ns1", ns.Slug)
				req.Equal("ns1 name", ns.Name)
			},
		},

		{
			name: "modules; no namespace",
			file: "modules_no_ns",
			post: func(req *require.Assertions, err error) {
				req.Error(err)
				req.Contains(err.Error(), "prepare corteza::compose:module")
				req.Contains(err.Error(), "compose namespace unresolved")
			},
		},

		{
			name: "modules",
			file: "modules",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
				)
			},
			check: func(req *require.Assertions) {
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 100, "mod1")
				req.NoError(err)
				req.NotNil(mod)

				mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, ctypes.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
				req.NoError(err)

				req.Equal("mod1", mod.Handle)
				req.Equal("mod1 name", mod.Name)
				req.Equal(uint64(100), mod.NamespaceID)

				ff := mod.Fields
				req.Len(ff, 1)
				req.Equal("f1", ff[0].Name)
				req.Equal("f1 label", ff[0].Label)
				req.Equal("String", ff[0].Kind)
			},
		},

		{
			name: "modules; conditional",
			file: "modules_conditional",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),

					storeComposeModule(ctx, s, 100, 200, "mod1", "mod1 name"),
					storeComposeModuleField(ctx, s, 200, 300, "f1"),

					storeComposeModule(ctx, s, 100, 201, "mod2", "mod2 name"),
					storeComposeModuleField(ctx, s, 201, 301, "f1"),
				)
			},
			check: func(req *require.Assertions) {
				mod1, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(mod1)

				mod2, err := store.LookupComposeModuleByID(ctx, s, 201)
				req.NoError(err)
				req.NotNil(mod2)

				// The first one overwrites merge alg to replace, the second one defaults to skip
				req.Equal("mod1 name (EDITED)", mod1.Name)
				req.Equal("mod2 name", mod2.Name)
			},
		},

		{
			name: "modules; expressions",
			file: "modules_expressions",
			check: func(req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "crm")
				req.NoError(err)
				req.NotNil(ns)

				mod, err := loadComposeModuleFull(ctx, s, req, ns.ID, "Account")
				req.NotNil(mod)
				req.NoError(err)
				req.Len(mod.Fields, 2)

				// Check the full thing
				mfF := mod.Fields[0]
				req.Equal("a > b", mfF.Expressions.ValueExpr)
				req.Subset(mfF.Expressions.Sanitizers, []string{"trim(value)"})
				v := mfF.Expressions.Validators[0]
				req.Equal("a == \"\"", v.Test)
				req.Equal("Value should not be empty", v.Error)

				// Check the other validator form
				mfV := mod.Fields[1]
				v = mfV.Expressions.Validators[0]
				req.Equal("value == \"\"", v.Test)
				req.Equal("Value should be filled", v.Error)
			},
		},

		{
			name: "charts; no namespace",
			file: "charts_no_ns",
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.Contains(err.Error(), "prepare corteza::compose:chart")
				req.Contains(err.Error(), "compose namespace unresolved")
			},
		},

		{
			name: "charts; no module",
			file: "charts_no_mod",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.Contains(err.Error(), "prepare corteza::compose:chart")
				req.Contains(err.Error(), "compose module unresolved")
			},
		},

		{
			name: "charts",
			file: "charts",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
					storeComposeModule(ctx, s, 100, 200, "mod1"),
				)
			},
			check: func(req *require.Assertions) {
				chr, err := store.LookupComposeChartByNamespaceIDHandle(ctx, s, 100, "c1")
				req.NoError(err)
				req.NotNil(chr)

				req.Equal("c1", chr.Handle)
				req.Equal("c1 name", chr.Name)

				req.Len(chr.Config.Reports, 1)

				req.Equal(uint64(200), chr.Config.Reports[0].ModuleID)
			},
		},

		{
			name: "charts; axis",
			file: "charts_axis",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
					storeComposeModule(ctx, s, 100, 200, "mod1"),
				)
			},
			check: func(req *require.Assertions) {
				chr, err := store.LookupComposeChartByNamespaceIDHandle(ctx, s, 100, "c1")
				req.NoError(err)
				req.NotNil(chr)

				req.Equal("c1", chr.Handle)
				req.Equal("c1 name", chr.Name)

				req.Len(chr.Config.Reports, 1)

				rep := chr.Config.Reports[0]
				req.Equal("Euro", rep.YAxis["label"])
				req.Equal(true, rep.YAxis["beginAtZero"])
			},
		},

		{
			name: "pages; no namespace",
			file: "pages_no_ns",
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.Contains(err.Error(), "prepare corteza::compose:page")
				req.Contains(err.Error(), "compose namespace unresolved")
			},
		},

		{
			name: "record pages; no module",
			file: "pages_r_no_mod",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
				)
			},

			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.Contains(err.Error(), "prepare corteza::compose:page")
				req.Contains(err.Error(), "compose module unresolved")
			},
		},

		{
			name: "pages",
			file: "pages",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
				)
			},
			check: func(req *require.Assertions) {
				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, 100, "pg1")
				req.NoError(err)
				req.NotNil(pg)

				req.Equal("pg1", pg.Handle)
				req.Equal("pg1 name", pg.Title)
				req.Equal(uint64(0), pg.ModuleID)
				req.Equal(uint64(100), pg.NamespaceID)

				req.Len(pg.Blocks, 1)
				req.Equal("block1", pg.Blocks[0].Title)
			},
		},

		{
			name: "record page",
			file: "pages_r",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
					storeComposeModule(ctx, s, 100, 200, "mod1"),
				)
			},
			check: func(req *require.Assertions) {
				pg, err := store.LookupComposePageByNamespaceIDHandle(ctx, s, 100, "pg1")
				req.NoError(err)
				req.NotNil(pg)

				req.Equal("pg1", pg.Handle)
				req.Equal("pg1 name", pg.Title)
				req.Equal(uint64(100), pg.NamespaceID)
				req.Equal(uint64(200), pg.ModuleID)

				req.Len(pg.Blocks, 1)
				req.Equal("block1", pg.Blocks[0].Title)
			},
		},

		{
			name: "applications",
			file: "applications",
			check: func(req *require.Assertions) {
				apps, _, err := store.SearchApplications(ctx, s, systypes.ApplicationFilter{
					Name: "app1",
				})
				req.NoError(err)
				req.NotNil(apps)
				req.Len(apps, 1)
				app := apps[0]

				req.Equal("app1", app.Name)
			},
		},

		{
			name: "users",
			file: "users",
			check: func(req *require.Assertions) {
				u, err := store.LookupUserByHandle(ctx, s, "u1")
				req.NoError(err)
				req.NotNil(u)

				req.Equal("u1", u.Handle)
				req.Equal("u1 name", u.Name)
				req.Equal("u1@example.tld", u.Email)
			},
		},

		{
			name: "workflows",
			file: "workflows",
			check: func(req *require.Assertions) {
				ww, err := store.LookupAutomationWorkflowByHandle(ctx, s, "testko_wf")
				req.NoError(err)
				req.NotNil(ww)

				req.Equal("testko_wf", ww.Handle)
				req.Equal("name here", ww.Meta.Name)
				req.Equal("description here", ww.Meta.Description)
				req.Len(ww.Steps, 2)
				req.Len(ww.Paths, 1)

				tt, _, err := store.SearchAutomationTriggers(ctx, s, atypes.TriggerFilter{
					WorkflowID: []uint64{ww.ID},
				})
				req.NoError(err)
				req.NotNil(tt)
				req.Len(tt, 1)
			},
		},

		{
			name: "roles",
			file: "roles",
			check: func(req *require.Assertions) {
				r, err := store.LookupRoleByHandle(ctx, s, "r1")
				req.NoError(err)
				req.NotNil(r)

				req.Equal("r1", r.Handle)
				req.Equal("r1 name", r.Name)
			},
		},

		{
			name: "settings",
			file: "settings",
			check: func(req *require.Assertions) {
				ss, _, err := store.SearchSettings(ctx, s, systypes.SettingsFilter{})
				req.NoError(err)
				req.NotNil(ss)
				req.Len(ss, 3)
			},
		},

		{
			name: "rbac rules; no role",
			file: "rbac_rules_no_role",
			post: func(req *require.Assertions, err error) {
				req.Error(err)
				req.Contains(err.Error(), "prepare rbac-rule")
				req.Contains(err.Error(), "role unresolved")
			},
		},

		{
			name: "records; no namespace",
			file: "records_no_ns",
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.Contains(err.Error(), "prepare corteza::compose:record")
				req.Contains(err.Error(), "compose namespace unresolved")
			},
		},

		{
			name: "records; no module",
			file: "records_no_mod",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.Contains(err.Error(), "prepare corteza::compose:record")
				req.Contains(err.Error(), "compose module unresolved")
			},
		},

		{
			name: "records",
			file: "records",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
					storeComposeModule(ctx, s, 100, 200, "mod1"),
					storeComposeModuleField(ctx, s, 200, 300, "f1"),
				)
			},
			check: func(req *require.Assertions) {
				m, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(m)

				rr, _, err := store.SearchComposeRecords(ctx, s, m, ctypes.RecordFilter{ModuleID: m.ID, NamespaceID: m.NamespaceID})
				req.NoError(err)
				req.NotNil(rr)
				req.Len(rr, 1)

				r := rr[0]
				req.Equal(uint64(100), r.NamespaceID)
				req.Equal(uint64(200), r.ModuleID)

				req.Len(r.Values, 1)
				v := r.Values[0]
				req.Equal("v1", v.Value)
			},
		},

		{
			name: "records; multiple",
			file: "records_multi",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),
					storeComposeModule(ctx, s, 100, 200, "mod1"),
					storeComposeModuleField(ctx, s, 200, 300, "f1"),

					storeComposeModule(ctx, s, 100, 201, "mod2"),
					storeComposeModuleField(ctx, s, 201, 301, "f1"),
				)
			},
			check: func(req *require.Assertions) {
				mod1, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(mod1)

				mod2, err := store.LookupComposeModuleByID(ctx, s, 201)
				req.NoError(err)
				req.NotNil(mod2)

				rr, _, err := store.SearchComposeRecords(ctx, s, mod1, ctypes.RecordFilter{ModuleID: mod1.ID, NamespaceID: mod1.NamespaceID})
				req.NoError(err)
				req.NotNil(rr)
				req.Len(rr, 1)
				req.Equal("mod1 f1 v1", rr[0].Values[0].Value)

				rr, _, err = store.SearchComposeRecords(ctx, s, mod2, ctypes.RecordFilter{ModuleID: mod2.ID, NamespaceID: mod2.NamespaceID})
				req.NoError(err)
				req.NotNil(rr)
				req.Len(rr, 1)
				req.Equal("mod2 f1 v1", rr[0].Values[0].Value)
			},
		},

		{
			name: "records; conditional",
			file: "records_conditional",
			pre: func() (err error) {
				return collect(
					storeComposeNamespace(ctx, s, 100, "ns1"),

					storeComposeModule(ctx, s, 100, 200, "mod1"),
					storeComposeModuleField(ctx, s, 200, 300, "f1"),
					storeComposeRecord(ctx, s, 100, 200, 400, "existing value"),

					storeComposeModule(ctx, s, 100, 201, "mod2"),
					storeComposeModuleField(ctx, s, 201, 301, "f1"),
				)
			},
			check: func(req *require.Assertions) {
				mod1, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(mod1)

				mod2, err := store.LookupComposeModuleByID(ctx, s, 201)
				req.NoError(err)
				req.NotNil(mod2)

				rr1, _, err := store.SearchComposeRecords(ctx, s, mod1, ctypes.RecordFilter{ModuleID: mod1.ID, NamespaceID: mod1.NamespaceID})
				req.NoError(err)
				req.NotNil(rr1)
				req.Len(rr1, 1)
				req.Equal("existing value", rr1[0].Values.FilterByName("f1")[0].Value)

				rr2, _, err := store.SearchComposeRecords(ctx, s, mod2, ctypes.RecordFilter{ModuleID: mod2.ID, NamespaceID: mod2.NamespaceID})
				req.NoError(err)
				req.NotNil(rr2)
				req.Len(rr2, 1)
				req.Equal("f1 value", rr2[0].Values.FilterByName("f1")[0].Value)
			},
		},

		{
			name:  "settings",
			file:  "settings",
			asDir: true,
			check: func(req *require.Assertions) {
				ss, _, err := store.SearchSettings(ctx, s, systypes.SettingsFilter{})
				req.NoError(err)
				req.NotNil(ss)
				req.Len(ss, 6)

				rs := []string{ss[0].Name, ss[1].Name, ss[2].Name, ss[3].Name, ss[4].Name, ss[5].Name}
				req.Subset(rs, []string{"s1.opt.1", "s1.opt.2", "compose.ui.record-toolbar", "compose.ui.record-otherbar", "s2.opt.1", "s2.opt.2"})

				req.Equal("{\"hideBack\":true}", ss[2].Value.String())
				req.Equal("{\"hideBack\":true,\"hideCreate\":true}", ss[3].Value.String())
			},
		},
	}

	for _, c := range cases {
		f := c.file + ".yaml"
		t.Run(fmt.Sprintf("%s; testdata/%s/%s", c.name, namespace, f), func(t *testing.T) {
			// before running tests, reset ID and truncate store
			ni = 10
			truncateStore(ctx, s, t)

			req := require.New(t)

			if c.pre != nil {
				err = c.pre()
				req.NoError(err)
			}

			var nn []resource.Interface
			var err error
			if c.asDir {
				nn, err = decodeDirectory(ctx, path.Join(namespace, c.file))
			} else {
				nn, err = decodeYaml(ctx, namespace, f)
			}
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

		})
	}
}

func TestYamlStore_ComposeRecordRBAC(t *testing.T) {
	var (
		ctx       = context.Background()
		namespace = "base"
		s         = initServices(ctx, t)
		err       error
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	// Allow rec. field values to be updated for this user.
	rbac.Global().Grant(
		ctx,
		rbac.AllowRule(
			auth.ServiceUser().Roles()[0],
			ctypes.ModuleFieldRbacResource(0, 0, 0),
			"record.value.update",
		),
	)

	f := "rbac_rules_compose_records.yaml"
	t.Run(fmt.Sprintf("testdata/%s/%s", namespace, f), func(t *testing.T) {
		// before running tests, reset ID and truncate store
		ni = 10
		truncateStore(ctx, s, t)

		req := require.New(t)

		var nn []resource.Interface
		nn, err = decodeYaml(ctx, namespace, f)
		req.NoError(err)

		err = encode(ctx, s, nn)
		req.NoError(err)

		rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		req.NoError(err)

		// compose stuff
		ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
		req.NoError(err)

		mod1, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "mod1")
		req.NoError(err)

		rcs, _, err := store.SearchComposeRecords(ctx, s, mod1, ctypes.RecordFilter{})
		rec1 := rcs[0]
		req.NoError(err)

		role1, err := store.LookupRoleByHandle(ctx, s, "r1")
		req.NoError(err)
		role2, err := store.LookupRoleByHandle(ctx, s, "r2")
		req.NoError(err)

		// size
		req.Len(rr, 8)
		req.Len(rr.FilterAccess(rbac.Allow), 4)
		req.Len(rr.FilterAccess(rbac.Deny), 4)

		// allows
		var rule *rbac.Rule
		const tpl = "corteza::compose:record/%s/%s/%s"

		// allows
		rule = rr[0]
		req.Equal(role1.ID, rule.RoleID)
		req.Equal("allow.op1", rule.Operation)
		req.Equal(rbac.Allow, rule.Access)
		req.Equal(fmt.Sprintf(tpl, "*", "*", "*"), rule.Resource)

		rule = rr[1]
		req.Equal(role1.ID, rule.RoleID)
		req.Equal("allow.op2", rule.Operation)
		req.Equal(rbac.Allow, rule.Access)
		req.Equal(fmt.Sprintf(tpl, strconv.FormatUint(ns.ID, 10), "*", "*"), rule.Resource)

		rule = rr[2]
		req.Equal(role1.ID, rule.RoleID)
		req.Equal("allow.op3", rule.Operation)
		req.Equal(rbac.Allow, rule.Access)
		req.Equal(fmt.Sprintf(tpl, strconv.FormatUint(ns.ID, 10), strconv.FormatUint(mod1.ID, 10), "*"), rule.Resource)

		rule = rr[3]
		req.Equal(role1.ID, rule.RoleID)
		req.Equal("allow.op4", rule.Operation)
		req.Equal(rbac.Allow, rule.Access)
		req.Equal(fmt.Sprintf(tpl, strconv.FormatUint(ns.ID, 10), strconv.FormatUint(mod1.ID, 10), strconv.FormatUint(rec1.ID, 10)), rule.Resource)

		// denies
		rule = rr[4]
		req.Equal(role2.ID, rule.RoleID)
		req.Equal("deny.op1", rule.Operation)
		req.Equal(rbac.Deny, rule.Access)
		req.Equal(fmt.Sprintf(tpl, "*", "*", "*"), rule.Resource)

		rule = rr[5]
		req.Equal(role2.ID, rule.RoleID)
		req.Equal("deny.op2", rule.Operation)
		req.Equal(rbac.Deny, rule.Access)
		req.Equal(fmt.Sprintf(tpl, strconv.FormatUint(ns.ID, 10), "*", "*"), rule.Resource)

		rule = rr[6]
		req.Equal(role2.ID, rule.RoleID)
		req.Equal("deny.op3", rule.Operation)
		req.Equal(rbac.Deny, rule.Access)
		req.Equal(fmt.Sprintf(tpl, strconv.FormatUint(ns.ID, 10), strconv.FormatUint(mod1.ID, 10), "*"), rule.Resource)

		rule = rr[7]
		req.Equal(role2.ID, rule.RoleID)
		req.Equal("deny.op4", rule.Operation)
		req.Equal(rbac.Deny, rule.Access)
		req.Equal(fmt.Sprintf(tpl, strconv.FormatUint(ns.ID, 10), strconv.FormatUint(mod1.ID, 10), strconv.FormatUint(rec1.ID, 10)), rule.Resource)
	})
}
