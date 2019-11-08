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

func (h helper) repoApplication() repository.ApplicationRepository {
	return repository.Application(context.Background(), db())
}

func (h helper) repoMakeApplication(name string) *types.Application {
	a, err := h.
		repoApplication().
		Create(&types.Application{Name: name})
	h.a.NoError(err)

	return a
}

func TestApplicationRead(t *testing.T) {
	h := newHelper(t)

	a := h.repoMakeApplication("one-app")

	h.apiInit().
		Get(fmt.Sprintf("/application/%d", a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, a.Name)).
		Assert(jsonpath.Equal(`$.response.applicationID`, fmt.Sprintf("%d", a.ID))).
		End()
}

func TestApplicationList(t *testing.T) {
	h := newHelper(t)

	h.repoMakeApplication("app")
	h.repoMakeApplication("app")

	h.apiInit().
		Get("/application/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApplicationList_filterForbiden(t *testing.T) {
	h := newHelper(t)

	h.repoMakeApplication("app")
	f := h.repoMakeApplication("app_forbiden")

	h.deny(types.ApplicationPermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/application/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="app_forbiden"]`)).
		End()
}

func TestApplicationCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/application/").
		FormData("name", "my-app").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoPermissions")).
		End()
}

func TestApplicationCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "application.create")

	h.apiInit().
		Post("/application/").
		FormData("name", "my-app").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestApplicationUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	a := h.repoMakeApplication("one-app")

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", a.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoPermissions")).
		End()
}

func TestApplicationUpdate(t *testing.T) {
	h := newHelper(t)
	a := h.repoMakeApplication("one-app")
	h.allow(types.ApplicationPermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Put(fmt.Sprintf("/application/%d", a.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	a, err := h.repoApplication().FindByID(a.ID)
	h.a.NoError(err)
	h.a.NotNil(a)
	h.a.Equal("changed-name", a.Name)
}

func TestApplicationDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	a := h.repoMakeApplication("one-app")

	h.apiInit().
		Delete(fmt.Sprintf("/application/%d", a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoPermissions")).
		End()
}

func TestApplicationDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.ApplicationPermissionResource.AppendWildcard(), "delete")

	a := h.repoMakeApplication("one-app")

	h.apiInit().
		Delete(fmt.Sprintf("/application/%d", a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	a, err := h.repoApplication().FindByID(a.ID)
	h.a.Error(err, "system.repository.ApplicationNotFound")
}
