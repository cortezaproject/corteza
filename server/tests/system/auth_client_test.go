package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearAuthClients() {
	h.noError(store.TruncateAuthClients(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeAuthClient(ss ...string) *types.AuthClient {
	res := &types.AuthClient{
		ID:        id.Next(),
		CreatedAt: time.Now(),
	}

	if len(ss) > 0 {
		res.Handle = ss[0]
	} else {
		res.Handle = "n_" + rs()
	}

	h.a.NoError(store.CreateAuthClient(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) lookupAuthClientByID(id uint64) *types.AuthClient {
	res, err := store.LookupAuthClientByID(context.Background(), service.DefaultStore, id)
	h.noError(err)
	return res
}

func (h helper) lookupAuthClientByHandle(handle string) *types.AuthClient {
	res, err := store.LookupAuthClientByHandle(context.Background(), service.DefaultStore, handle)
	h.noError(err)
	return res
}

func TestAuthClientList(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()

	h.repoMakeAuthClient()
	h.repoMakeAuthClient()

	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-clients.search")
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read")

	h.apiInit().
		Get("/auth/clients/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestAuthClientRead(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()

	client := h.repoMakeAuthClient()

	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/auth/clients/%d", client.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAuthClientCreate(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()

	handle := rs()

	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-clients.search")
	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-client.create")

	h.apiInit().
		Post("/auth/clients/").
		Header("Accept", "application/json").
		FormData("handle", handle).
		FormData("scope", "profile api").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res := h.lookupAuthClientByHandle(handle)
	h.a.NotNil(res)
	h.a.Equal(handle, res.Handle)
}

func TestAuthClientUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()

	client := h.repoMakeAuthClient()
	client.Handle = rs()

	helpers.AllowMe(h, types.AuthClientRbacResource(0), "update")

	h.apiInit().
		Put(fmt.Sprintf("/auth/clients/%d", client.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(client)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res := h.lookupAuthClientByHandle(client.Handle)
	h.a.NotNil(res)
}

func TestAuthClientDelete(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()

	client := h.repoMakeAuthClient()

	helpers.AllowMe(h, types.AuthClientRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/auth/clients/%d", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res := h.lookupAuthClientByID(client.ID)
	h.a.NotNil(res)
	h.a.NotNil(res.DeletedAt)
}

func TestAuthClientUnDelete(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()

	client := h.repoMakeAuthClient()

	helpers.AllowMe(h, types.AuthClientRbacResource(0), "delete")

	h.apiInit().
		Post(fmt.Sprintf("/auth/clients/%d/undelete", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res := h.lookupAuthClientByID(client.ID)
	h.a.NotNil(res)
	h.a.Nil(res.DeletedAt)
}
