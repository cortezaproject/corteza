package system

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/scim"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/go-chi/chi"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
)

var (
	scimRoutes chi.Router
)

// apitest basics, initialize, set handler, add auth
func (h helper) scimApiInit() *apitest.APITest {
	InitTestApp()

	if scimRoutes == nil {
		scimRoutes = chi.NewRouter()
		scimRoutes.Use(server.BaseMiddleware(false, logger.Default())...)
		scim.Routes(scimRoutes)
	}

	return apitest.
		New().
		Handler(scimRoutes)
}

func TestScimUserGet(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())

	h.scimApiInit().
		Get(fmt.Sprintf("/Users/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Contains(`$.schemas`, "urn:ietf:params:scim:schemas:core:2.0:User")).
		Assert(jsonpath.Equal(`$.id`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestScimUserCreate(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.scimApiInit().
		Post("/Users").
		JSON(`{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "userName": "foo",
  "nickName": "baz",
  "emails": [
    {
      "value": "foo@bar.com",
      "primary": true
    },
    {
      "value": "bar@foo.com"
    }
  ]
}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err := store.LookupUserByEmail(context.Background(), service.DefaultStore, "foo@bar.com")
	h.a.NoError(err)
	h.a.Equal("foo", u.Username)
	h.a.Equal("baz", u.Handle)
}

func TestScimUserCreateNoEmail(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.scimApiInit().
		Post("/Users").
		JSON(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func TestScimUserCreateOverwrite(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail("foo@bar.com")

	h.scimApiInit().
		Post("/Users").
		JSON(`{"userName":"UPDATED","emails":[{"value":"foo@bar.com"}],"schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err := store.LookupUserByEmail(context.Background(), service.DefaultStore, "foo@bar.com")
	h.a.NoError(err)
	h.a.Equal("UPDATED", u.Username)
}

func TestScimUserExternalID(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	h.scimApiInit().
		Post("/Users").
		JSON(`{"userName":"foo","emails":[{"value":"foo@bar.com"}],"externalId":"foo42","schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err := store.LookupUserByEmail(context.Background(), service.DefaultStore, "foo@bar.com")
	h.a.NoError(err)
	h.a.Equal("foo", u.Username)

	h.scimApiInit().
		Post("/Users").
		JSON(`{"userName":"baz","emails":[{"value":"baz@bar.com"}],"externalId":"foo42","schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err = store.LookupUserByEmail(context.Background(), service.DefaultStore, "baz@bar.com")
	h.a.NoError(err)
	h.a.Equal("baz", u.Username)

}

func TestScimUserReplace(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())

	h.scimApiInit().
		Put(fmt.Sprintf("/Users/%d", u.ID)).
		JSON(`{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "userName": "bar",
  "emails": [
    {
      "value": "foo@bar.com"
    }
  ]
}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	u, err := store.LookupUserByID(context.Background(), service.DefaultStore, u.ID)
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.Equal("foo@bar.com", u.Email)
}

func TestScimUserPassword(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	service.CurrentSettings.Auth.Internal.Enabled = true
	auth := service.Auth()

	h.scimApiInit().
		Post("/Users").
		JSON(`{"password":"foo$bar$baz 42","emails":[{"value":"baz@bar.com"}],"externalId":"foo42","schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err := auth.InternalLogin(context.Background(), "baz@bar.com", "foo$bar$baz 42")
	h.a.NoError(err)
	h.a.NotNil(u)
}

func TestScimUserDelete(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	u := h.createUserWithEmail(h.randEmail())

	h.scimApiInit().
		Delete(fmt.Sprintf("/Users/%d", u.ID)).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func TestScimGroupGet(t *testing.T) {
	h := newHelper(t)
	h.clearRoles()

	u := h.repoMakeRole()

	h.scimApiInit().
		Get(fmt.Sprintf("/Groups/%d", u.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Contains(`$.schemas`, "urn:ietf:params:scim:schemas:core:2.0:Group")).
		Assert(jsonpath.Equal(`$.id`, fmt.Sprintf("%d", u.ID))).
		End()
}

func TestScimGroupCreate(t *testing.T) {
	h := newHelper(t)
	h.clearRoles()

	h.scimApiInit().
		Post("/Groups").
		JSON(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":"foo"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err := store.LookupRoleByName(context.Background(), service.DefaultStore, "foo")
	h.a.NoError(err)
	h.a.Equal("foo", u.Name)
}

func TestScimGroupExternalId(t *testing.T) {
	h := newHelper(t)
	h.clearRoles()

	h.scimApiInit().
		Post("/Groups").
		JSON(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":"foo","externalId":"grp42"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err := store.LookupRoleByName(context.Background(), service.DefaultStore, "foo")
	h.a.NoError(err)
	h.a.Equal("foo", u.Name)

	h.scimApiInit().
		Post("/Groups").
		JSON(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":"bar","externalId":"grp42"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u, err = store.LookupRoleByName(context.Background(), service.DefaultStore, "bar")
	h.a.NoError(err)
	h.a.Equal("bar", u.Name)
}

func TestScimGroupReplace(t *testing.T) {
	h := newHelper(t)
	h.clearRoles()

	u := h.repoMakeRole()

	h.scimApiInit().
		Put(fmt.Sprintf("/Groups/%d", u.ID)).
		JSON(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":"bar"}`).
		Expect(t).
		End()

	u, err := store.LookupRoleByID(context.Background(), service.DefaultStore, u.ID)
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.Equal("bar", u.Name)
}

func TestScimGroupDelete(t *testing.T) {
	h := newHelper(t)
	h.clearRoles()

	u := h.repoMakeRole(h.randEmail())

	h.scimApiInit().
		Delete(fmt.Sprintf("/Groups/%d", u.ID)).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}
