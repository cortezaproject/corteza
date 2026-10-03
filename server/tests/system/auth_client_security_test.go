package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// Auth clients can impersonate users and force roles on everyone that uses them,
// both need the same permissions as doing that directly

func TestAuthClientCreateImpersonateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-client.create")

	h.apiInit().
		Post("/auth/clients/").
		Header("Accept", "application/json").
		JSON(helpers.JSON(&types.AuthClient{
			Handle:     "handle_" + rs(),
			Meta:       &types.AuthClientMeta{Name: rs()},
			Scope:      "profile api",
			ValidGrant: "client_credentials",
			Security:   &types.AuthClientSecurity{ImpersonateUser: id.Next()},
		})).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToImpersonate")).
		End()
}

func TestAuthClientCreateImpersonate(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-client.create")
	helpers.AllowMe(h, types.UserRbacResource(0), "impersonate")

	h.apiInit().
		Post("/auth/clients/").
		Header("Accept", "application/json").
		JSON(helpers.JSON(&types.AuthClient{
			Handle:     "handle_" + rs(),
			Meta:       &types.AuthClientMeta{Name: rs()},
			Scope:      "profile api",
			ValidGrant: "client_credentials",
			Security:   &types.AuthClientSecurity{ImpersonateUser: id.Next()},
		})).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAuthClientCreateImpersonateMyself(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-client.create")

	h.apiInit().
		Post("/auth/clients/").
		Header("Accept", "application/json").
		JSON(helpers.JSON(&types.AuthClient{
			Handle:     "handle_" + rs(),
			Meta:       &types.AuthClientMeta{Name: rs()},
			Scope:      "profile api",
			ValidGrant: "client_credentials",
			Security:   &types.AuthClientSecurity{ImpersonateUser: h.cUser.ID},
		})).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAuthClientCreateForcedRolesForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	helpers.AllowMe(h, types.ComponentRbacResource(), "auth-client.create")

	h.apiInit().
		Post("/auth/clients/").
		Header("Accept", "application/json").
		JSON(helpers.JSON(&types.AuthClient{
			Handle:     "handle_" + rs(),
			Meta:       &types.AuthClientMeta{Name: rs()},
			Scope:      "profile api",
			ValidGrant: "authorization_code",
			Security:   &types.AuthClientSecurity{ForcedRoles: []string{fmt.Sprintf("%d", id.Next())}},
		})).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToForceRoles")).
		End()
}

func TestAuthClientUpdateForcedRolesForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeAuthClient()
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "update")

	client.Meta = &types.AuthClientMeta{Name: rs()}
	client.Security = &types.AuthClientSecurity{ForcedRoles: []string{fmt.Sprintf("%d", id.Next())}}

	h.apiInit().
		Put(fmt.Sprintf("/auth/clients/%d", client.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(client)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToForceRoles")).
		End()

	stored := h.lookupAuthClientByID(client.ID)
	h.a.True(stored.Security == nil || len(stored.Security.ForcedRoles) == 0)
}

func TestAuthClientUpdateForcedRoles(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeAuthClient()
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "update")
	helpers.AllowMe(h, types.RoleRbacResource(0), "members.manage")

	client.Meta = &types.AuthClientMeta{Name: rs()}
	client.Security = &types.AuthClientSecurity{ForcedRoles: []string{fmt.Sprintf("%d", id.Next())}}

	h.apiInit().
		Put(fmt.Sprintf("/auth/clients/%d", client.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(client)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

// Secret is as good as the client itself, reading the client is not enough to get or replace it

func TestAuthClientExposeSecretReadOnly(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeAuthClient()
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/auth/clients/%d/secret", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToUpdate")).
		End()
}

func TestAuthClientRegenerateSecretReadOnly(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeAuthClient()
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read")

	secret := h.lookupAuthClientByID(client.ID).Secret

	h.apiInit().
		Post(fmt.Sprintf("/auth/clients/%d/secret", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToUpdate")).
		End()

	h.a.Equal(secret, h.lookupAuthClientByID(client.ID).Secret)
}

func TestAuthClientExposeSecret(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeAuthClient()
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read")
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "update")

	h.apiInit().
		Get(fmt.Sprintf("/auth/clients/%d/secret", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

// A client that impersonates another user is as good as that user: its
// secret and its configuration are off limits without impersonate rights
func (h helper) repoMakeImpersonatingAuthClient(userID uint64) *types.AuthClient {
	client := h.repoMakeAuthClient()
	client.ValidGrant = "client_credentials"
	client.Security = &types.AuthClientSecurity{ImpersonateUser: userID}
	h.noError(store.UpdateAuthClient(context.Background(), service.DefaultStore, client))
	return client
}

func TestAuthClientExposeSecretImpersonatingForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeImpersonatingAuthClient(id.Next())
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read", "update")

	h.apiInit().
		Get(fmt.Sprintf("/auth/clients/%d/secret", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToImpersonate")).
		End()
}

func TestAuthClientExposeSecretImpersonating(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeImpersonatingAuthClient(id.Next())
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read", "update")
	helpers.AllowMe(h, types.UserRbacResource(0), "impersonate")

	h.apiInit().
		Get(fmt.Sprintf("/auth/clients/%d/secret", client.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestAuthClientUpdateImpersonatingForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAuthClients()
	client := h.repoMakeImpersonatingAuthClient(id.Next())
	helpers.AllowMe(h, types.AuthClientRbacResource(0), "read", "update")

	h.apiInit().
		Put(fmt.Sprintf("/auth/clients/%d", client.ID)).
		JSON(fmt.Sprintf(`{"handle": "%s", "validGrant": "client_credentials", "security": {"impersonateUser": "%d"}, "meta": {"name": "renamed"}}`, client.Handle, client.Security.ImpersonateUser)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("auth-client.errors.notAllowedToImpersonate")).
		End()
}
