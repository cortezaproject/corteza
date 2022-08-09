package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearDalConnections() {
	cc, _, err := store.SearchDalConnections(context.Background(), service.DefaultStore, types.DalConnectionFilter{})
	h.noError(err)

	for _, c := range cc {
		if c.Type == types.DalPrimaryConnectionResourceType {
			continue
		}
		h.noError(store.DeleteDalConnection(context.Background(), service.DefaultStore, c))
	}
}

func (h helper) getPrimaryConnection() *types.DalConnection {
	cc, _, err := store.SearchDalConnections(context.Background(), service.DefaultStore, types.DalConnectionFilter{Type: types.DalPrimaryConnectionResourceType})
	h.a.NoError(err)

	if len(cc) != 1 {
		h.a.FailNow("invalid state: no or too many primary connections")
	}

	return cc[0]
}

func (h helper) createDalConnection(res *types.DalConnection) *types.DalConnection {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.Meta.Name == "" {
		res.Meta.Name = "Test Connection"
	}
	if res.Handle == "" {
		res.Handle = "test_connection"
	}
	if res.Type == "" {
		res.Type = types.DalConnectionResourceType
	}
	if res.Meta.Ownership == "" {
		res.Meta.Ownership = "tester"
	}

	if res.Config.DAL.ModelIdent == "" {
		res.Config.DAL.ModelIdent = "compose_records_{{namespace}}_{{module}}"
	}
	if res.Config.DAL.AttributeIdent == "" {
		res.Config.DAL.AttributeIdent = "values"
	}
	if res.Config.Connection.Params == nil {
		res.Config.Connection = dal.NewDSNConnection("sqlite3://file::memory:?cache=shared&mode=memory")
	}

	if len(res.Config.DAL.Operations) == 0 {
		res.Config.DAL.Operations = dal.FullOperations()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	if res.CreatedBy == 0 {
		res.CreatedBy = h.cUser.ID
	}

	h.a.NoError(service.DefaultStore.CreateDalConnection(context.Background(), res))
	h.a.NoError(service.DefaultDalConnection.ReloadConnections(context.Background()))
	return res
}

func Test_dal_connection_list(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connections.search")
	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "read")

	h.apiInit().
		Get("/dal/connections/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func Test_dal_connection_list_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	h.apiInit().
		Get("/dal/connections/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToSearch")).
		End()
}

func Test_dal_connection_list_forbidden_read(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connections.search")

	h.apiInit().
		Get("/dal/connections/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.response.set", 0)).
		End()
}

func Test_dal_connection_create(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connection.create")

	h.apiInit().
		Post("/dal/connections/").
		Body(loadScenarioRequest(t, "generic.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		End()
}

func Test_dal_connection_create_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	h.apiInit().
		Post("/dal/connections/").
		Body(loadScenarioRequest(t, "generic.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToCreate")).
		End()
}

func Test_dal_connection_update(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update")

	h.apiInit().
		Put(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadScenarioRequest(t, "generic.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.handle", "test_connection_edited")).
		End()
}

func Test_dal_connection_update_primary(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.getPrimaryConnection()

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update")

	// a bit of a problem with testing primary connection update
	//
	// when using (for running tests) anything else than connection params specified
	// in the generic.json scenario, the update will fail
	// with "can not update connection parameters for primary ..."
	//
	// see Update on dalConnection service.

	h.apiInit().
		Put(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadScenarioRequest(t, "generic.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.meta.name", "Primary Connection EDITED")).
		End()
}

func Test_dal_connection_update_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	h.apiInit().
		Put(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadScenarioRequest(t, "generic.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToUpdate")).
		End()
}

func Test_dal_connection_read(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		End()
}

func Test_dal_connection_read_forbiden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	h.apiInit().
		Get(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToRead")).
		End()
}

func Test_dal_connection_delete(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_connection_delete_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	h.apiInit().
		Delete(fmt.Sprintf("/dal/connections/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToDelete")).
		End()
}

func Test_dal_connection_undelete(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle:    "test_connection",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "delete")

	h.apiInit().
		Post(fmt.Sprintf("/dal/connections/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_connection_undelete_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle:    "test_connection",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	h.apiInit().
		Post(fmt.Sprintf("/dal/connections/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalConnection.errors.notAllowedToUndelete")).
		End()
}
