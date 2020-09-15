package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) randEmail() string {
	return fmt.Sprintf("%s@test.tld", rs())
}

func (h helper) createUserWithEmail(email string) *types.User {
	return h.createUser(&types.User{Email: email})
}

func (h helper) createUser(user *types.User) *types.User {
	if user.ID == 0 {
		user.ID = id.Next()
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateUser(context.Background(), user))
	return user
}

func (h helper) clearUsers() {
	h.noError(store.TruncateUsers(context.Background(), service.DefaultStore))
}

func TestUserRead(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())

	h.apiInit().
		Get(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Contains(`$.response.email`, "####")).
		Assert(jsonpath.Equal(`$.response.userID`, fmt.Sprintf("%d", u.ID))).
		End()

	u = h.createUserWithEmail(h.randEmail())
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
	h.clearUsers()

	h.secCtx()

	seedCount := 5
	for i := 0; i < seedCount; i++ {
		h.createUserWithEmail(h.randEmail())
	}

	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	h.apiInit().
		Get("/users/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filter`)).
		Assert(jsonpath.Present(`$.response.set`)).
		Assert(jsonpath.Len(`$.response.set`, seedCount)).
		End()
}

func TestUserList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	h.createUserWithEmail("usr")
	f := h.createUserWithEmail(h.randEmail())

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
	h.clearUsers()

	h.secCtx()

	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	h.apiInit().
		Get("/users/").
		Query("query", h.randEmail()).
		Query("email", h.randEmail()).
		Query("name", "John Doe").
		Query("handle", "jdoe").
		Query("username", "jdoe").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filter`)).
		Assert(jsonpath.Len(`$.response.set`, 0)).
		End()
}

func TestUserListQueryEmail(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.secCtx()
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")
	h.allow(types.UserPermissionResource.AppendWildcard(), "unmask.email")

	ee := h.randEmail()
	h.createUserWithEmail(ee)

	h.apiInit().
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
	h.clearUsers()

	h.secCtx()
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	ee := h.randEmail()
	h.createUser(&types.User{
		Email:    "test@test.tld",
		Username: ee,
	})

	h.apiInit().
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
	h.clearUsers()

	h.secCtx()
	h.allow(types.UserPermissionResource.AppendWildcard(), "read")

	h.createUser(&types.User{
		Email:  "test@test.tld",
		Handle: "johnDoe",
	})

	h.apiInit().
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
	h.clearUsers()

	h.secCtx()

	newUserWeCanAccess := h.createUserWithEmail(h.randEmail())
	h.allow(newUserWeCanAccess.PermissionResource(), "read")

	// And one we can not access
	h.createUserWithEmail(h.randEmail())

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
		Assert(jsonpath.Present(`$.response.filter`)).
		Assert(jsonpath.Present(`$.response.set`)).
		Assert(jsonpath.Len(`$.response.set`, 1)).
		End().
		JSON(&aux)
}

func TestUserCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.apiInit().
		Post("/users/").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create users")).
		End()
}

func TestUserCreate(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.SystemPermissionResource, "user.create")

	email := h.randEmail()

	h.apiInit().
		Post("/users/").
		FormData("email", email).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestUserUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", u.ID)).
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this user")).
		End()
}

func TestUserUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())
	h.allow(types.UserPermissionResource.AppendWildcard(), "update")

	newEmail := h.randEmail()

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", u.ID)).
		FormData("email", newEmail).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestUserDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this user")).
		End()
}

func TestUserDelete(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.UserPermissionResource.AppendWildcard(), "delete")

	u := h.createUserWithEmail(h.randEmail())

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
