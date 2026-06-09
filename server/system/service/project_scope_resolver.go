package service

import (
	"context"
	"strconv"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// projectScopeResolver implements scope.ProjectResolver against the store.
//
// Lives in the service layer because it must import system/types + store, which
// the scope package may not (import cycle).
type projectScopeResolver struct {
	store store.Storer
}

// ProjectScopeResolver returns a scope.ProjectResolver backed by the default
// store. Wire it into scope.ProjectScopeMiddleware at router setup.
func ProjectScopeResolver() scope.ProjectResolver {
	return &projectScopeResolver{store: DefaultStore}
}

// Resolve maps a project handle or numeric ID to a project within the tenant,
// authorises the user, and returns the project ID plus capabilities.
func (r *projectScopeResolver) Resolve(ctx context.Context, tenantID, userID uint64, handleOrID string) (uint64, scope.Capabilities, error) {
	var zero scope.Capabilities

	p, err := r.loadProject(ctx, handleOrID)
	if err != nil {
		return 0, zero, err
	}

	// Project must belong to the current tenant.
	//
	// TODO(multi-tenancy): tenant is mocked to 0. Project.TenantID defaults to
	// 0 too, so this passes until real tenants land.
	if p.TenantID != tenantID {
		return 0, zero, ProjectErrNotFound()
	}

	caps, ok, err := r.authorise(ctx, p, userID)
	if err != nil {
		return 0, zero, err
	}
	if !ok {
		return 0, zero, ProjectErrNotAllowedToRead()
	}

	return p.ID, toScopeCaps(caps), nil
}

// loadProject resolves the reference as a numeric ID first, then by handle.
func (r *projectScopeResolver) loadProject(ctx context.Context, ref string) (*types.Project, error) {
	if id, err := strconv.ParseUint(ref, 10, 64); err == nil && id != 0 {
		p, e := store.LookupProjectByID(ctx, r.store, id)
		if errors.IsNotFound(e) {
			return nil, ProjectErrNotFound()
		}
		return p, e
	}

	p, err := store.LookupProjectByHandle(ctx, r.store, ref)
	if errors.IsNotFound(err) {
		return nil, ProjectErrNotFound()
	}
	return p, err
}

// authorise resolves the user's capabilities on the project.
//
//  1. direct ProjectMember record → use its role preset
//  2. no direct membership but project visibility is open → tenant member
//     fallback with the project's default member role (or member)
//  3. otherwise → denied
//
// TODO(multi-tenancy): step 2 should additionally verify the user holds a
// TenantMembership once tenant scoping lands. For now any authenticated user is
// treated as a tenant member of the single mocked tenant.
func (r *projectScopeResolver) authorise(ctx context.Context, p *types.Project, userID uint64) (types.ProjectCapabilities, bool, error) {
	var zero types.ProjectCapabilities

	m, err := store.LookupProjectMemberByProjectIDUserID(ctx, r.store, p.ID, userID)
	if err == nil && m != nil && m.DeletedAt == nil {
		return m.RolePreset.Capabilities(), true, nil
	}
	if err != nil && !errors.IsNotFound(err) {
		return zero, false, err
	}

	// No direct membership — fall back to tenant-level access for open projects.
	if p.Config.Visibility == types.ProjectVisibilityOpen {
		role := p.Config.DefaultMemberRole
		if !role.Valid() {
			role = types.ProjectRoleMember
		}
		return role.Capabilities(), true, nil
	}

	return zero, false, nil
}

func toScopeCaps(c types.ProjectCapabilities) scope.Capabilities {
	return scope.Capabilities{
		CanRead:            c.CanRead,
		CanWrite:           c.CanWrite,
		CanRequestApproval: c.CanRequestApproval,
		CanGrantApproval:   c.CanGrantApproval,
	}
}
