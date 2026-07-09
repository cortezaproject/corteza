package scope

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

// Capabilities is the runtime capability set attached to context by the project
// scope middleware. Handlers read it to enforce per-operation capability checks.
//
// It mirrors system/types.ProjectCapabilities but is declared here to keep the
// scope package free of a dependency on system/types (which would create an
// import cycle). The project resolver maps between the two.
type Capabilities struct {
	CanRead            bool
	CanWrite           bool
	CanRequestApproval bool
	CanGrantApproval   bool
}

type projectCtxKey struct{}

// SetCapabilitiesToContext stores resolved project capabilities in ctx.
func SetCapabilitiesToContext(ctx context.Context, c Capabilities) context.Context {
	return context.WithValue(ctx, projectCtxKey{}, c)
}

// GetCapabilitiesFromContext returns the capabilities stored in ctx.
// Zero-value (no capabilities) when unset.
func GetCapabilitiesFromContext(ctx context.Context) Capabilities {
	if c, ok := ctx.Value(projectCtxKey{}).(Capabilities); ok {
		return c
	}
	return Capabilities{}
}

// Capability identifies a single project capability for context enforcement.
type Capability uint8

const (
	CapRead Capability = iota
	CapWrite
	CapRequestApproval
	CapGrantApproval
)

// Has reports whether c grants the given capability.
func (c Capabilities) Has(cap Capability) bool {
	switch cap {
	case CapRead:
		return c.CanRead
	case CapWrite:
		return c.CanWrite
	case CapRequestApproval:
		return c.CanRequestApproval
	case CapGrantApproval:
		return c.CanGrantApproval
	}
	return false
}

// RequireTenantMembership returns nil when the request is operating within a
// known tenant scope. Tenant 0 is the mocked single-tenant and is always
// considered active; when real multi-tenancy lands this function is the single
// place to add the membership-record lookup.
func RequireTenantMembership(ctx context.Context) error {
	if GetScopeFromContext(ctx).TenantID == 0 {
		// Tenant 0 = mocked global tenant or no scope set (bootstrap/migration
		// paths). Both cases are permitted until multi-tenancy is active.
		// TODO(multi-tenancy): verify TenantMembership record for sc.TenantID.
		return nil
	}
	return nil
}

// RequireCapability returns nil when the capabilities on ctx grant cap, otherwise
// an Unauthorized error. Handlers and services call this to enforce the specific
// capability an operation needs. Middleware only resolves and stores caps; it
// does not enforce them.
//
// When the request is not inside a project scope (ProjectID == 0 on the context
// scope) the check is skipped — system-level and bootstrap paths have no project
// scope and must not be blocked by a zero-capability context.
func RequireCapability(ctx context.Context, cap Capability) error {
	if GetScopeFromContext(ctx).ProjectID == 0 {
		return nil
	}
	if GetCapabilitiesFromContext(ctx).Has(cap) {
		return nil
	}
	return errors.Unauthorized("insufficient project capability")
}

// ProjectResolver resolves a project from the request and authorises the user
// against it. Implemented in a package that may import system/types + store;
// injected here to avoid an import cycle.
type ProjectResolver interface {
	// Resolve maps a project handle (or ID) to a projectID within the given
	// tenant, verifies the user has access, and returns the project ID plus the
	// capabilities granted to the user.
	//
	// Returns a NotFound error when the project does not exist in the tenant,
	// or an Unauthorized error when the user has no access.
	Resolve(ctx context.Context, tenantID, userID uint64, handleOrID string) (projectID uint64, caps Capabilities, err error)
}

// TenantValidator reports whether a tenant is active. Injected to avoid an
// import cycle. The middleware caches results to avoid hitting it per request.
type TenantValidator interface {
	IsActive(ctx context.Context, tenantID uint64) (bool, error)
}

// tenantStatusCache is a tiny TTL cache keyed by tenantID.
type tenantStatusCache struct {
	ttl time.Duration
	mu  sync.RWMutex
	m   map[uint64]tenantStatusEntry
}

type tenantStatusEntry struct {
	active  bool
	expires time.Time
}

func newTenantStatusCache(ttl time.Duration) *tenantStatusCache {
	return &tenantStatusCache{ttl: ttl, m: make(map[uint64]tenantStatusEntry)}
}

func (c *tenantStatusCache) get(tenantID uint64) (active, ok bool) {
	c.mu.RLock()
	e, found := c.m[tenantID]
	c.mu.RUnlock()
	if !found || time.Now().After(e.expires) {
		return false, false
	}
	return e.active, true
}

func (c *tenantStatusCache) set(tenantID uint64, active bool) {
	c.mu.Lock()
	c.m[tenantID] = tenantStatusEntry{active: active, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

// TenantScopeMiddleware extracts the tenant from the token, validates it, and
// sets the tenant on the request scope.
//
// validator may be nil, in which case tenants are assumed active.
//
// TODO(multi-tenancy): tenant is mocked to 0 system-wide. TenantFromToken
// returns 0 for current tokens, which is treated as the single active tenant.
func TenantScopeMiddleware(validator TenantValidator) func(http.Handler) http.Handler {
	cache := newTenantStatusCache(30 * time.Second)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, _, _ := jwtauth.FromContext(ctx)
			tenantID := auth.TenantFromToken(token)

			// Validate tenant is active (cached). Tenant 0 is the mocked
			// system tenant and is always considered active.
			if tenantID != 0 && validator != nil {
				active, ok := cache.get(tenantID)
				if !ok {
					var err error
					if active, err = validator.IsActive(ctx, tenantID); err != nil {
						errors.ProperlyServeHTTP(w, r, err, false)
						return
					}
					cache.set(tenantID, active)
				}
				if !active {
					errors.ProperlyServeHTTP(w, r, errors.Unauthorized("tenant is not active"), false)
					return
				}
			}

			ctx = SetScopeToContext(ctx, Scope{TenantID: tenantID})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ProjectScopeMiddleware resolves the project from the URL path, checks the
// user's membership, attaches capabilities to context, and sets ProjectID on
// the scope.
//
// Mounted only on project-scoped route groups. urlParam is the chi URL
// parameter name carrying the project handle or ID (e.g. "projectID").
func ProjectScopeMiddleware(resolver ProjectResolver, urlParam string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			sc := GetScopeFromContext(ctx)

			ref := chi.URLParam(r, urlParam)
			if ref == "" {
				errors.ProperlyServeHTTP(w, r, errors.NotFound("project reference missing from request path"), false)
				return
			}

			userID := auth.GetIdentityFromContext(ctx).Identity()

			projectID, caps, err := resolver.Resolve(ctx, sc.TenantID, userID, ref)
			if err != nil {
				errors.ProperlyServeHTTP(w, r, err, false)
				return
			}

			sc.ProjectID = projectID
			ctx = SetScopeToContext(ctx, sc)
			ctx = SetCapabilitiesToContext(ctx, caps)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
