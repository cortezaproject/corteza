package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// userGroupTree stores a root group with a child and a grandchild and registers
// them with the org tree the way the service does on startup
func (h helper) userGroupTree() (root, child, grandchild *types.UserGroup) {
	ctx := context.Background()
	h.noError(store.TruncateUserGroups(ctx, service.DefaultStore))

	mk := func(handle string, parent *types.UserGroup) *types.UserGroup {
		g := &types.UserGroup{ID: id.Next(), Handle: handle, CreatedAt: time.Now(), Config: &types.UserGroupConfig{}}
		if parent != nil {
			g.Config.Paths = []types.UserGroupPath{{SelfID: parent.ID, Name: "reports to"}}
		}
		h.noError(store.CreateUserGroup(ctx, service.DefaultStore, g))
		return g
	}

	root = mk("ug_root", nil)
	child = mk("ug_child", root)
	grandchild = mk("ug_grandchild", child)

	paths := func(g *types.UserGroup) (pp []rbac.GroupNodePath) {
		for _, p := range g.Config.Paths {
			pp = append(pp, rbac.GroupNodePath{SelfID: id.MustNumID(p.SelfID), Name: p.Name})
		}
		return
	}

	h.noError(rbac.Global().UpdateUserGroups(
		rbac.ConvUserGroup(id.MustNumID(root.ID), root.Handle, nil, nil, nil),
		rbac.ConvUserGroup(id.MustNumID(child.ID), child.Handle, nil, nil, paths(child)),
		rbac.ConvUserGroup(id.MustNumID(grandchild.ID), grandchild.Handle, nil, nil, paths(grandchild)),
	))

	return
}

func (h helper) userGroupUpdatePayload(g *types.UserGroup, parent *types.UserGroup) string {
	return helpers.JSON(map[string]interface{}{
		"handle": g.Handle,
		"config": map[string]interface{}{
			"path": []map[string]interface{}{
				{"selfID": fmt.Sprintf("%d", parent.ID), "name": "reports to"},
			},
		},
	})
}

func TestUserGroupDelete_withChildGroups(t *testing.T) {
	h := newHelper(t)
	_, child, _ := h.userGroupTree()
	defer func() { _ = rbac.Global().UpdateUserGroups() }()

	helpers.AllowMe(h, types.UserGroupRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/user-groups/%d", child.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("user-group.errors.hasChildGroups")).
		End()

	// the refused delete must not leave the group soft-deleted
	g, err := store.LookupUserGroupByID(context.Background(), service.DefaultStore, child.ID)
	h.noError(err)
	h.a.Nil(g.DeletedAt)
}

func TestUserGroupDelete_withMembers(t *testing.T) {
	h := newHelper(t)
	_, _, grandchild := h.userGroupTree()
	defer func() { _ = rbac.Global().UpdateUserGroups() }()

	h.clearUsers()
	h.createUser(&types.User{Email: "member@example.tld", Handle: "member", UserGroupID: grandchild.ID})

	helpers.AllowMe(h, types.UserGroupRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/user-groups/%d", grandchild.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("user-group.errors.hasMembers")).
		End()

	g, err := store.LookupUserGroupByID(context.Background(), service.DefaultStore, grandchild.ID)
	h.noError(err)
	h.a.Nil(g.DeletedAt)
}

func TestUserGroupDelete_leaf(t *testing.T) {
	h := newHelper(t)
	_, _, grandchild := h.userGroupTree()
	defer func() { _ = rbac.Global().UpdateUserGroups() }()

	helpers.AllowMe(h, types.UserGroupRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/user-groups/%d", grandchild.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	g, err := store.LookupUserGroupByID(context.Background(), service.DefaultStore, grandchild.ID)
	h.noError(err)
	h.a.NotNil(g.DeletedAt)
}

func TestUserGroupUpdate_cyclicParent(t *testing.T) {
	h := newHelper(t)
	_, child, grandchild := h.userGroupTree()
	defer func() { _ = rbac.Global().UpdateUserGroups() }()

	helpers.AllowMe(h, types.UserGroupRbacResource(0), "update")

	// child would report to its own grandchild
	h.apiInit().
		Put(fmt.Sprintf("/user-groups/%d", child.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(h.userGroupUpdatePayload(child, grandchild)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("user-group.errors.cyclicPath")).
		End()

	// a group can not report to itself either
	h.apiInit().
		Put(fmt.Sprintf("/user-groups/%d", child.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(h.userGroupUpdatePayload(child, child)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("user-group.errors.cyclicPath")).
		End()
}

func TestUserGroupUpdate_reparent(t *testing.T) {
	h := newHelper(t)
	root, _, grandchild := h.userGroupTree()
	defer func() { _ = rbac.Global().UpdateUserGroups() }()

	helpers.AllowMe(h, types.UserGroupRbacResource(0), "update")

	// moving the grandchild directly under the root is a valid change
	h.apiInit().
		Put(fmt.Sprintf("/user-groups/%d", grandchild.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(h.userGroupUpdatePayload(grandchild, root)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	g, err := store.LookupUserGroupByID(context.Background(), service.DefaultStore, grandchild.ID)
	h.noError(err)
	h.a.Equal(root.ID, g.Config.Paths[0].SelfID)
}
