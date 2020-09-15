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

func (h helper) clearApplications() {
	h.noError(store.TruncateApplications(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeApplication(ss ...string) *types.Application {
	var res = &types.Application{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Unify:     &types.ApplicationUnify{},
	}

	if len(ss) > 0 {
		res.Name = ss[0]
	} else {
		res.Name = "n_" + rs()
	}

	h.a.NoError(store.CreateApplication(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) lookupApplicationByID(ID uint64) *types.Application {
	res, err := store.LookupApplicationByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestApplicationRead(t *testing.T) {
	h := newHelper(t)

	u := h.repoMakeApplication()

	h.apiInit().
		Get(fmt.Sprintf("/application/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, u.Name)).
		Assert(jsonpath.Equal(`$.response.applicationID`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestApplicationList(t *testing.T) {
	h := newHelper(t)

	h.repoMakeApplication(h.randEmail())
	h.repoMakeApplication(h.randEmail())

	h.apiInit().
		Get("/application/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApplicationList_filterForbidden(t *testing.T) {
	h := newHelper(t)

	// @todo this can be a problematic test because it leaves
	//       behind applications that are not denied this context
	//       db purge might be needed

	h.repoMakeApplication("application")
	f := h.repoMakeApplication()

	h.deny(types.ApplicationPermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/application/").
		Query("name", f.Name).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.name=="%s"]`, f.Name))).
		End()
}

func TestApplicationCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/application/").
		FormData("name", rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create applications")).
		End()
}

func TestApplicationCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "application.create")

	h.apiInit().
		Post("/application/").
		FormData("name", rs()).
		FormData("handle", "handle_"+rs()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApplicationUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeApplication()

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", u.ID)).
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this application")).
		End()
}

func TestApplicationUpdate(t *testing.T) {
	h := newHelper(t)
	res := h.repoMakeApplication()
	h.allow(types.ApplicationPermissionResource.AppendWildcard(), "update")

	newName := "updated-" + rs()
	newHandle := "updated-" + rs()

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", res.ID)).
		FormData("name", newName).
		FormData("handle", newHandle).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupApplicationByID(res.ID)
	h.a.NotNil(res)
	h.a.Equal(newName, res.Name)
}

func TestApplicationDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeApplication()

	h.apiInit().
		Delete(fmt.Sprintf("/application/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this application")).
		End()
}

func TestApplicationDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ApplicationPermissionResource.AppendWildcard(), "delete")

	res := h.repoMakeApplication()

	h.apiInit().
		Delete(fmt.Sprintf("/application/%d", res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupApplicationByID(res.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}
