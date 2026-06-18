package scope

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Provider builds and registers one component into a scope's container.
// Apply must call scope.Set[Iface](c, instance) with the interface it owns.
// Whether a provider is tenant- or project-level is determined solely by
// which Register method it is passed to.
type Provider struct {
	Name  string // for logging / error context
	Apply func(ctx context.Context, sc Scope, c *Container) error
}

// Runtime holds a Scope and its associated service container.
// Callers must never cache *Runtime pointers — the registry owns lifecycle.
type Runtime struct {
	Scope Scope

	c        *Container
	cancel   context.CancelFunc
	lastUsed time.Time
	mu       sync.RWMutex
}

func (r *Runtime) Container() *Container { return r.c }

func (r *Runtime) Close() error {
	if r.cancel != nil {
		r.cancel()
	}
	return r.c.close(context.Background())
}

func (r *Runtime) touch() {
	r.mu.Lock()
	r.lastUsed = time.Now()
	r.mu.Unlock()
}

// tenantEntry holds a tenant-level Runtime and its project Runtimes.
type tenantEntry struct {
	runtime  *Runtime
	projects map[uint64]*Runtime
	mu       sync.RWMutex
}

// ScopeRegistry is the central two-level runtime index keyed by (tenantID, projectID).
type ScopeRegistry struct {
	baseCtx context.Context
	tenants map[uint64]*tenantEntry
	mu      sync.RWMutex

	// Guarded by mu. May grow after boot (runtime-pluggable services);
	// late-registered providers apply only to runtimes built afterwards.
	tenantProviders  []Provider
	projectProviders []Provider
}

func NewScopeRegistry(baseCtx context.Context) *ScopeRegistry {
	return &ScopeRegistry{
		baseCtx: baseCtx,
		tenants: make(map[uint64]*tenantEntry),
	}
}

func (r *ScopeRegistry) RegisterTenant(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tenantProviders = append(r.tenantProviders, p)
}

func (r *ScopeRegistry) RegisterProject(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projectProviders = append(r.projectProviders, p)
}

// buildRuntime creates a new Runtime for sc, applying each provider in order.
// On provider error the partially-built runtime is closed and the error returned.
func (r *ScopeRegistry) buildRuntime(sc Scope, providers []Provider) (*Runtime, error) {
	ctx, cancel := context.WithCancel(r.baseCtx)
	rt := &Runtime{Scope: sc, c: newContainer(), cancel: cancel}
	for _, p := range providers {
		if err := p.Apply(ctx, sc, rt.c); err != nil {
			_ = rt.Close()
			return nil, fmt.Errorf("scope: provider %q failed for %v: %w", p.Name, sc, err)
		}
	}
	rt.touch()
	return rt, nil
}

// Tenant returns the tenant-level Runtime, lazy-initialising on first access.
func (r *ScopeRegistry) Tenant(tenantID uint64) (*Runtime, error) {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return nil, err
	}
	entry.runtime.touch()
	return entry.runtime, nil
}

// Project returns the project-level Runtime, lazy-initialising on first access.
func (r *ScopeRegistry) Project(tenantID, projectID uint64) (*Runtime, error) {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return nil, err
	}

	// Fast path.
	entry.mu.RLock()
	rt, ok := entry.projects[projectID]
	entry.mu.RUnlock()
	if ok {
		rt.touch()
		return rt, nil
	}

	// Snapshot project providers before taking entry lock to avoid lock-order
	// inversion: GC holds r.mu → entry.mu; this path must not hold entry.mu → r.mu.
	r.mu.RLock()
	providers := make([]Provider, len(r.projectProviders))
	copy(providers, r.projectProviders)
	r.mu.RUnlock()

	// Slow path — init under write lock (double-check).
	entry.mu.Lock()
	defer entry.mu.Unlock()

	if rt, ok = entry.projects[projectID]; ok {
		rt.touch()
		return rt, nil
	}

	sc := Scope{TenantID: tenantID, ProjectID: projectID}
	rt, err = r.buildRuntime(sc, providers)
	if err != nil {
		return nil, fmt.Errorf("scope: init project runtime %d/%d: %w", tenantID, projectID, err)
	}
	entry.projects[projectID] = rt
	return rt, nil
}

// Projects returns a snapshot of all live project Runtimes for tenantID.
func (r *ScopeRegistry) Projects(tenantID uint64) ([]*Runtime, error) {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return nil, err
	}

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	out := make([]*Runtime, 0, len(entry.projects))
	for _, rt := range entry.projects {
		out = append(out, rt)
	}
	return out, nil
}

// EvictProject removes the project Runtime and closes it.
func (r *ScopeRegistry) EvictProject(tenantID, projectID uint64) error {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return err
	}
	entry.mu.Lock()
	rt := entry.projects[projectID]
	delete(entry.projects, projectID)
	entry.mu.Unlock()
	if rt != nil {
		return rt.Close()
	}
	return nil
}

// EvictTenant removes the tenant Runtime and all its projects, closing each.
// Project runtimes are closed before the tenant runtime.
func (r *ScopeRegistry) EvictTenant(tenantID uint64) {
	r.mu.Lock()
	entry := r.tenants[tenantID]
	delete(r.tenants, tenantID)
	r.mu.Unlock()
	if entry == nil {
		return
	}
	for _, rt := range entry.projects {
		_ = rt.Close()
	}
	_ = entry.runtime.Close()
}

// GC evicts Runtimes idle longer than maxIdle. Map entries are deleted under
// the lock; Close is called after all locks are released so slow teardown does
// not block scope lookups. Project runtimes are closed before their tenant runtime.
func (r *ScopeRegistry) GC(maxIdle time.Duration) {
	type evictedTenant struct {
		projects []*Runtime
		tenant   *Runtime
	}
	var evicted []evictedTenant

	now := time.Now()

	r.mu.Lock()
	for tid, entry := range r.tenants {
		entry.mu.Lock()

		var projVictims []*Runtime
		for pid, rt := range entry.projects {
			rt.mu.RLock()
			idle := now.Sub(rt.lastUsed)
			rt.mu.RUnlock()
			if idle > maxIdle {
				delete(entry.projects, pid)
				projVictims = append(projVictims, rt)
			}
		}

		entry.runtime.mu.RLock()
		tenantIdle := now.Sub(entry.runtime.lastUsed)
		entry.runtime.mu.RUnlock()

		var tenantVictim *Runtime
		if tenantIdle > maxIdle && len(entry.projects) == 0 {
			delete(r.tenants, tid)
			tenantVictim = entry.runtime
		}

		entry.mu.Unlock()

		if len(projVictims) > 0 || tenantVictim != nil {
			evicted = append(evicted, evictedTenant{projects: projVictims, tenant: tenantVictim})
		}
	}
	r.mu.Unlock()

	for _, ev := range evicted {
		for _, rt := range ev.projects {
			_ = rt.Close()
		}
		if ev.tenant != nil {
			_ = ev.tenant.Close()
		}
	}
}

// tenantEntry returns or lazy-inits the tenantEntry for tenantID.
func (r *ScopeRegistry) tenantEntry(tenantID uint64) (*tenantEntry, error) {
	r.mu.RLock()
	e, ok := r.tenants[tenantID]
	r.mu.RUnlock()
	if ok {
		return e, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if e, ok = r.tenants[tenantID]; ok {
		return e, nil
	}

	providers := make([]Provider, len(r.tenantProviders))
	copy(providers, r.tenantProviders)

	sc := Scope{TenantID: tenantID}
	rt, err := r.buildRuntime(sc, providers)
	if err != nil {
		return nil, fmt.Errorf("scope: init tenant runtime %d: %w", tenantID, err)
	}

	e = &tenantEntry{
		runtime:  rt,
		projects: make(map[uint64]*Runtime),
	}
	r.tenants[tenantID] = e
	return e, nil
}
