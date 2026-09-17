package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

type (
	quietBus struct{}

	// records every org tree write, so a refused change can be shown to leave the tree alone
	groupTreeSpy struct{ calls []string }

	userByID struct {
		UserService
		store store.Storer
	}
)

func (quietBus) WaitFor(context.Context, eventbus.Event) error { return nil }
func (quietBus) Dispatch(context.Context, eventbus.Event)      {}

func (spy *groupTreeSpy) UpdateUserGroups(...rbac.GroupMembers) error { return nil }
func (spy *groupTreeSpy) AssignGroupMembers(g id.ID, mm ...id.ID) error {
	spy.calls = append(spy.calls, fmt.Sprintf("assign %d", g.Num()))
	return nil
}
func (spy *groupTreeSpy) AddNode(g id.ID, _ string, _ ...rbac.GroupNodePath) error {
	spy.calls = append(spy.calls, fmt.Sprintf("add %d", g.Num()))
	return nil
}
func (spy *groupTreeSpy) UpdateNode(g id.ID, _ string, _ ...rbac.GroupNodePath) error {
	spy.calls = append(spy.calls, fmt.Sprintf("update %d", g.Num()))
	return nil
}
func (spy *groupTreeSpy) RemoveNode(g id.ID) error {
	spy.calls = append(spy.calls, fmt.Sprintf("remove %d", g.Num()))
	return nil
}

func (u userByID) FindByID(ctx context.Context, ID uint64) (*types.User, error) {
	return store.LookupUserByID(ctx, u.store, ID)
}

func newUserGroupTestService(t *testing.T) (*userGroup, *groupTreeSpy, store.Storer, context.Context) {
	t.Helper()

	var (
		req        = require.New(t)
		ctx        = context.Background()
		testRoleID = nextID()
		ac         = rbac.NewService(zap.NewNop(), nil)
	)

	s, err := sqlite.ConnectInMemory(ctx)
	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))

	ac.UpdateRoles(rbac.CommonRole.Make(testRoleID, "test-role"))
	req.NoError(ac.Grant(ctx,
		rbac.AllowRule(testRoleID, types.UserGroupRbacResource(0), "read"),
		rbac.AllowRule(testRoleID, types.UserGroupRbacResource(0), "update"),
		rbac.AllowRule(testRoleID, types.UserGroupRbacResource(0), "delete"),
		rbac.AllowRule(testRoleID, types.UserGroupRbacResource(0), "members.manage"),
	))

	caller := &types.User{ID: nextID()}
	caller.SetRoles(testRoleID)
	ctx = a.SetIdentityToContext(ctx, caller)

	spy := &groupTreeSpy{}
	svc := &userGroup{
		ac:    &accessControl{rbac: ac},
		store: s,
		services: &userGroupServices{
			eventbus: quietBus{},
			rbac:     spy,
			user:     userByID{store: s},
		},
	}

	return svc, spy, s, ctx
}

func seedUserGroupRow(t *testing.T, s store.Storer, deleted bool, parents ...uint64) *types.UserGroup {
	t.Helper()
	g := &types.UserGroup{ID: nextID(), CreatedAt: *now(), Config: &types.UserGroupConfig{}}
	g.Handle = fmt.Sprintf("g%d", g.ID)
	for _, p := range parents {
		g.Config.Paths = append(g.Config.Paths, types.UserGroupPath{SelfID: p})
	}
	if deleted {
		g.DeletedAt = now()
	}
	require.NoError(t, store.CreateUserGroup(context.Background(), s, g))
	return g
}

func seedGroupMember(t *testing.T, s store.Storer, group uint64, deleted bool) *types.User {
	t.Helper()
	u := &types.User{ID: nextID(), CreatedAt: *now(), UserGroupID: group}
	u.Email = fmt.Sprintf("u%d@us.er", u.ID)
	if deleted {
		u.DeletedAt = now()
	}
	require.NoError(t, store.CreateUser(context.Background(), s, u))
	return u
}

func sameError(t *testing.T, want, got error) {
	t.Helper()
	require.Error(t, got)
	require.Equal(t, want.Error(), got.Error())
}

func withPaths(g *types.UserGroup, paths ...types.UserGroupPath) *types.UserGroup {
	upd := g.Clone()
	upd.Config = &types.UserGroupConfig{Paths: paths}
	return upd
}

// The org tree refuses what it cannot index, and a refusal that lands after the
// write leaves the store and the tree disagreeing — the next start then fails.
// Every such change is refused before anything is written.
func TestUserGroup_RefusesWhatTheOrgTreeCannotTake(t *testing.T) {
	svc, spy, s, ctx := newUserGroupTestService(t)

	var (
		req  = require.New(t)
		root = seedUserGroupRow(t, s, false)
	)

	svc.services.rootUserGroup = id.MustNumID(root.ID)

	stored := func(g *types.UserGroup) *types.UserGroup {
		out, err := store.LookupUserGroupByID(ctx, s, g.ID)
		req.NoError(err)
		return out
	}

	t.Run("delete while members remain", func(t *testing.T) {
		g := seedUserGroupRow(t, s, false, root.ID)
		seedGroupMember(t, s, g.ID, false)
		spy.calls = nil

		sameError(t, UserGroupErrHasMembers(), svc.DeleteByID(ctx, g.ID))
		req.Nil(stored(g).DeletedAt)
		req.Empty(spy.calls)
	})

	t.Run("delete while groups report to it", func(t *testing.T) {
		g := seedUserGroupRow(t, s, false, root.ID)
		seedUserGroupRow(t, s, false, g.ID)
		spy.calls = nil

		sameError(t, UserGroupErrHasChildGroups(), svc.DeleteByID(ctx, g.ID))
		req.Nil(stored(g).DeletedAt)
		req.Empty(spy.calls)
	})

	t.Run("delete when only deleted users and deleted groups point at it", func(t *testing.T) {
		g := seedUserGroupRow(t, s, false, root.ID)
		seedGroupMember(t, s, g.ID, true)
		seedUserGroupRow(t, s, true, g.ID)

		req.NoError(svc.DeleteByID(ctx, g.ID))
		req.NotNil(stored(g).DeletedAt)
	})

	t.Run("report to a deleted or missing group", func(t *testing.T) {
		g := seedUserGroupRow(t, s, false, root.ID)
		gone := seedUserGroupRow(t, s, true, root.ID)
		spy.calls = nil

		_, err := svc.Update(ctx, withPaths(g, types.UserGroupPath{SelfID: gone.ID}))
		sameError(t, UserGroupErrParentNotFound(), err)

		_, err = svc.Update(ctx, withPaths(g, types.UserGroupPath{SelfID: nextID()}))
		sameError(t, UserGroupErrParentNotFound(), err)

		req.Equal(root.ID, stored(g).Config.Paths[0].SelfID)
		req.Empty(spy.calls)
	})

	t.Run("report to itself or to a group below it", func(t *testing.T) {
		a := seedUserGroupRow(t, s, false, root.ID)
		b := seedUserGroupRow(t, s, false, a.ID)
		c := seedUserGroupRow(t, s, false, b.ID)
		spy.calls = nil

		_, err := svc.Update(ctx, withPaths(a, types.UserGroupPath{SelfID: c.ID}))
		sameError(t, UserGroupErrCyclicPath(), err)

		_, err = svc.Update(ctx, withPaths(a, types.UserGroupPath{SelfID: a.ID}))
		sameError(t, UserGroupErrCyclicPath(), err)

		req.Equal(root.ID, stored(a).Config.Paths[0].SelfID)
		req.Empty(spy.calls)
	})

	t.Run("two parent links with the same relationship name", func(t *testing.T) {
		other := seedUserGroupRow(t, s, false, root.ID)
		g := seedUserGroupRow(t, s, false, root.ID)

		_, err := svc.Update(ctx, withPaths(g,
			types.UserGroupPath{SelfID: root.ID, Name: "read"},
			types.UserGroupPath{SelfID: other.ID, Name: "read"},
		))
		sameError(t, UserGroupErrDuplicatePathName(), err)
		req.Len(stored(g).Config.Paths, 1)
	})

	t.Run("a group other than the root with no parent", func(t *testing.T) {
		g := seedUserGroupRow(t, s, false, root.ID)

		_, err := svc.Update(ctx, withPaths(g))
		sameError(t, UserGroupErrMissingSelfID(), err)
	})

	t.Run("the root group saves without a parent", func(t *testing.T) {
		upd := withPaths(root)
		upd.Meta = &types.UserGroupMeta{Short: "Renamed root"}
		spy.calls = nil

		_, err := svc.Update(ctx, upd)
		req.NoError(err)
		req.Equal("Renamed root", stored(root).Meta.Short)
		req.Equal([]string{fmt.Sprintf("update %d", root.ID)}, spy.calls)
	})

	t.Run("restore a group whose parent is deleted", func(t *testing.T) {
		gone := seedUserGroupRow(t, s, true, root.ID)
		g := seedUserGroupRow(t, s, true, gone.ID)
		spy.calls = nil

		sameError(t, UserGroupErrParentNotFound(), svc.UndeleteByID(ctx, g.ID))
		req.NotNil(stored(g).DeletedAt)
		req.Empty(spy.calls)
	})

	t.Run("add a member to a deleted group", func(t *testing.T) {
		home := seedUserGroupRow(t, s, false, root.ID)
		gone := seedUserGroupRow(t, s, true, root.ID)
		u := seedGroupMember(t, s, home.ID, false)
		spy.calls = nil

		sameError(t, UserGroupErrNotFound(), svc.onMemberAdd(ctx, &userGroupActionProps{}, gone.ID, u.ID))

		after, err := store.LookupUserByID(ctx, s, u.ID)
		req.NoError(err)
		req.Equal(home.ID, after.UserGroupID)
		req.Empty(spy.calls)
	})

	t.Run("create under a deleted group", func(t *testing.T) {
		gone := seedUserGroupRow(t, s, true, root.ID)

		err := svc.validate(ctx, &types.UserGroup{Config: &types.UserGroupConfig{Paths: []types.UserGroupPath{{SelfID: gone.ID}}}})
		sameError(t, UserGroupErrParentNotFound(), err)
	})
}

// A group left pointing at a deleted parent by an earlier delete must not keep
// the org tree from being built at start.
func TestUserGroup_LiveParentPaths(t *testing.T) {
	const root, a, b, gone = 1, 2, 3, 9

	live := map[uint64]bool{root: true, a: true, b: true}
	group := func(ID uint64, parents ...uint64) *types.UserGroup {
		g := &types.UserGroup{ID: ID, Config: &types.UserGroupConfig{}}
		for _, p := range parents {
			g.Config.Paths = append(g.Config.Paths, types.UserGroupPath{SelfID: p})
		}
		return g
	}

	t.Run("live links are kept", func(t *testing.T) {
		kept, dropped := liveParentPaths(group(b, a, root), live, root)
		require.Equal(t, []types.UserGroupPath{{SelfID: a}, {SelfID: root}}, kept)
		require.Empty(t, dropped)
	})

	t.Run("a dead link is dropped", func(t *testing.T) {
		kept, dropped := liveParentPaths(group(b, gone, a), live, root)
		require.Equal(t, []types.UserGroupPath{{SelfID: a}}, kept)
		require.Equal(t, []uint64{gone}, dropped)
	})

	t.Run("a group with no live parent goes under the root", func(t *testing.T) {
		kept, dropped := liveParentPaths(group(b, gone), live, root)
		require.Equal(t, []types.UserGroupPath{{SelfID: root}}, kept)
		require.Equal(t, []uint64{gone}, dropped)
	})

	t.Run("the root keeps no parents", func(t *testing.T) {
		kept, dropped := liveParentPaths(group(root), live, root)
		require.Empty(t, kept)
		require.Empty(t, dropped)
	})
}

// A reporting cycle saved before cycles were refused cuts its groups, and every
// group below them, off from the root; the org tree then has no node for them.
func TestUserGroup_BreakReportingCycles(t *testing.T) {
	const root, a, b, c, d = 1, 2, 3, 4, 5

	path := func(ii ...uint64) []types.UserGroupPath {
		out := make([]types.UserGroupPath, len(ii))
		for i, ID := range ii {
			out[i] = types.UserGroupPath{SelfID: ID}
		}
		return out
	}

	t.Run("a tree without cycles is left alone", func(t *testing.T) {
		paths := map[uint64][]types.UserGroupPath{root: nil, a: path(root), b: path(a), c: path(a, b)}
		require.Empty(t, breakReportingCycles(paths, root))
		require.Equal(t, path(a, b), paths[c])
	})

	t.Run("the lowest group of a cycle goes under the root, its descendants keep their links", func(t *testing.T) {
		paths := map[uint64][]types.UserGroupPath{root: nil, a: path(c), b: path(a), c: path(b), d: path(c)}

		require.Equal(t, []uint64{a}, breakReportingCycles(paths, root))
		require.Equal(t, path(root), paths[a])
		require.Equal(t, path(a), paths[b])
		require.Equal(t, path(b), paths[c])
		require.Equal(t, path(c), paths[d])
	})

	t.Run("two separate cycles are both broken", func(t *testing.T) {
		paths := map[uint64][]types.UserGroupPath{root: nil, a: path(b), b: path(a), c: path(d), d: path(c)}

		require.Equal(t, []uint64{a, c}, breakReportingCycles(paths, root))
	})
}
