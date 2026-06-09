package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// projectResolver implements scope.ProjectResolver. It maps a project handle (or
// numeric ID) to a projectID within a tenant, authorises the user against it,
// and returns the capabilities granted by the user's access path.
//
// Access follows docs/multi-tenancy-project-access.md:
//   - Project membership (explicit ProjectMember record) → capabilities from RolePreset.
//   - Tenant membership + open visibility → read-only capabilities.
//   - Neither → Unauthorized (403). All denials are indistinguishable on purpose.
//
// TenantMembership is not yet modelled. Until it lands, every authenticated user
// is treated as a member of tenant 0 (the mocked system tenant); see memberOfTenant.
type projectResolver struct {
	store store.Storer

	projectCache *resolverCache[string, *types.Project]
	memberCache  *resolverCache[string, *types.ProjectMember]
}

// NewProjectResolver builds a resolver over the given store with the default
// cache TTLs from the design doc (project handle 60s, membership 30s).
func NewProjectResolver(s store.Storer) scope.ProjectResolver {
	return &projectResolver{
		store:        s,
		projectCache: newResolverCache[string, *types.Project](60 * time.Second),
		memberCache:  newResolverCache[string, *types.ProjectMember](30 * time.Second),
	}
}

func (r *projectResolver) Resolve(ctx context.Context, tenantID, userID uint64, ref string) (projectID uint64, caps scope.Capabilities, err error) {
	p, err := r.loadProject(ctx, tenantID, ref)
	if err != nil {
		return 0, scope.Capabilities{}, err
	}

	// Suspended / archived projects are inaccessible. Do not leak status — same
	// shape as a missing project.
	if p.Status != types.ProjectStatusActive {
		return 0, scope.Capabilities{}, errNoProjectAccess()
	}

	// 1. Explicit project membership wins and carries an explicit RolePreset.
	if m, e := r.loadMember(ctx, p.ID, userID); e != nil {
		return 0, scope.Capabilities{}, e
	} else if m != nil {
		return p.ID, toScopeCapabilities(m.RolePreset.Capabilities()), nil
	}

	// 2. Tenant membership grants read-only access to open-visibility projects.
	if p.Config.Visibility == types.ProjectVisibilityOpen && r.memberOfTenant(tenantID, userID) {
		return p.ID, toScopeCapabilities(types.ProjectRoleMember.Capabilities()), nil
	}

	// 3. No access path. 403 in all denial cases.
	return 0, scope.Capabilities{}, errNoProjectAccess()
}

// loadProject resolves ref (numeric ID or handle) to a project within tenantID.
// A project from another tenant is treated as not found.
func (r *projectResolver) loadProject(ctx context.Context, tenantID uint64, ref string) (*types.Project, error) {
	key := strconv.FormatUint(tenantID, 10) + ":" + ref
	if p, ok := r.projectCache.get(key); ok {
		return p, nil
	}

	var (
		p   *types.Project
		err error
	)

	// Prefer ID when ref is purely numeric, fall back to handle otherwise.
	if id, e := strconv.ParseUint(ref, 10, 64); e == nil && id != 0 {
		p, err = store.LookupProjectByID(ctx, r.store, id)
	} else {
		p, err = store.LookupProjectByHandle(ctx, r.store, ref)
	}

	if errors.IsNotFound(err) {
		return nil, errNoProjectAccess()
	} else if err != nil {
		return nil, err
	}

	// Enforce tenant isolation: never resolve a project from another tenant.
	if p.TenantID != tenantID {
		return nil, errNoProjectAccess()
	}

	r.projectCache.set(key, p)
	return p, nil
}

// loadMember returns the active ProjectMember for (userID, projectID), or nil
// when none exists. Soft-deleted records are ignored.
func (r *projectResolver) loadMember(ctx context.Context, projectID, userID uint64) (*types.ProjectMember, error) {
	key := strconv.FormatUint(userID, 10) + ":" + strconv.FormatUint(projectID, 10)
	if m, ok := r.memberCache.get(key); ok {
		return m, nil
	}

	m, err := store.LookupProjectMemberByProjectIDUserID(ctx, r.store, projectID, userID)
	if errors.IsNotFound(err) {
		r.memberCache.set(key, nil)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if m.DeletedAt != nil {
		m = nil
	}

	r.memberCache.set(key, m)
	return m, nil
}

// memberOfTenant reports whether the user is a member of the tenant.
//
// TenantMembership is not modelled yet. Tenant 0 is the mocked system tenant and
// every authenticated user belongs to it; any other tenant has no members until
// the model lands.
func (r *projectResolver) memberOfTenant(tenantID, userID uint64) bool {
	return tenantID == 0 && userID != 0
}

// EvictProject clears cached project + membership entries on mutation. Call on
// project update/delete and on membership change.
func (r *projectResolver) EvictProject(tenantID, projectID uint64) {
	r.projectCache.clear()
	r.memberCache.clear()
}

// errNoProjectAccess is the single denial error used for missing project,
// inactive project, and no-access. Callers must not be able to distinguish them.
func errNoProjectAccess() error {
	return errors.Unauthorized("no access to project")
}

func toScopeCapabilities(c types.ProjectCapabilities) scope.Capabilities {
	return scope.Capabilities{
		CanRead:            c.CanRead,
		CanWrite:           c.CanWrite,
		CanRequestApproval: c.CanRequestApproval,
		CanGrantApproval:   c.CanGrantApproval,
	}
}

// resolverCache is a tiny TTL cache. nil values are cacheable (negative caching).
type resolverCache[K comparable, V any] struct {
	ttl time.Duration
	mu  sync.RWMutex
	m   map[K]cacheEntry[V]
}

type cacheEntry[V any] struct {
	val     V
	expires time.Time
}

func newResolverCache[K comparable, V any](ttl time.Duration) *resolverCache[K, V] {
	return &resolverCache[K, V]{ttl: ttl, m: make(map[K]cacheEntry[V])}
}

func (c *resolverCache[K, V]) get(key K) (V, bool) {
	c.mu.RLock()
	e, ok := c.m[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expires) {
		var zero V
		return zero, false
	}
	return e.val, true
}

func (c *resolverCache[K, V]) set(key K, val V) {
	c.mu.Lock()
	c.m[key] = cacheEntry[V]{val: val, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *resolverCache[K, V]) clear() {
	c.mu.Lock()
	c.m = make(map[K]cacheEntry[V])
	c.mu.Unlock()
}
