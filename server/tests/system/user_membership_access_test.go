package system

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

// Roles of other users are only visible to users that can read the user and the roles
func TestUserMemberListForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	helpers.DenyMe(h, types.UserRbacResource(0), "read")

	u := h.createUserWithEmail(h.randEmail())
	h.createRoleMember(u.ID, h.repoMakeRole(h.randEmail()).ID)

	h.apiInit().
		Get(fmt.Sprintf("/users/%d/membership", u.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("user.errors.notAllowedToRead")).
		End()
}

func TestUserMemberListUnreadableRoles(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	helpers.AllowMe(h, types.UserRbacResource(0), "read")
	helpers.DenyMe(h, types.RoleRbacResource(0), "read")

	u := h.createUserWithEmail(h.randEmail())
	h.createRoleMember(u.ID, h.repoMakeRole(h.randEmail()).ID)

	h.apiInit().
		Get(fmt.Sprintf("/users/%d/membership", u.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response`, 0)).
		End()
}
