package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
)

// stores the user that is making the requests
func (h helper) storeCurrentUser() {
	h.cUser.Email = h.randEmail()
	h.cUser.Handle = "self_" + rs()
	h.cUser.Name = "Self"
	h.cUser.CreatedAt = time.Now()

	h.noError(store.CreateUser(context.Background(), service.DefaultStore, h.cUser))
}

// Users without permissions to update users must not be able
// to move themselves into another user group and inherit its roles
func TestUserSelfUpdateUserGroup(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.storeCurrentUser()

	var (
		ctx = context.Background()

		groupID  = id.Next()
		roleID   = id.Next()
		secret   = types.UserRbacResource(id.Next())
		asMyself = rbac.ParamsToSession(ctx, h.cUser.ID, h.roleID)
	)

	// privileged user group; its role is allowed to delete some user
	h.noError(store.TruncateUserGroups(ctx, service.DefaultStore))
	h.noError(store.CreateUserGroup(ctx, service.DefaultStore, &types.UserGroup{
		ID:        groupID,
		Handle:    "privileged",
		CreatedAt: time.Now(),
	}))

	helpers.UpdateRBAC(h.roleID, roleID)
	h.noError(rbac.Global().UpdateUserGroups(
		rbac.ConvUserGroup(id.MustNumID(groupID), "privileged", nil, []id.ID{id.MustNumID(roleID)}, nil),
	))
	h.noError(rbac.Global().Grant(ctx, rbac.AllowRule(roleID, secret, "delete")))

	defer func() {
		_ = rbac.Global().UpdateUserGroups()
	}()

	h.a.NotEqual(rbac.Allow, rbac.Global().Check(asMyself, "delete", rbac.NewResource(secret)), "precondition: no access before the update")

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", h.cUser.ID)).
		Header("Accept", "application/json").
		FormData("email", h.cUser.Email).
		FormData("handle", h.cUser.Handle).
		FormData("name", "Renamed").
		FormData("userGroupID", fmt.Sprintf("%d", groupID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := store.LookupUserByID(ctx, service.DefaultStore, h.cUser.ID)
	h.noError(err)

	h.a.Equal("Renamed", u.Name, "profile fields can still be updated")
	h.a.NotEqual(rbac.Allow, rbac.Global().Check(asMyself, "delete", rbac.NewResource(secret)), "must not inherit roles of the user group")
	h.a.Zero(u.UserGroupID, "user group must not change")
}

// Users without permissions to update users can only change their profile
func TestUserSelfUpdateRestrictedFields(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.storeCurrentUser()

	var (
		ctx   = context.Background()
		email = h.cUser.Email
	)

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", h.cUser.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(&types.User{
			Email:    "vip@corp.example",
			Username: "vip",
			Handle:   "renamed_handle",
			Name:     "Renamed",
			Kind:     types.UserKind("bot"),
			Labels:   map[string]labelTypes.LabelValue{"team": {Val: "finance"}},
		})).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := store.LookupUserByID(ctx, service.DefaultStore, h.cUser.ID)
	h.noError(err)

	h.a.Equal("Renamed", u.Name)
	h.a.Equal("renamed_handle", u.Handle)

	h.a.Equal(email, u.Email, "email must not change")
	h.a.Empty(u.Username, "username must not change")
	h.a.Equal(types.NormalUser, u.Kind, "kind must not change")
	h.a.Empty(helpers.LoadLabelsFromStore(t, service.DefaultStore, u.LabelResourceKind(), u.ID), "labels must not change")
}

// Users with permissions to update can change everything on themselves
func TestUserSelfUpdateWithPermissions(t *testing.T) {
	h := newHelper(t)
	h.clearUsers()
	h.storeCurrentUser()
	helpers.AllowMe(h, types.UserRbacResource(0), "update")

	newEmail := h.randEmail()

	h.apiInit().
		Put(fmt.Sprintf("/users/%d", h.cUser.ID)).
		Header("Accept", "application/json").
		JSON(helpers.JSON(&types.User{
			Email:  newEmail,
			Handle: h.cUser.Handle,
			Name:   "Renamed",
			Labels: map[string]labelTypes.LabelValue{"team": {Val: "finance"}},
		})).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	u, err := store.LookupUserByID(context.Background(), service.DefaultStore, h.cUser.ID)
	h.noError(err)

	h.a.Equal(newEmail, u.Email)
	h.a.Equal("finance", helpers.LoadLabelsFromStore(t, service.DefaultStore, u.LabelResourceKind(), u.ID)["team"].Val)
}
