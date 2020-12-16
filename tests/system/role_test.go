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
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func (h helper) clearRoles() {
	h.noError(store.TruncateRoles(context.Background(), service.DefaultStore))
}

func (h helper) clearRoleMembers() {
	h.noError(store.TruncateRoleMembers(context.Background(), service.DefaultStore))
}

func (h helper) createRole(res *types.Role) *types.Role {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateRole(context.Background(), res))
	return res
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

	h.deny(types.RoleRBACResource.AppendID(f.ID), "read")

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
		Header("Accept", "application/json").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create roles")).
		End()
}

func TestRoleCreateNotUnique(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemRBACResource, "role.create")

	role := h.repoMakeRole()
	h.apiInit().
		Post("/roles/").
		Header("Accept", "application/json").
		FormData("name", rs()).
		FormData("handle", role.Handle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("role handle not unique")).
		End()

	h.apiInit().
		Post("/roles/").
		Header("Accept", "application/json").
		FormData("name", role.Name).
		FormData("handle", "handle_"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("role name not unique")).
		End()

}

func TestRoleCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemRBACResource, "role.create")

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
		Header("Accept", "application/json").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this role")).
		End()
}

func TestRoleUpdate(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeRole()
	h.allow(types.RoleRBACResource.AppendWildcard(), "update")

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
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this role")).
		End()
}

func TestRoleDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.RoleRBACResource.AppendWildcard(), "delete")

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

func TestRoleLabels(t *testing.T) {
	h := newHelper(t)
	h.clearRoles()

	h.allow(types.SystemRBACResource, "role.create")
	h.allow(types.RoleRBACResource.AppendWildcard(), "read")
	h.allow(types.RoleRBACResource.AppendWildcard(), "update")
	h.allow(types.RoleRBACResource.AppendWildcard(), "delete")

	var (
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Role{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/roles/",
			types.Role{Labels: map[string]string{"foo": "bar", "bar": "42"}},
			payload,
		)
		req.NotZero(payload.ID)

		h.a.Equal(payload.Labels["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Labels["bar"], "42",
			"labels must contain bar with value 42")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")

		ID = payload.ID
	})

	t.Run("update", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req     = require.New(t)
			payload = &types.Role{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("PUT /roles/%d", ID),
			&types.Role{Labels: map[string]string{"foo": "baz", "baz": "123"}},
			payload,
		)
		req.NotZero(payload.ID)
		//req.Nil(payload.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Labels["foo"], "baz",
			"labels must contain foo with value baz")
		req.NotContains(payload.Labels, "bar",
			"labels must not contain bar")
		req.Equal(payload.Labels["baz"], "123",
			"labels must contain baz with value 123")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)
			set = types.RoleSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/roles/", &set, url.Values{"labels": []string{"baz=123"}})
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
