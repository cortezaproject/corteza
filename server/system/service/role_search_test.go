package service

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// A search by member or by user group lists that holder's roles, not every role.
func TestRole_SearchByMemberOrGroupListsOnlyTheirRoles(t *testing.T) {
	svc, _, s, ctx := newRoleTestService(t)

	var (
		req     = require.New(t)
		userID  = nextID()
		groupID = nextID()

		userRole  = seedTestRole(t, s, &types.Role{Handle: "user-held", Name: "User held"})
		groupRole = seedTestRole(t, s, &types.Role{Handle: "group-held", Name: "Group held"})
		_         = seedTestRole(t, s, &types.Role{Handle: "unheld", Name: "Unheld"})
	)

	req.NoError(store.CreateRoleMember(ctx, s, &types.RoleMember{
		RoleID:   userRole.ID,
		Resource: "corteza::system:user/" + strconv.FormatUint(userID, 10),
	}))
	req.NoError(store.CreateRoleMember(ctx, s, &types.RoleMember{
		RoleID:   groupRole.ID,
		Resource: "corteza::system:user-group/" + strconv.FormatUint(groupID, 10),
	}))

	search := func(f types.RoleFilter) []string {
		req.NoError(svc.beforeSearch(ctx, &f))
		set, _, err := store.SearchRoles(ctx, s, f)
		req.NoError(err)

		handles := make([]string, 0, len(set))
		for _, r := range set {
			handles = append(handles, r.Handle)
		}
		return handles
	}

	req.Equal([]string{"user-held"}, search(types.RoleFilter{MemberID: userID}))
	req.Equal([]string{"group-held"}, search(types.RoleFilter{UserGroupID: groupID}))
	req.Empty(search(types.RoleFilter{UserGroupID: nextID()}), "a group with no roles lists none")

	f := types.RoleFilter{MemberID: userID, UserGroupID: groupID}
	req.Error(svc.beforeSearch(ctx, &f), "member and group together are refused")
}
