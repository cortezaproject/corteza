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
	// records org tree calls with the group the store holds for the user at
	// that moment, so a call made before the write is visible as such
	orgTreeSpy struct {
		store store.Storer
		calls []string
	}

	avatarStub struct{ AttachmentService }
)

func (spy *orgTreeSpy) stored(userID id.ID) uint64 {
	u, err := store.LookupUserByID(context.Background(), spy.store, userID.Num())
	if err != nil {
		return 0
	}
	return u.UserGroupID
}

func (spy *orgTreeSpy) AssignGroupMembers(group id.ID, members ...id.ID) error {
	for _, m := range members {
		spy.calls = append(spy.calls, fmt.Sprintf("assign %d (stored %d)", group.Num(), spy.stored(m)))
	}
	return nil
}

func (spy *orgTreeSpy) RemoveGroupMembers(group id.ID, members ...id.ID) error {
	for _, m := range members {
		spy.calls = append(spy.calls, fmt.Sprintf("remove %d (stored %d)", group.Num(), spy.stored(m)))
	}
	return nil
}

func (spy *orgTreeSpy) CloneRulesByRoleID(context.Context, uint64, ...uint64) error { return nil }

func (spy *orgTreeSpy) AddGroupRole(group id.ID, roles ...id.ID) error {
	for _, r := range roles {
		spy.calls = append(spy.calls, fmt.Sprintf("add role %d to %d", r.Num(), group.Num()))
	}
	return nil
}

func (spy *orgTreeSpy) RemoveGroupRole(group id.ID, roles ...id.ID) error {
	for _, r := range roles {
		spy.calls = append(spy.calls, fmt.Sprintf("remove role %d from %d", r.Num(), group.Num()))
	}
	return nil
}

func (avatarStub) CreateAvatarInitialsAttachment(context.Context, string, string, string) (*types.Attachment, error) {
	return &types.Attachment{ID: nextID(), Meta: types.AttachmentMeta{Original: types.AttachmentFileMeta{Image: &types.AttachmentImageMeta{}}}}, nil
}

func (avatarStub) FindByID(context.Context, uint64) (*types.Attachment, error) {
	return &types.Attachment{Meta: types.AttachmentMeta{Labels: map[string]string{"key": types.AttachmentKindAvatar}}}, nil
}

func seedUserGroup(t *testing.T, s store.Storer, deleted bool) uint64 {
	t.Helper()
	g := &types.UserGroup{ID: nextID(), Handle: fmt.Sprintf("g%d", nextID()), CreatedAt: *now()}
	if deleted {
		g.DeletedAt = now()
	}
	require.NoError(t, store.CreateUserGroup(context.Background(), s, g))
	return g.ID
}

func newUserTreeTestService(t *testing.T) (*user, *orgTreeSpy, store.Storer, context.Context) {
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
		rbac.AllowRule(testRoleID, types.ComponentRbacResource(), "user.create"),
		rbac.AllowRule(testRoleID, types.UserRbacResource(0), "update"),
	))

	caller := &types.User{ID: nextID()}
	caller.SetRoles(testRoleID)
	ctx = a.SetIdentityToContext(ctx, caller)

	spy := &orgTreeSpy{store: s}

	svc := &user{
		ac:    &accessControl{rbac: ac},
		store: s,
		services: &userServices{
			settings: &types.AppSettings{},
			eventbus: eventbus.New(),
			att:      avatarStub{},
			groups:   spy,
		},
	}

	return svc, spy, s, ctx
}

// A user's group decides which group roles RBAC evaluates for them, so the
// live tree follows every write of it, and only once the write has landed.
func TestUser_GroupChangeReachesOrgTree(t *testing.T) {
	svc, spy, s, ctx := newUserTreeTestService(t)

	var (
		req     = require.New(t)
		home    = seedUserGroup(t, s, false)
		target  = seedUserGroup(t, s, false)
		deleted = seedUserGroup(t, s, true)
	)

	req.NoError(svc.onCreate(ctx, &types.User{Email: "mover@us.er", Handle: "mover", UserGroupID: home}))

	u, err := store.LookupUserByEmail(ctx, s, "mover@us.er")
	req.NoError(err)
	req.Equal([]string{fmt.Sprintf("assign %d (stored %d)", home, home)}, spy.calls, "create places the user")

	update := func(group uint64, name string) error {
		cur, err := store.LookupUserByID(ctx, s, u.ID)
		req.NoError(err)
		upd := cur.Clone()
		upd.UserGroupID = group
		upd.Name = name
		_, err = svc.Update(ctx, upd)
		return err
	}

	spy.calls = nil
	req.NoError(update(target, ""))
	req.Equal([]string{fmt.Sprintf("assign %d (stored %d)", target, target)}, spy.calls, "a move lands after the write")

	spy.calls = nil
	req.NoError(update(target, "Renamed"))
	req.Empty(spy.calls, "an edit that keeps the group leaves the tree alone")

	spy.calls = nil
	req.NoError(update(0, ""))
	req.Equal([]string{fmt.Sprintf("remove %d (stored 0)", target)}, spy.calls, "no group takes the user out")

	spy.calls = nil
	req.Error(update(deleted, ""), "a deleted group is refused")
	req.Error(update(nextID(), ""), "a missing group is refused")
	req.Empty(spy.calls)

	stored, err := store.LookupUserByID(ctx, s, u.ID)
	req.NoError(err)
	req.Zero(stored.UserGroupID, "a refused move writes nothing")
}

// Group roles are what a group grants its members, so adding or removing one
// reaches the live tree as well as the store.
func TestRole_GroupRoleChangeReachesOrgTree(t *testing.T) {
	svc, _, s, ctx := newRoleTestService(t)

	var (
		req     = require.New(t)
		spy     = &orgTreeSpy{store: s}
		r       = seedTestRole(t, s, &types.Role{Handle: "grouped", Name: "Grouped"})
		group   = seedUserGroup(t, s, false)
		deleted = seedUserGroup(t, s, true)
	)

	svc.services.rbac = spy

	req.NoError(svc.onMemberAddGroup(ctx, &roleActionProps{}, r.ID, group))
	req.Equal([]string{fmt.Sprintf("add role %d to %d", r.ID, group)}, spy.calls)

	spy.calls = nil
	req.NoError(svc.onMemberRemoveGroup(ctx, &roleActionProps{}, r.ID, group))
	req.Equal([]string{fmt.Sprintf("remove role %d from %d", r.ID, group)}, spy.calls)

	spy.calls = nil
	req.Error(svc.onMemberAddGroup(ctx, &roleActionProps{}, r.ID, deleted), "a deleted group is refused")
	req.Error(svc.onMemberAddGroup(ctx, &roleActionProps{}, r.ID, nextID()), "a missing group is refused")
	req.Empty(spy.calls)

	mm, _, err := store.SearchRoleMembers(ctx, s, types.RoleMemberFilter{RoleID: r.ID})
	req.NoError(err)
	req.Empty(mm, "a refused add writes nothing")
}
