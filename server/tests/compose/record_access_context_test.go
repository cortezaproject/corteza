package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func crissCrossUserRoles(ctx context.Context, s store.Storer, h helper, uu systemTypes.UserSet, rr systemTypes.RoleSet) (map[string]*systemTypes.User, map[string]*systemTypes.Role) {
	ux := make(map[string]*systemTypes.User)
	for _, r := range uu {
		ux[r.Name] = r
	}

	rx := make(map[string]*systemTypes.Role)
	for _, r := range rr {
		rx[r.Handle] = r
	}

	for _, u := range uu {
		if _, ok := rx[u.Name]; !ok {
			h.a.FailNow(fmt.Sprintf("corresponding role not found for user: r: %s; u: %s", u.Name, u.Handle))
		}
		rID := rx[u.Name].ID
		h.a.NoError(store.CreateRoleMember(ctx, s, &systemTypes.RoleMember{
			RoleID: rID,
			UserID: u.ID,
		}))
		u.SetRoles(rID)
	}

	return ux, rx
}

func Test_record_access_context(t *testing.T) {
	ctx, h, s := setup(t)
	loadScenario(ctx, defStore, t, h)

	// setup
	rr, _, err := store.SearchRoles(ctx, s, systemTypes.RoleFilter{})
	h.a.NoError(err)
	uu, _, err := store.SearchUsers(ctx, s, systemTypes.UserFilter{})
	h.a.NoError(err)

	ux, rx := crissCrossUserRoles(ctx, s, h, uu, rr)

	ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
	h.a.NoError(err)

	mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "mod1")
	h.a.NoError(err)

	records, _, err := store.SearchComposeRecords(ctx, s, mod, types.RecordFilter{})
	h.a.NoError(err)
	rec := records[0]

	for _, r := range rr {
		helpers.Allow(r, types.NamespaceRbacResource(0), "read")
		helpers.Allow(r, types.ModuleRbacResource(0, 0), "read")
	}

	helpers.DenyMe(h, rec.RbacResource(), "read")
	helpers.Allow(rx["owner"], rec.RbacResource(), "read")
	helpers.Allow(rx["creator"], rec.RbacResource(), "read")
	helpers.Allow(rx["updater"], rec.RbacResource(), "read")
	helpers.Allow(rx["deleter"], rec.RbacResource(), "read")

	h.a.NoError(service.UpdateRbacRoles(ctx, testApp.Log, rbac.Global(), nil, nil, nil))
	rbac.Global().Reload(ctx)

	t.Run("generic user with no ctx role", func(t *testing.T) {
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertError("record.errors.notAllowedToRead")).
			End()
	})

	t.Run("user with owner ctx role", func(t *testing.T) {
		h.identityToHelper(ux["owner"])
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})

	t.Run("user with creator ctx role", func(t *testing.T) {
		h.identityToHelper(ux["creator"])
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})

	t.Run("user with updater ctx role", func(t *testing.T) {
		h.identityToHelper(ux["updater"])
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})

	t.Run("user with deleter ctx role", func(t *testing.T) {
		h.identityToHelper(ux["deleter"])
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})

	t.Run("user with creator ctx role", func(t *testing.T) {
		h.identityToHelper(ux["creator"])
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})

	t.Run("user with creator ctx role", func(t *testing.T) {
		h.identityToHelper(ux["creator"])
		h.apiInit().
			Get(fmt.Sprintf("/namespace/%d/module/%d/record/%d", mod.NamespaceID, mod.ID, rec.ID)).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	})
}
