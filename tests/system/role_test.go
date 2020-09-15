package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearRoles() {
	h.noError(store.TruncateRoles(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeRole(ss ...string) *types.Role {
	var r = &types.Role{
		ID:        id.Next(),
		CreatedAt: time.Now(),
	}

	if len(ss) > 1 {
		r.Handle = ss[1]
	} else {
		r.Handle = "h_" + rs()

	}
	if len(ss) > 0 {
		r.Name = ss[0]
	} else {
		r.Name = "n_" + rs()
	}

	h.a.NoError(store.CreateRole(context.Background(), service.DefaultStore, r))

	return r
}

func (h helper) lookupRoleByID(ID uint64) *types.Role {
	res, err := store.LookupRoleByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestRoleRead(t *testing.T) {
	h := newHelper(t)

	u := h.repoMakeRole()

	h.apiInit().
		Get(fmt.Sprintf("/roles/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, u.Name)).
		Assert(jsonpath.Equal(`$.response.roleID`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestRoleList(t *testing.T) {
	h := newHelper(t)

	h.repoMakeRole(h.randEmail())
	h.repoMakeRole(h.randEmail())

	h.apiInit().
		Get("/roles/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestRoleList_filterForbidden(t *testing.T) {
	h := newHelper(t)

	// @todo this can be a problematic test because it leaves
	//       behind roles that are not denied this context
	//       db purge might be needed

	h.repoMakeRole("role")
	f := h.repoMakeRole()

	h.deny(types.RolePermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/roles/").
		Query("handle", f.Handle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.handle=="%s"]`, f.Handle))).
		End()
}

func TestRoleCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/roles/").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create roles")).
		End()
}

func TestRoleCreateNotUnique(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "role.create")

	role := h.repoMakeRole()
	h.apiInit().
		Post("/roles/").
		FormData("name", rs()).
		FormData("handle", role.Handle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("role handle not unique")).
		End()

	h.apiInit().
		Post("/roles/").
		FormData("name", role.Name).
		FormData("handle", "handle_"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("role name not unique")).
		End()

}

func TestRoleCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "role.create")

	h.apiInit().
		Post("/roles/").
		FormData("name", rs()).
		FormData("handle", "handle_"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestRoleUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeRole()

	h.apiInit().
		Put(fmt.Sprintf("/roles/%d", u.ID)).
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this role")).
		End()
}

func TestRoleUpdate(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeRole()
	h.allow(types.RolePermissionResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/roles/%d", res.ID)).
		FormData("name", newName).
		FormData("handle", newHandle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupRoleByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(newName, res.Name)
	h.a.Equal(newHandle, res.Handle)
}

func TestRoleDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeRole()

	h.apiInit().
		Delete(fmt.Sprintf("/roles/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this role")).
		End()
}

func TestRoleDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.RolePermissionResource.AppendWildcard(), "delete")

	res := h.repoMakeRole()

	h.apiInit().
		Delete(fmt.Sprintf("/roles/%d", res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupRoleByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}
