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
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
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

	service.CurrentSettings.Privacy.Mask.Email = true
	service.CurrentSettings.Privacy.Mask.Name = true
	defer func() {
		service.CurrentSettings.Privacy.Mask.Email = false
		service.CurrentSettings.Privacy.Mask.Name = false
	}()

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
	h.allow(types.UserRBACResource.AppendWildcard(), "unmask.email")

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

	h.allow(types.UserRBACResource.AppendWildcard(), "read")

	h.apiInit().
		Get("/users/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filter`)).
		Assert(jsonpath.Present(`$.response.set`)).
		End()
}

func TestUserListWithPaging(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.secCtx()

	seedCount := 40
	for i := 0; i < seedCount; i++ {
		h.createUserWithEmail(h.randEmail())
	}

	h.allow(types.UserRBACResource.AppendWildcard(), "read")

	var aux = struct {
		Response struct {
			Filter struct {
				NextPage *string
				PrevPage *string
			}
		}
	}{}

	h.apiInit().
		Get("/users/").
		Query("limit", "13").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.filter`)).
		Assert(jsonpath.Present(`$.response.set`)).
		Assert(jsonpath.Len(`$.response.set`, 13)).
		Assert(jsonpath.Present(`$.response.filter.nextPage`)).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response.Filter.NextPage)

	h.apiInit().
		Get("/users/").
		Header("Accept", "application/json").
		Query("limit", "13").
		Query("pageCursor", *aux.Response.Filter.NextPage).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.response.set`, 13)).
		Assert(jsonpath.Present(`$.response.filter.prevPage`)).
		Assert(jsonpath.Present(`$.response.filter.nextPage`)).
		End().
		JSON(&aux)

	h.a.NotNil(aux.Response.Filter.PrevPage)
	h.a.NotNil(aux.Response.Filter.NextPage)

}

func TestUserList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.UserRBACResource.AppendWildcard(), "read")

	h.createUserWithEmail("usr")
	f := h.createUserWithEmail(h.randEmail())

	h.deny(types.UserRBACResource.AppendID(f.ID), "read")

	h.apiInit().
		Get("/users/").
		Header("Accept", "application/json").
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

	h.allow(types.UserRBACResource.AppendWildcard(), "read")

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
		End()
}

func TestUserListQueryEmail(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.secCtx()
	h.allow(types.UserRBACResource.AppendWildcard(), "read")
	h.allow(types.UserRBACResource.AppendWildcard(), "unmask.email")

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
	h.allow(types.UserRBACResource.AppendWildcard(), "read")

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
	h.allow(types.UserRBACResource.AppendWildcard(), "read")

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
	h.allow(newUserWeCanAccess.RBACResource(), "read")

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
		Header("Accept", "application/json").
		FormData("email", h.randEmail()).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create users")).
		End()
}

func TestUserCreate(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.SystemRBACResource, "user.create")

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
		Header("Accept", "application/json").
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
	h.allow(types.UserRBACResource.AppendWildcard(), "update")

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
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this user")).
		End()
}

func TestUserDelete(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.UserRBACResource.AppendWildcard(), "delete")

	u := h.createUserWithEmail(h.randEmail())

	h.apiInit().
		Delete(fmt.Sprintf("/users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestUserLabels(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.allow(types.SystemRBACResource, "user.create")
	h.allow(types.UserRBACResource.AppendWildcard(), "read")
	h.allow(types.UserRBACResource.AppendWildcard(), "update")
	h.allow(types.UserRBACResource.AppendWildcard(), "delete")

	var (
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.User{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/users/",
			types.User{Email: h.randEmail(), Labels: map[string]string{"foo": "bar", "bar": "42"}},
			payload,
		)
		req.NotZero(payload.ID)

		h.a.Equal(payload.Labels["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Labels["bar"], "42",
			"labels must contain bar with value 42")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")

		ID = payload.ID
	})

	t.Run("update", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req     = require.New(t)
			payload = &types.User{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("PUT /users/%d", ID),
			&types.User{ID: ID, Email: h.randEmail(), Labels: map[string]string{"foo": "baz", "baz": "123"}},
			payload,
		)
		req.NotZero(payload.ID)
		//req.Nil(payload.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Labels["foo"], "baz",
			"labels must contain foo with value baz")
		req.NotContains(payload.Labels, "bar",
			"labels must not contain bar")
		req.Equal(payload.Labels["baz"], "123",
			"labels must contain baz with value 123")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)
			set = types.UserSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/users/", &set, url.Values{"labels": []string{"baz=123"}})
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
