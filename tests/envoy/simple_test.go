package envoy

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	st "github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestSimpleCases(t *testing.T) {
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
			name:  "simple namespaces",
			suite: "simple",
			file:  "namespaces",
			pre: func() (err error) {
				return s.TruncateComposeNamespaces(ctx)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
				req.NoError(err)
				req.NotNil(ns)
				req.Equal("ns1", ns.Slug)
				req.Equal("ns1 name", ns.Name)
			},
		},

		{
			name:  "simple mods; no namespace",
			suite: "simple",
			file:  "modules_no_ns",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeModuleFields(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)
				req.True(strings.Contains(err.Error(), "prepare compose module"))
				req.True(strings.Contains(err.Error(), "compose namespace unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "simple mods",
			suite: "simple",
			file:  "modules",
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
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, 100, "mod1")
				req.NoError(err)
				req.NotNil(mod)

				mod.Fields, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
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
			name:  "simple mods; conditional",
			suite: "simple",
			file:  "modules_conditional",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeRecords(ctx, nil),
					s.TruncateComposeModuleFields(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeNamespaces(ctx),

					storeNamespace(ctx, s, 100, "ns1"),

					storeModule(ctx, s, 100, 200, "mod1", "mod1 name"),
					storeModuleField(ctx, s, 200, 300, "f1"),

					storeModule(ctx, s, 100, 201, "mod2", "mod2 name"),
					storeModuleField(ctx, s, 201, 301, "f1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
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
			name:  "simple charts; no ns",
			suite: "simple",
			file:  "charts_no_ns",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.True(strings.Contains(err.Error(), "prepare compose chart"))
				req.True(strings.Contains(err.Error(), "compose namespace unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "simple charts; no mod",
			suite: "simple",
			file:  "charts_no_mod",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.True(strings.Contains(err.Error(), "prepare compose chart"))
				req.True(strings.Contains(err.Error(), "compose module unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "simple charts",
			suite: "simple",
			file:  "charts",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
					storeModule(ctx, s, 100, 200, "mod1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
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
			name:  "simple pages; no ns",
			suite: "simple",
			file:  "pages_no_ns",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),
					s.TruncateComposePages(ctx),
				)
			},

			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.True(strings.Contains(err.Error(), "prepare compose page"))
				req.True(strings.Contains(err.Error(), "compose namespace unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "record pages; no mod",
			suite: "simple",
			file:  "pages_r_no_mod",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),
					s.TruncateComposePages(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
				)
			},

			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.True(strings.Contains(err.Error(), "prepare compose page"))
				req.True(strings.Contains(err.Error(), "compose module unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "simple pages",
			suite: "simple",
			file:  "pages",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
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
			name:  "record page",
			suite: "simple",
			file:  "pages_r",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeCharts(ctx),
					s.TruncateComposePages(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
					storeModule(ctx, s, 100, 200, "mod1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
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
			name:  "applications",
			suite: "simple",
			file:  "applications",
			pre: func() (err error) {
				return ce(
					s.TruncateApplications(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},

			check: func(req *require.Assertions) {
				apps, _, err := store.SearchApplications(ctx, s, st.ApplicationFilter{
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
			name:  "users",
			suite: "simple",
			file:  "users",
			pre: func() (err error) {
				return ce(
					s.TruncateUsers(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
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
			name:  "roles",
			suite: "simple",
			file:  "roles",
			pre: func() (err error) {
				return ce(
					s.TruncateRoles(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				r, err := store.LookupRoleByHandle(ctx, s, "r1")
				req.NoError(err)
				req.NotNil(r)

				req.Equal("r1", r.Handle)
				req.Equal("r1 name", r.Name)
			},
		},

		{
			name:  "settings",
			suite: "simple",
			file:  "settings",
			pre: func() (err error) {
				return ce(
					s.TruncateSettings(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				ss, _, err := store.SearchSettings(ctx, s, st.SettingsFilter{})
				req.NoError(err)
				req.NotNil(ss)
				req.Len(ss, 3)
			},
		},

		{
			name:  "rbac rules; no role",
			suite: "simple",
			file:  "rbac_rules_no_role",
			pre: func() (err error) {
				return ce(
					s.TruncateRoles(ctx),
					s.TruncateRbacRules(ctx),
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeModuleFields(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)
				req.True(strings.Contains(err.Error(), "prepare rbac rule"))
				req.True(strings.Contains(err.Error(), "role unresolved"))
			},
			check: func(req *require.Assertions) {
			},
		},

		{
			name:  "rbac rules",
			suite: "simple",
			file:  "rbac_rules",
			pre: func() (err error) {
				return ce(
					s.TruncateRoles(ctx),
					s.TruncateRbacRules(ctx),
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeModuleFields(ctx),

					storeRole(ctx, s, 100, "r1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
				req.NoError(err)
				req.NotNil(rr)
				req.Len(rr, 3)
			},
		},

		{
			name:  "simple records; no ns",
			suite: "simple",
			file:  "records_no_ns",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.True(strings.Contains(err.Error(), "prepare compose record"))
				req.True(strings.Contains(err.Error(), "compose namespace unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "simple records; no mod",
			suite: "simple",
			file:  "records_no_mod",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.Error(err)

				req.True(strings.Contains(err.Error(), "prepare compose record"))
				req.True(strings.Contains(err.Error(), "compose module unresolved"))
			},
			check: func(req *require.Assertions) {},
		},

		{
			name:  "simple records",
			suite: "simple",
			file:  "records",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),

					storeNamespace(ctx, s, 100, "ns1"),
					storeModule(ctx, s, 100, 200, "mod1"),
					storeModuleField(ctx, s, 200, 300, "f1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				m, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(m)

				rr, _, err := store.SearchComposeRecords(ctx, s, m, types.RecordFilter{ModuleID: m.ID, NamespaceID: m.NamespaceID})
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
			name:  "simple records; multiple",
			suite: "simple",
			file:  "records_multi",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeModuleFields(ctx),
					s.TruncateComposeRecords(ctx, nil),

					storeNamespace(ctx, s, 100, "ns1"),
					storeModule(ctx, s, 100, 200, "mod1"),
					storeModuleField(ctx, s, 200, 300, "f1"),

					storeModule(ctx, s, 100, 201, "mod2"),
					storeModuleField(ctx, s, 201, 301, "f1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				mod1, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(mod1)

				mod2, err := store.LookupComposeModuleByID(ctx, s, 201)
				req.NoError(err)
				req.NotNil(mod2)

				rr, _, err := store.SearchComposeRecords(ctx, s, mod1, types.RecordFilter{ModuleID: mod1.ID, NamespaceID: mod1.NamespaceID})
				req.NoError(err)
				req.NotNil(rr)
				req.Len(rr, 1)
				req.Equal("mod1 f1 v1", rr[0].Values[0].Value)

				rr, _, err = store.SearchComposeRecords(ctx, s, mod2, types.RecordFilter{ModuleID: mod2.ID, NamespaceID: mod2.NamespaceID})
				req.NoError(err)
				req.NotNil(rr)
				req.Len(rr, 1)
				req.Equal("mod2 f1 v1", rr[0].Values[0].Value)
			},
		},

		{
			name:  "simple records; conditional",
			suite: "simple",
			file:  "records_conditional",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeRecords(ctx, nil),
					s.TruncateComposeModuleFields(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeNamespaces(ctx),

					storeNamespace(ctx, s, 100, "ns1"),

					storeModule(ctx, s, 100, 200, "mod1"),
					storeModuleField(ctx, s, 200, 300, "f1"),
					storeRecord(ctx, s, 100, 200, 400, "existing value"),

					storeModule(ctx, s, 100, 201, "mod2"),
					storeModuleField(ctx, s, 201, 301, "f1"),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {
				mod1, err := store.LookupComposeModuleByID(ctx, s, 200)
				req.NoError(err)
				req.NotNil(mod1)

				mod2, err := store.LookupComposeModuleByID(ctx, s, 201)
				req.NoError(err)
				req.NotNil(mod2)

				rr1, _, err := store.SearchComposeRecords(ctx, s, mod1, types.RecordFilter{ModuleID: mod1.ID, NamespaceID: mod1.NamespaceID})
				req.NoError(err)
				req.NotNil(rr1)
				req.Len(rr1, 1)
				req.Equal("existing value", rr1[0].Values.FilterByName("f1")[0].Value)

				rr2, _, err := store.SearchComposeRecords(ctx, s, mod2, types.RecordFilter{ModuleID: mod2.ID, NamespaceID: mod2.NamespaceID})
				req.NoError(err)
				req.NotNil(rr2)
				req.Len(rr2, 1)
				req.Equal("f1 value", rr2[0].Values.FilterByName("f1")[0].Value)
			},
		},

		{
			name:  "base access control",
			suite: "simple",
			file:  "access_control_base",
			pre: func() (err error) {
				return ce(
					s.TruncateComposeRecords(ctx, nil),
					s.TruncateComposeModuleFields(ctx),
					s.TruncateComposeModules(ctx),
					s.TruncateComposeNamespaces(ctx),
					s.TruncateComposePages(ctx),
					s.TruncateComposeCharts(ctx),
					s.TruncateUsers(ctx),
					s.TruncateRoles(ctx),
					s.TruncateMessagingChannels(ctx),
					s.TruncateRbacRules(ctx),
				)
			},
			post: func(req *require.Assertions, err error) {
				req.NoError(err)
			},
			check: func(req *require.Assertions) {

				role, err := store.LookupRoleByHandle(ctx, s, "everyone")
				req.NoError(err)
				req.NotNil(role)

				rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
				req.NoError(err)
				req.Len(rr, 18)

				// Check that the role is ok
				rr.Walk(func(r *rbac.Rule) error {
					req.Equal(role.ID, r.RoleID)
					return nil
				})

				resources := []string{
					"compose:namespace:",
					"compose:namespace:",
					"compose:module:",
					"compose:module:",
					"compose:page:",
					"compose:page:",
					"compose:chart:",
					"compose:chart:",
					"messaging:channel:",
					"messaging:channel:",
					"system:role:",
					"system:role:",
					"system:user:",
					"system:user:",
					"system:application:",
					"system:application:",
					"compose",
					"compose",
				}

				for i, res := range resources {
					req.Equal(rbac.Resource(res), rr[i].Resource.TrimID())
					if i%2 == 0 {
						req.Equal(rbac.Operation("op1"), rr[i].Operation)
						req.Equal(rbac.Allow, rr[i].Access)
					} else {
						req.Equal(rbac.Operation("op2"), rr[i].Operation)
						req.Equal(rbac.Deny, rr[i].Access)
					}
				}
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
