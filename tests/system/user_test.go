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

func (h helper) randEmail() string {
	return fmt.Sprintf("%s@test.tld", rs())
}

func (h helper) repoUser() repository.UserRepository {
	return repository.User(context.Background(), db())
}

func (h helper) repoMakeUser(email string) *types.User {
	u, err := h.
		repoUser().
		Create(&types.User{Email: email})
	h.a.NoError(err)

	return u
}

func TestUserRead(t *testing.T) {
	h := newHelper(t)

	u := h.repoMakeUser(h.randEmail())

	h.apiInit().
		Get(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.email`, u.Email)).
		Assert(jsonpath.Equal(`$.response.userID`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestUserList(t *testing.T) {
	h := newHelper(t)

	h.repoMakeUser(h.randEmail())
	h.repoMakeUser(h.randEmail())

	h.apiInit().
		Get("/users/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestUserCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/users/").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoCreatePermissions")).
		End()
}

func TestUserCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "user.create")

	h.apiInit().
		Post("/users/").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestUserUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeUser(h.randEmail())

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", u.ID)).
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoUpdatePermissions")).
		End()
}

func TestUserUpdate(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeUser(h.randEmail())
	h.allow(types.UserPermissionResource.AppendWildcard(), "update")

	newEmail := "updated-" + u.Email

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", u.ID)).
		FormData("email", newEmail).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := h.repoUser().FindByID(u.ID)
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.Equal(newEmail, u.Email)
}

func TestUserDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeUser(h.randEmail())

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("system.service.NoPermissions")).
		End()
}

func TestUserDelete(t *testing.T) {
	h := newHelper(t)
	h.allow(types.UserPermissionResource.AppendWildcard(), "delete")

	u := h.repoMakeUser(h.randEmail())

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := h.repoUser().FindByID(u.ID)
	h.a.Error(err, "system.repository.UserNotFound")
}
