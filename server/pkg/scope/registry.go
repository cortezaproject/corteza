package scope

import (
	"fmt"
	"sync"
	"time"
)

// Runtime holds a Scope and its associated services.
// S is the caller-defined services type; the registry imposes no constraints on it.
// Scope.TenantOnly() distinguishes tenant-level from project-level runtimes.
type Runtime[S any] struct {
	Scope    Scope
	Services S

	lastUsed time.Time
	mu       sync.RWMutex
}

func (r *Runtime[S]) touch() {
	r.mu.Lock()
	r.lastUsed = time.Now()
	r.mu.Unlock()
}

// tenantEntry holds a tenant-level Runtime and its project Runtimes.
type tenantEntry[S any] struct {
	runtime  *Runtime[S]
	projects map[uint64]*Runtime[S]
	mu       sync.RWMutex
}

// ScopeRegistry is the global singleton two-level nested runtime index.
// Callers must never cache *Runtime pointers — the registry owns lifecycle.
type ScopeRegistry[S any] struct {
	tenants map[uint64]*tenantEntry[S]
	mu      sync.RWMutex

	newTenantRuntime  func(tenantID uint64) (*Runtime[S], error)
	newProjectRuntime func(tenantID, projectID uint64) (*Runtime[S], error)
}

// NewScopeRegistry creates the registry. newTenant and newProject factories are
// called on first access for a given scope and must be safe for concurrent use.
func NewScopeRegistry[S any](
	newTenant func(tenantID uint64) (*Runtime[S], error),
	newProject func(tenantID, projectID uint64) (*Runtime[S], error),
) *ScopeRegistry[S] {
	return &ScopeRegistry[S]{
		tenants:           make(map[uint64]*tenantEntry[S]),
		newTenantRuntime:  newTenant,
		newProjectRuntime: newProject,
	}
}

// Tenant returns the tenant-level Runtime, lazy-initialising on first access.
func (r *ScopeRegistry[S]) Tenant(tenantID uint64) (*Runtime[S], error) {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return nil, err
	}
	entry.runtime.touch()
	return entry.runtime, nil
}

// Project returns the project-level Runtime, lazy-initialising on first access.
func (r *ScopeRegistry[S]) Project(tenantID, projectID uint64) (*Runtime[S], error) {
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

	// Slow path — init under write lock (double-check).
	entry.mu.Lock()
	defer entry.mu.Unlock()

	if rt, ok = entry.projects[projectID]; ok {
		rt.touch()
		return rt, nil
	}

	rt, err = r.newProjectRuntime(tenantID, projectID)
	if err != nil {
		return nil, fmt.Errorf("scope: init project runtime %d/%d: %w", tenantID, projectID, err)
	}
	entry.projects[projectID] = rt
	rt.touch()
	return rt, nil
}

// Projects returns a snapshot of all live project Runtimes for tenantID.
// Used for tenant-admin fan-out and scheduler iteration.
// Callers must not cache the returned pointers.
func (r *ScopeRegistry[S]) Projects(tenantID uint64) ([]*Runtime[S], error) {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return nil, err
	}

	entry.mu.RLock()
	defer entry.mu.RUnlock()

	out := make([]*Runtime[S], 0, len(entry.projects))
	for _, rt := range entry.projects {
		out = append(out, rt)
	}
	return out, nil
}

// EvictProject removes the project Runtime. Safe to call on project deletion.
func (r *ScopeRegistry[S]) EvictProject(tenantID, projectID uint64) error {
	entry, err := r.tenantEntry(tenantID)
	if err != nil {
		return err
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	delete(entry.projects, projectID)
	return nil
}

// EvictTenant removes the tenant Runtime and all its projects.
func (r *ScopeRegistry[S]) EvictTenant(tenantID uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tenants, tenantID)
}

// GC evicts Runtimes idle longer than maxIdle. Call from a background goroutine.
// Tenant runtime is only evicted after all its project runtimes are gone.
func (r *ScopeRegistry[S]) GC(maxIdle time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for tid, entry := range r.tenants {
		entry.mu.Lock()

		for pid, rt := range entry.projects {
			rt.mu.RLock()
			idle := now.Sub(rt.lastUsed)
			rt.mu.RUnlock()
			if idle > maxIdle {
				delete(entry.projects, pid)
			}
		}

		entry.runtime.mu.RLock()
		tenantIdle := now.Sub(entry.runtime.lastUsed)
		entry.runtime.mu.RUnlock()

		if tenantIdle > maxIdle && len(entry.projects) == 0 {
			delete(r.tenants, tid)
		}

		entry.mu.Unlock()
	}
}

// tenantEntry returns or lazy-inits the tenantEntry for tenantID.
func (r *ScopeRegistry[S]) tenantEntry(tenantID uint64) (*tenantEntry[S], error) {
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

	rt, err := r.newTenantRuntime(tenantID)
	if err != nil {
		return nil, fmt.Errorf("scope: init tenant runtime %d: %w", tenantID, err)
	}

	e = &tenantEntry[S]{
		runtime:  rt,
		projects: make(map[uint64]*Runtime[S]),
	}
	r.tenants[tenantID] = e
	return e, nil
}
