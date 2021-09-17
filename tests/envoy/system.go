package envoy

import (
	"context"
	"testing"

	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func sTestUser(ctx context.Context, t *testing.T, s store.Storer, pfx string) *types.User {
	usr := &types.User{
		ID:       su.NextID(),
		Username: pfx + "_user_u",
		Email:    pfx + "_user@test.tld",
		Name:     pfx + " user",
		Handle:   pfx + "_user",
		Kind:     types.NormalUser,
		Meta: &types.UserMeta{
			Avatar: "avatar",
		},

		EmailConfirmed: true,

		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateUser(ctx, s, usr)
	if err != nil {
		t.Fatal(err)
	}

	return usr
}

func sTestTemplate(ctx context.Context, t *testing.T, s store.Storer, pfx string) *types.Template {
	tpl := &types.Template{
		ID:      su.NextID(),
		Handle:  pfx + "_template",
		Type:    types.DocumentTypeHTML,
		Partial: true,
		Meta: types.TemplateMeta{
			Short:       pfx + "_short",
			Description: pfx + "_description",
		},
		Template:  pfx + "_template content",
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateTemplate(ctx, s, tpl)
	if err != nil {
		t.Fatal(err)
	}

	return tpl
}

func sTestRole(ctx context.Context, t *testing.T, s store.Storer, pfx string) *types.Role {
	rl := &types.Role{
		ID:        su.NextID(),
		Name:      pfx + " role",
		Handle:    pfx + "_role",
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateRole(ctx, s, rl)
	if err != nil {
		t.Fatal(err)
	}

	return rl
}

func sTestApplication(ctx context.Context, t *testing.T, s store.Storer, usrID uint64, pfx string) *types.Application {
	app := &types.Application{
		ID:      su.NextID(),
		Name:    pfx + " application",
		OwnerID: usrID,
		Enabled: true,
		Unify: &types.ApplicationUnify{
			Name:   "name",
			Listed: true,
			Icon:   "icon",
			Logo:   "logo",
			Url:    "url",
			Config: "{\"config\": \"config\"}",
		},
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateApplication(ctx, s, app)
	if err != nil {
		t.Fatal(err)
	}

	return app
}

func sTestSettings(ctx context.Context, t *testing.T, s store.Storer, usrID uint64, pfx string) types.SettingValueSet {
	ss := types.SettingValueSet{
		{
			Name:      pfx + "_setting_1",
			Value:     []byte(`10`),
			UpdatedAt: updatedAt,
			UpdatedBy: usrID,
		},
		{
			Name:      pfx + "_setting_2",
			Value:     []byte(`20`),
			UpdatedAt: updatedAt,
			UpdatedBy: usrID,
		},
	}

	err := store.CreateSetting(ctx, s, ss...)
	if err != nil {
		t.Fatal(err)
	}

	return ss
}

func sTestRbac(ctx context.Context, t *testing.T, s store.Storer, roleID uint64) rbac.RuleSet {
	rr := rbac.RuleSet{
		{
			RoleID:    roleID,
			Resource:  "corteza::compose/",
			Operation: "read",
			Access:    rbac.Allow,
		},
		{
			RoleID:    roleID,
			Resource:  types.ComponentRbacResource(),
			Operation: "read",
			Access:    rbac.Deny,
		},
		{
			RoleID:    roleID,
			Resource:  types.RoleRbacResource(0),
			Operation: "read",
			Access:    rbac.Allow,
		},
		{
			RoleID:    roleID,
			Resource:  types.RoleRbacResource(roleID),
			Operation: "read",
			Access:    rbac.Deny,
		},
	}

	err := store.CreateRbacRule(ctx, s, rr...)
	if err != nil {
		t.Fatal(err)
	}

	return rr
}

func sTestReport(ctx context.Context, t *testing.T, s store.Storer, pfx string) *types.Report {
	r := &types.Report{
		ID:     su.NextID(),
		Handle: pfx + "_report",
		Meta:   &types.ReportMeta{Name: pfx + " report", Description: "testing"},

		Sources: types.ReportDataSourceSet{{
			Meta: map[string]interface{}{"key1": "value1"},
			Step: &report.StepDefinition{
				Kind: "Load",
				Load: &report.LoadStepDefinition{
					Name:   "Test",
					Source: "test",
					Definition: map[string]interface{}{
						"k1": "v1",
					},
					Columns: report.FrameColumnSet{{Name: "col1", Label: "col1 label"}},
				},
			},
		}},
		Projections: types.ReportProjectionSet{{
			Title:       "title",
			Description: "description",
			Key:         "key",
			Kind:        "kind",
			Options: map[string]interface{}{
				"k1": "v1",
			},
			Elements: []interface{}{
				map[string]interface{}{"k1": "v1"},
			},
			XYWH:   [4]int{1, 2, 3, 4},
			Layout: "layout",
		}},
	}

	err := store.CreateReport(ctx, s, r)
	if err != nil {
		t.Fatal(err)
	}

	return r
}

func sTestAPIGatewayRoute(ctx context.Context, t *testing.T, s store.Storer, r string) *types.ApigwRoute {
	gwr := &types.ApigwRoute{
		ID:       su.NextID(),
		Endpoint: "/testing/" + r,
		Method:   "POST",
		Enabled:  true,
		Group:    0,
		Meta: types.ApigwRouteMeta{
			Debug: true,
			Async: true,
		},
	}

	err := store.CreateApigwRoute(ctx, s, gwr)
	if err != nil {
		t.Fatal(err)
	}

	return gwr
}

func sTestAPIGatewayFilter(ctx context.Context, t *testing.T, s store.Storer, routeID uint64, pfx string) *types.ApigwFilter {
	gwf := &types.ApigwFilter{
		ID:     su.NextID(),
		Route:  routeID,
		Weight: 0,
		Ref:    pfx + "_ref",
		Kind:   pfx + "_kind",
		Params: map[string]interface{}{
			"param1": "value1",
		},
	}

	err := store.CreateApigwFilter(ctx, s, gwf)
	if err != nil {
		t.Fatal(err)
	}

	return gwf
}
