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
	return h.repoSaveUser(&types.User{Email: email})
}

func (h helper) repoSaveUser(user *types.User) *types.User {
	u, err := h.
		repoUser().
		Create(user)
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
		Assert(jsonpath.Contains(`$.response.email`, "####")).
		Assert(jsonpath.Equal(`$.response.userID`, fmt.Sprintf("%d", u.ID))).
		End()

	u = h.repoMakeUser(h.randEmail())
	h.allow(types.UserPermissionResource.AppendWildcard(), "unmask.email")

	h.apiInit().
		Get(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.email`, u.Email)).
		End()
}

func TestUserListAll(t *testing.T) {
	h := newHelper(t)

	h.secCtx()

	seedCount := 5
	for i := 0; i < seedCount; i++ {
		h.repoMakeUser(h.randEmail())
	}

	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	aux := struct {
		Response *struct{ Filter *types.UserFilter }
	}{}

	h.apiInit().
		Get("/users/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response)
	h.a.NotNil(aux.Response.Filter)

	// we need to test with >= because we're not running this inside a transaction.
	h.a.GreaterOrEqual(int(aux.Response.Filter.Count), seedCount)
}

func TestUserList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	h.repoMakeUser("usr")
	f := h.repoMakeUser(h.randEmail())

	h.deny(types.UserPermissionResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/users/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(fmt.Sprintf(`$.response.set[? @.email=="%s"]`, f.Email))).
		End()
}

func TestUserListQuery(t *testing.T) {
	h := newHelper(t)

	h.secCtx()

	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	aux := struct {
		Response *struct{ Filter *types.UserFilter }
	}{}

	h.apiInit().
		Debug().
		Get("/users/").
		Query("query", h.randEmail()).
		Query("email", h.randEmail()).
		Query("name", "John Doe").
		Query("handle", "jdoe").
		Query("username", "jdoe").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response)
	h.a.NotNil(aux.Response.Filter)
	h.a.GreaterOrEqual(int(aux.Response.Filter.Count), 0)
}

func TestUserListQueryEmail(t *testing.T) {
	h := newHelper(t)

	h.secCtx()
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")
	h.allow(types.UserPermissionResource.AppendWildcard(), "unmask.email")

	ee := h.randEmail()
	h.repoMakeUser(ee)

	h.apiInit().
		Debug().
		Get("/users/").
		Query("email", ee).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.set != null`)).
		End()
}

func TestUserListQueryUsername(t *testing.T) {
	h := newHelper(t)

	h.secCtx()
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	ee := h.randEmail()
	h.repoSaveUser(&types.User{
		Email:    "test@test.tld",
		Username: ee,
	})

	h.apiInit().
		Debug().
		Get("/users/").
		Query("username", ee).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.set != null`)).
		End()
}

func TestUserListQueryHandle(t *testing.T) {
	h := newHelper(t)

	h.secCtx()
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	h.repoSaveUser(&types.User{
		Email:  "test@test.tld",
		Handle: "johnDoe",
	})

	h.apiInit().
		Debug().
		Get("/users/").
		Query("handle", "johnDoe").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.set != null`)).
		End()
}

func TestUserListWithOneAllowed(t *testing.T) {
	h := newHelper(t)

	h.secCtx()

	newUserWeCanAccess := h.repoMakeUser(h.randEmail())
	h.allow(newUserWeCanAccess.PermissionResource(), "read")

	// And one we can not access
	h.repoMakeUser(h.randEmail())

	aux := struct {
		Response *struct {
			Filter *types.UserFilter
		}
	}{}

	h.apiInit().
		Get("/users/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response)
	h.a.NotNil(aux.Response.Filter)

	// we need to test with >= because we're not running this inside a transaction.
	h.a.Equal(1, int(aux.Response.Filter.Count))
}

func TestUserCreateForbidden(t *testing.T) {
	h := newHelper(t)

	h.apiInit().
		Post("/users/").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create user")).
		End()
}

func TestUserCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.SystemPermissionResource, "user.create")

	email := h.randEmail()

	h.apiInit().
		Post("/users/").
		FormData("email", email).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := h.repoUser().FindByEmail(email)
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.True(u.EmailConfirmed)
}

func TestUserUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeUser(h.randEmail())

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", u.ID)).
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update user")).
		End()
}

func TestUserUpdate(t *testing.T) {
	h := newHelper(t)
	u := h.repoMakeUser(h.randEmail())
	h.allow(types.UserPermissionResource.AppendWildcard(), "update")

	newEmail := h.randEmail()

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
		Assert(helpers.AssertError("not allowed to delete user")).
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
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.NotNil(u.DeletedAt)
}
