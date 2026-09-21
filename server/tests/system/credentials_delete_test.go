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
)

func (h helper) repoMakeCredentials(ownerID uint64) *types.Credential {
	c := &types.Credential{
		ID:          id.Next(),
		OwnerID:     ownerID,
		Kind:        "password",
		Credentials: "hash",
		CreatedAt:   time.Now(),
	}

	h.noError(store.CreateCredential(context.Background(), service.DefaultStore, c))
	return c
}

func (h helper) lookupCredentials(ID uint64) *types.Credential {
	c, err := store.LookupCredentialByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return c
}

// Users can manage their own credentials but must not be
// able to use that to remove credentials of someone else
func TestUserCredentialsDeleteForeign(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.storeCurrentUser()

	victim := h.createUserWithEmail(h.randEmail())
	c := h.repoMakeCredentials(victim.ID)

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d/credentials/%d", h.cUser.ID, c.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("credentials.errors.notFound")).
		End()

	h.a.Nil(h.lookupCredentials(c.ID).DeletedAt, "credentials of another user must not be removed")
}

func TestUserCredentialsDeleteOwn(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.storeCurrentUser()

	c := h.repoMakeCredentials(h.cUser.ID)

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d/credentials/%d", h.cUser.ID, c.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupCredentials(c.ID).DeletedAt)
}

func TestUserCredentialsDeleteForeignForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.storeCurrentUser()

	victim := h.createUserWithEmail(h.randEmail())
	c := h.repoMakeCredentials(victim.ID)

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d/credentials/%d", victim.ID, c.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("credentials.errors.notAllowedToManage")).
		End()

	h.a.Nil(h.lookupCredentials(c.ID).DeletedAt)
}
