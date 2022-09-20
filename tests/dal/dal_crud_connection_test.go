package dal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func Test_dal_crud_connection_list(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMeDalConnectionSearch(h)

	h.apiInit().
		Get("/system/dal/connections/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		// This will be the primary one
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func Test_dal_crud_connection_list_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	h.apiInit().
		Get("/system/dal/connections/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToSearch")).
		End()
}

func Test_dal_crud_connection_list_forbidden_read(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connections.search")

	h.apiInit().
		Get("/system/dal/connections/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.response.set", 0)).
		End()
}

func Test_dal_crud_connection_create(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connection.create")

	h.apiInit().
		Post("/system/dal/connections/").
		Body(loadRequestFromGenerics(t, "ok_connection.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		End()
}

func Test_dal_crud_connection_create_invalid_type(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connection.create")

	h.apiInit().
		Post("/system/dal/connections/").
		Body(loadRequestFromGenerics(t, "nok_connection_invalid_type.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertErrorP("corteza::system:primary-dal-connection")).
		End()
}

func Test_dal_crud_connection_create_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	h.apiInit().
		Post("/system/dal/connections/").
		Body(loadRequestFromGenerics(t, "ok_connection.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToCreate")).
		End()
}

func Test_dal_crud_connection_update(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update")
	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "dal-config.manage")

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadRequestFromGenerics(t, "ok_connection_update.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.handle", "test_connection_edited")).
		End()
}

func Test_dal_crud_connection_update_primary(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.getPrimaryConnection()

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update")

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadRequestFromScenario(t, "connection.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.meta.name", "Primary Connection EDITED")).
		End()
}

func Test_dal_crud_connection_update_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadRequestFromGenerics(t, "ok_connection_update.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToUpdate")).
		End()
}

func Test_dal_crud_connection_read(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		End()
}

func Test_dal_crud_connection_read_forbiden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	h.apiInit().
		Get(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToRead")).
		End()
}

func Test_dal_crud_connection_delete(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_crud_connection_delete_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	h.apiInit().
		Delete(fmt.Sprintf("/system/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToDelete")).
		End()
}

func Test_dal_crud_connection_undelete(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle:    "test_connection",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "delete")

	h.apiInit().
		Post(fmt.Sprintf("/system/dal/connections/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_crud_connection_undelete_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createDalConnection(&types.DalConnection{
		Handle:    "test_connection",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	h.apiInit().
		Post(fmt.Sprintf("/system/dal/connections/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToUndelete")).
		End()
}
