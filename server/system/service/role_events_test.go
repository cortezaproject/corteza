package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

type (
	// what the store shows about a role at the moment its change is announced
	announcedRole struct {
		event    string
		handle   string
		deleted  bool
		archived bool
	}

	// stands in for the eventbus so dispatch order is observable; the real one
	// runs handlers in a goroutine, which says nothing about when it was called
	roleEventSpy struct {
		store  store.Storer
		roleID uint64
		seen   []announcedRole
	}
)

func (*roleEventSpy) WaitFor(context.Context, eventbus.Event) error { return nil }

func (spy *roleEventSpy) Dispatch(ctx context.Context, ev eventbus.Event) {
	r, err := store.LookupRoleByID(ctx, spy.store, spy.roleID)
	if err != nil {
		spy.seen = append(spy.seen, announcedRole{event: ev.EventType()})
		return
	}

	spy.seen = append(spy.seen, announcedRole{
		event:    ev.EventType(),
		handle:   r.Handle,
		deleted:  r.DeletedAt != nil,
		archived: r.ArchivedAt != nil,
	})
}

// newRoleTestService wires the role service against an in-memory store, with a
// caller allowed to do anything to roles.
func newRoleTestService(t *testing.T) (*role, *roleEventSpy, store.Storer, context.Context) {
	t.Helper()

	var (
		req = require.New(t)
		ctx = context.Background()

		testRoleID = nextID()
		ac         = rbac.NewService(zap.NewNop(), nil)
	)

	s, err := sqlite.ConnectInMemory(ctx)
	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))

	ac.UpdateRoles(rbac.CommonRole.Make(testRoleID, "test-role"))
	req.NoError(ac.Grant(ctx,
		rbac.AllowRule(testRoleID, types.RoleRbacResource(0), "read"),
		rbac.AllowRule(testRoleID, types.RoleRbacResource(0), "update"),
		rbac.AllowRule(testRoleID, types.RoleRbacResource(0), "delete"),
		rbac.AllowRule(testRoleID, types.RoleRbacResource(0), "members.manage"),
		rbac.AllowRule(testRoleID, types.UserRbacResource(0), "read"),
	))

	caller := &types.User{ID: nextID()}
	caller.SetRoles(testRoleID)
	ctx = a.SetIdentityToContext(ctx, caller)

	spy := &roleEventSpy{store: s}

	svc := &role{
		ac:    &accessControl{rbac: ac},
		store: s,
		services: &roleServices{
			eventbus: spy,
			system:   make(map[string]bool),
			closed:   make(map[string]bool),
		},
	}

	return svc, spy, s, ctx
}

func seedTestRole(t *testing.T, s store.Storer, r *types.Role) *types.Role {
	t.Helper()
	r.ID = nextID()
	r.CreatedAt = *now()
	require.NoError(t, store.CreateRole(context.Background(), s, r))
	return r
}

// The RBAC role registry rebuilds itself by re-reading the roles table when it
// hears one of these events, so an event announced before its own write leaves
// the registry a change behind.
func TestRole_ChangeIsAnnouncedOnlyOnceItIsWritten(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		svc, spy, s, ctx := newRoleTestService(t)
		r := seedTestRole(t, s, &types.Role{Handle: "before", Name: "Before"})
		spy.roleID = r.ID

		_, err := svc.Update(ctx, &types.Role{ID: r.ID, Handle: "after", Name: "After"})
		require.NoError(t, err)

		require.Equal(t, []announcedRole{{event: "afterUpdate", handle: "after"}}, spy.seen)
	})

	t.Run("delete", func(t *testing.T) {
		svc, spy, s, ctx := newRoleTestService(t)
		r := seedTestRole(t, s, &types.Role{Handle: "doomed", Name: "Doomed"})
		spy.roleID = r.ID

		require.NoError(t, svc.DeleteByID(ctx, r.ID))
		require.Equal(t, []announcedRole{{event: "afterDelete", handle: "doomed", deleted: true}}, spy.seen)
	})

	t.Run("undelete", func(t *testing.T) {
		svc, spy, s, ctx := newRoleTestService(t)
		r := seedTestRole(t, s, &types.Role{Handle: "back", Name: "Back", DeletedAt: now()})
		spy.roleID = r.ID

		require.NoError(t, svc.UndeleteByID(ctx, r.ID))
		require.Equal(t, []announcedRole{{event: "afterUpdate", handle: "back"}}, spy.seen)
	})
}

// Membership answers "which roles does this user hold" for the user editor, so
// it has to agree with the session the login path builds.
func TestRole_MembershipSkipsDeletedAndArchivedRoles(t *testing.T) {
	svc, _, s, ctx := newRoleTestService(t)

	var (
		req      = require.New(t)
		userID   = nextID()
		resource = "corteza::system:user/" + strconv.FormatUint(userID, 10)

		active   = seedTestRole(t, s, &types.Role{Handle: "active", Name: "Active"})
		deleted  = seedTestRole(t, s, &types.Role{Handle: "gone", Name: "Gone", DeletedAt: now()})
		archived = seedTestRole(t, s, &types.Role{Handle: "parked", Name: "Parked", ArchivedAt: now()})
	)

	req.NoError(store.CreateUser(ctx, s, &types.User{ID: userID, Email: "held@us.er", CreatedAt: *now()}))

	for _, r := range []*types.Role{active, deleted, archived} {
		req.NoError(store.CreateRoleMember(ctx, s, &types.RoleMember{RoleID: r.ID, Resource: resource}))
	}

	mm, err := svc.onMembership(ctx, &roleActionProps{}, userID)
	req.NoError(err)

	held := make([]uint64, 0, len(mm))
	for _, m := range mm {
		held = append(held, m.RoleID)
	}

	req.Equal([]uint64{active.ID}, held)
}
