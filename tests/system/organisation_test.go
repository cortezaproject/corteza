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

func (h helper) repoOrganisation() repository.OrganisationRepository {
	return repository.Organisation(context.Background(), db())
}

func (h helper) repoMakeOrganisation(name string) *types.Organisation {
	a, err := h.
		repoOrganisation().
		Create(&types.Organisation{Name: name})
	h.a.NoError(err)

	return a
}

func TestOrganisationRead(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)

	a := h.repoMakeOrganisation("one-app")

	h.apiInit().
		Get(fmt.Sprintf("/organisation/%d", a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, a.Name)).
		Assert(jsonpath.Equal(`$.response.organisationID`, fmt.Sprintf("%d", a.ID))).
		End()
}

func TestOrganisationList(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)

	h.repoMakeOrganisation("app")
	h.repoMakeOrganisation("app")

	h.apiInit().
		Get("/organisation/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestOrganisationCreateForbidden(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)

	h.apiInit().
		Post("/organisation/").
		FormData("name", "my-app").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("Not allowed to create organisation")).
		End()
}

func TestOrganisationCreate(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "organisation.create")

	h.apiInit().
		Post("/organisation/").
		FormData("name", "my-app").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestOrganisationUpdateForbidden(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)
	a := h.repoMakeOrganisation("one-app")

	h.apiInit().
		Put(fmt.Sprintf("/organisation/%d", a.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("Not allowed to update organisation")).
		End()
}

func TestOrganisationUpdate(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)
	a := h.repoMakeOrganisation("one-app")
	h.allow(types.OrganisationPermissionResource.AppendWildcard(), "update")

	h.apiInit().
		Put(fmt.Sprintf("/organisation/%d", a.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	a, err := h.repoOrganisation().FindByID(a.ID)
	h.a.NoError(err)
	h.a.NotNil(a)
	h.a.Equal(a.Name, "changed-name")
}

func TestOrganisationDeleteForbidden(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)
	a := h.repoMakeOrganisation("one-app")

	h.apiInit().
		Delete(fmt.Sprintf("/organisation/%d", a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("Not allowed to delete organisation")).
		End()
}

func TestOrganisationDelete(t *testing.T) {
	t.Skip("pending implementation")

	h := newHelper(t)
	h.allow(types.OrganisationPermissionResource.AppendWildcard(), "delete")

	a := h.repoMakeOrganisation("one-app")

	h.apiInit().
		Delete(fmt.Sprintf("/organisation/%d", a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	a, err := h.repoOrganisation().FindByID(a.ID)
	h.a.Error(err, "system.repository.OrganisationNotFound")
}
