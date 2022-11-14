package system

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/scim"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi/v5"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

// apitest basics, initialize, set handler, add auth
func (h helper) scimApiInit(ffn ...func(*scim.Config)) *apitest.APITest {
	InitTestApp()
	var (
		scimConfig scim.Config
		scimRoutes = chi.NewRouter()
	)

	for _, fn := range ffn {
		fn(&scimConfig)
	}

	scimRoutes.Use(server.BaseMiddleware(false, logger.Default())...)
	scim.Routes(scimRoutes, scimConfig)

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
		Status(http.StatusInternalServerError).
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
		Status(http.StatusOK).
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
		Status(http.StatusOK).
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
	auth := service.Auth(service.AuthOptions{})

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
		Status(http.StatusOK).
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

func TestScimUserReplaceOnExternalId(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()

	// creating a new user and assigning an external ID label to it
	u := h.createUserWithEmail(h.randEmail())
	const externalId = `2819c223-7f76-453a-919d-413861904646`
	h.setLabel(u, "SCIM_externalId", externalId)

	h.scimApiInit(scimSetWithExternalId, scimSetWithUUIDValidator).
		Put(fmt.Sprintf("/Users/%s", externalId)).
		JSON(`{"emails":[{"value":"baz@bar.com"}],"externalId":"` + externalId + `","schemas":["urn:ietf:params:scim:schemas:core:2.0:User"]}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	u, err := store.LookupUserByID(context.Background(), service.DefaultStore, u.ID)
	h.a.NoError(err)
	h.a.NotNil(u)
	h.a.Equal("baz@bar.com", u.Email)
}

func TestScimPatchingGroupMembership(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.clearRoles()
	h.clearRoleMembers()

	isMember := func(r *types.Role, u *types.User) bool {
		mm, _, err := store.SearchRoleMembers(h.secCtx(), service.DefaultStore, types.RoleMemberFilter{RoleID: r.ID, UserID: u.ID})
		h.a.NoError(err)
		return len(mm) > 0
	}

	const (
		user1Id = `00000000-0000-0000-0000-000000000001`
		user2Id = `00000000-0000-0000-0000-000000000002`
		groupId = `00000000-0000-0000-0001-000000000001`
	)

	// creating a new user and assigning an external ID label to it
	u1 := h.createUserWithEmail(h.randEmail())
	h.setLabel(u1, "SCIM_externalId", user1Id)

	u2 := h.createUserWithEmail(h.randEmail())
	h.setLabel(u2, "SCIM_externalId", user2Id)

	r := h.createRole(&types.Role{})
	h.setLabel(r, "SCIM_externalId", groupId)

	h.a.False(isMember(r, u1))
	h.a.False(isMember(r, u2))

	// add only user #1
	h.scimApiInit(scimSetWithExternalId, scimSetWithUUIDValidator).
		Patch(fmt.Sprintf("/Groups/%s", groupId)).
		JSON(fmt.Sprintf(
			`{"Operations":[{"op":"add","path":"members[value eq \"%s\"]"}],"schemas":["urn:ietf:params:scim:schemas:core:2.0:PatchOp"]}`,
			user1Id,
		)).
		Expect(t).
		Status(http.StatusNoContent).
		End()

	h.a.True(isMember(r, u1))
	h.a.False(isMember(r, u2))

	// remove user #1, add user #2
	h.scimApiInit(scimSetWithExternalId, scimSetWithUUIDValidator).
		Patch(fmt.Sprintf("/Groups/%s", groupId)).
		JSON(fmt.Sprintf(
			`{"Operations":[{"op":"add","path":"members[value eq \"%[1]s\"]"},{"op":"remove","path":"members[value eq \"%[2]s\"]"}],"schemas":["urn:ietf:params:scim:schemas:core:2.0:PatchOp"]}`,
			user2Id,
			user1Id,
		)).
		Expect(t).
		Status(http.StatusNoContent).
		End()

	h.a.False(isMember(r, u1))
	h.a.True(isMember(r, u2))

}

func scimSetWithExternalId(c *scim.Config) {
	c.ExternalIdAsPrimary = true
}
func scimSetWithUUIDValidator(c *scim.Config) {
	c.ExternalIdValidator = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
}
