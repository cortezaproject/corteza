package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
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

	if res.Config.DAL == nil {
		res.Config.DAL = &types.DalConnectionConfigDAL{
			Type: "corteza::dal:connection:dsn",
			Params: map[string]any{
				"dsn": "sqlite3://file::memory:?cache=shared&mode=memory",
			},
		}
	}

	if res.Config.DAL.ModelIdent == "" {
		res.Config.DAL.ModelIdent = "compose_records_{{namespace}}_{{module}}"
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
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToSearch")).
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
	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "dal-config.manage")

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

func Test_dal_connection_create_invalid_type(t *testing.T) {
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
		Assert(helpers.AssertErrorP("corteza::system:primary-dal-connection")).
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
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToCreate")).
		End()
}

func Test_dal_connection_update(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	sl := h.createDalConnection(&types.DalConnection{
		Handle: "test_connection",
	})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update", "dal-config.manage")

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

func Test_dal_connection_create_dal_config_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connection.create")

	h.apiInit().
		Post("/dal/connections/").
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(`{"handle":"with_dal","type":"corteza::system:dal-connection","meta":{"name":"With DAL"},"config":{"dal":{"type":"corteza::dal:connection:dsn","params":{"dsn":"sqlite3://file::memory:?cache=shared&mode=memory"}}}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToCreate")).
		End()

	h.apiInit().
		Post("/dal/connections/").
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(`{"handle":"without_dal","type":"corteza::system:dal-connection","meta":{"name":"Without DAL"},"config":{}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.connectionID")).
		End()
}

func Test_dal_connection_update_keeps_meta_and_config(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	c := h.createDalConnection(&types.DalConnection{Handle: "test_connection"})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update", "dal-config.manage")

	const dsn = "sqlite3://file::memory:?cache=shared&mode=memory&_edited=1"

	h.apiInit().
		Put(fmt.Sprintf("/dal/connections/%d", c.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{
			"handle": "test_connection",
			"type": "corteza::system:primary-dal-connection",
			"meta": {"name": "Edited", "ownership": "Ops", "location": {"properties": {"name": "Ljubljana"}}},
			"config": {"dal": {"type": "corteza::dal:connection:dsn", "params": {"dsn": %q}, "modelIdent": "edited_{{module}}"}}
		}`, dsn)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.meta.name", "Edited")).
		Assert(jsonpath.Equal("$.response.config.dal.params.dsn", dsn)).
		End()

	stored, err := store.LookupDalConnectionByID(context.Background(), service.DefaultStore, c.ID)
	h.noError(err)
	h.a.Equal("Edited", stored.Meta.Name)
	h.a.Equal("Ops", stored.Meta.Ownership)
	h.a.Equal("Ljubljana", stored.Meta.Location.Properties.Name)
	h.a.Equal(dsn, stored.Config.DAL.Params["dsn"])
	h.a.Equal("edited_{{module}}", stored.Config.DAL.ModelIdent)
	h.a.Equal(types.DalConnectionResourceType, stored.Type, "type is fixed at create")
	h.a.Equal(c.CreatedBy, stored.CreatedBy)
}

func Test_dal_connection_update_hides_dal_config(t *testing.T) {
	h := newHelper(t)
	defer h.clearDalConnections()

	c := h.createDalConnection(&types.DalConnection{Handle: "test_connection"})

	helpers.AllowMe(h, types.DalConnectionRbacResource(0), "update")

	h.apiInit().
		Put(fmt.Sprintf("/dal/connections/%d", c.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(`{"handle": "test_connection", "type": "corteza::system:dal-connection", "meta": {"name": "Renamed"}, "config": {}}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.response.meta.name", "Renamed")).
		Assert(jsonpath.NotPresent("$.response.config.dal")).
		End()

	stored, err := store.LookupDalConnectionByID(context.Background(), service.DefaultStore, c.ID)
	h.noError(err)
	h.a.Equal(c.Config.DAL.Params["dsn"], stored.Config.DAL.Params["dsn"])
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
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToUpdate")).
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

func Test_dal_connection_read_forbidden(t *testing.T) {
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
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToRead")).
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
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToDelete")).
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
		Assert(helpers.AssertError("dal-connection.errors.notAllowedToUndelete")).
		End()
}
