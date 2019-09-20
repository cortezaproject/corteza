package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) repoRole() repository.RoleRepository {
	return repository.Role(context.Background(), db())
}

func (h helper) repoMakeRole(ss ...string) *types.Role {
	var r = &types.Role{}
	if len(ss) > 1 {
		r.Handle = ss[1]
	} else {
		r.Handle = rs()

	}
	if len(ss) > 0 {
		r.Name = ss[0]
	} else {
		r.Name = rs()
	}

	r, err := h.
		repoRole().
		Create(r)
	h.a.NoError(err)

	return r
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

func TestRoleCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/roles/").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoCreatePermissions")).
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
		Assert(helpers.AssertError("system.service.RoleHandleNotUnique")).
		End()

	h.apiInit().
		Post("/roles/").
		FormData("name", role.Name).
		FormData("handle", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.RoleNameNotUnique")).
		End()

}

func TestRoleCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "role.create")

	h.apiInit().
		Post("/roles/").
		FormData("name", rs()).
		FormData("handle", rs()).
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
		Assert(helpers.AssertError("system.service.NoUpdatePermissions")).
		End()
}

func TestRoleUpdate(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeRole()
	h.allow(types.RolePermissionResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/roles/%d", u.ID)).
		FormData("name", newName).
		FormData("handle", newHandle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := h.repoRole().FindByID(u.ID)
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.Equal(newName, u.Name)
	h.a.Equal(newHandle, u.Handle)
}

func TestRoleDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeRole()

	h.apiInit().
		Delete(fmt.Sprintf("/roles/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoPermissions")).
		End()
}

func TestRoleDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.RolePermissionResource.AppendWildcard(), "delete")

	u := h.repoMakeRole()

	h.apiInit().
		Delete(fmt.Sprintf("/roles/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := h.repoRole().FindByID(u.ID)
	h.a.Error(err, "system.repository.RoleNotFound")
}
