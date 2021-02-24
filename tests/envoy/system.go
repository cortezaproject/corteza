package envoy

import (
	"context"
	"testing"

	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
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
			Order:  0,
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
			Value:     []byte("{\"k\": \"vs1\"}"),
			UpdatedAt: updatedAt,
			UpdatedBy: usrID,
		},
		{
			Name:      pfx + "_setting_2",
			Value:     []byte("{\"k\": \"vs2\"}"),
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
			Resource:  "compose",
			Operation: "read",
			Access:    rbac.Allow,
		},
		{
			RoleID:    roleID,
			Resource:  "system",
			Operation: "read",
			Access:    rbac.Deny,
		},
		{
			RoleID:    roleID,
			Resource:  "system:role:*",
			Operation: "read",
			Access:    rbac.Deny,
		},
		{
			RoleID:    roleID,
			Resource:  types.RoleRBACResource.AppendID(roleID),
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
