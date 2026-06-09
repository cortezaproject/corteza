package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/crusttech/human/server/system/types"
)

func newTestProjectResolver(t *testing.T) (*projectResolver, store.Storer) {
	t.Helper()
	req := require.New(t)
	ctx := context.Background()

	s, err := sqlite.ConnectInMemory(ctx)
	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))

	return NewProjectResolver(s).(*projectResolver), s
}

func seedProject(t *testing.T, s store.Storer, p *types.Project) {
	t.Helper()
	if p.Status == "" {
		p.Status = types.ProjectStatusActive
	}
	require.NoError(t, store.CreateProject(context.Background(), s, p))
}

func seedMember(t *testing.T, s store.Storer, m *types.ProjectMember) {
	t.Helper()
	require.NoError(t, store.CreateProjectMember(context.Background(), s, m))
}

func TestProjectResolver_ProjectMembershipWins(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, Handle: "alpha", Config: types.ProjectConfig{Visibility: types.ProjectVisibilityInviteOnly}})
	seedMember(t, s, &types.ProjectMember{ID: nextID(), ProjectID: pid, UserID: uid, RolePreset: types.ProjectRoleDeveloper})

	gotID, caps, err := r.Resolve(ctx, 0, uid, "alpha")
	require.NoError(t, err)
	require.Equal(t, pid, gotID)
	// developer: read+write+request, no grant
	require.Equal(t, scope.Capabilities{CanRead: true, CanWrite: true, CanRequestApproval: true}, caps)
}

func TestProjectResolver_ResolveByNumericID(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, Handle: "byid", Config: types.ProjectConfig{Visibility: types.ProjectVisibilityInviteOnly}})
	seedMember(t, s, &types.ProjectMember{ID: nextID(), ProjectID: pid, UserID: uid, RolePreset: types.ProjectRoleMember})

	gotID, _, err := r.Resolve(ctx, 0, uid, strconv.FormatUint(pid, 10))
	require.NoError(t, err)
	require.Equal(t, pid, gotID)
}

func TestProjectResolver_TenantMemberOpenVisibility(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, Handle: "open", Config: types.ProjectConfig{Visibility: types.ProjectVisibilityOpen}})

	// No ProjectMember record; user is a tenant-0 member → read-only.
	gotID, caps, err := r.Resolve(ctx, 0, uid, "open")
	require.NoError(t, err)
	require.Equal(t, pid, gotID)
	require.Equal(t, scope.Capabilities{CanRead: true}, caps)
}

func TestProjectResolver_TenantMemberInviteOnlyDenied(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, Handle: "closed", Config: types.ProjectConfig{Visibility: types.ProjectVisibilityInviteOnly}})

	_, _, err := r.Resolve(ctx, 0, uid, "closed")
	require.True(t, errors.IsUnauthorized(err))
}

func TestProjectResolver_InactiveProjectDenied(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, Handle: "susp", Status: types.ProjectStatusSuspended, Config: types.ProjectConfig{Visibility: types.ProjectVisibilityOpen}})
	seedMember(t, s, &types.ProjectMember{ID: nextID(), ProjectID: pid, UserID: uid, RolePreset: types.ProjectRoleGovernanceOwner})

	_, _, err := r.Resolve(ctx, 0, uid, "susp")
	require.True(t, errors.IsUnauthorized(err))
}

func TestProjectResolver_MissingProjectDenied(t *testing.T) {
	ctx := context.Background()
	r, _ := newTestProjectResolver(t)

	_, _, err := r.Resolve(ctx, 0, nextID(), "ghost")
	require.True(t, errors.IsUnauthorized(err))
}

func TestProjectResolver_CrossTenantDenied(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, TenantID: 7, Handle: "other", Config: types.ProjectConfig{Visibility: types.ProjectVisibilityOpen}})

	// Caller is in tenant 0 — must not resolve a tenant-7 project.
	_, _, err := r.Resolve(ctx, 0, uid, "other")
	require.True(t, errors.IsUnauthorized(err))
}

func TestProjectResolver_SoftDeletedMemberFallsThrough(t *testing.T) {
	ctx := context.Background()
	r, s := newTestProjectResolver(t)

	pid, uid := nextID(), nextID()
	seedProject(t, s, &types.Project{ID: pid, Handle: "sd", Config: types.ProjectConfig{Visibility: types.ProjectVisibilityInviteOnly}})
	now := now()
	seedMember(t, s, &types.ProjectMember{ID: nextID(), ProjectID: pid, UserID: uid, RolePreset: types.ProjectRoleDeveloper, DeletedAt: now})

	// Deleted membership + invite-only → no access.
	_, _, err := r.Resolve(ctx, 0, uid, "sd")
	require.True(t, errors.IsUnauthorized(err))
}
