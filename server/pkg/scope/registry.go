package scope

import "sync"

// Runtime holds the sub-services for a single scope. Sub-services live in its
// Container, keyed by interface type (scope.Set / scope.Get). Build one with
// NewRuntime, populate its container, then hand it to Registry.Set.
type Runtime struct {
	Scope Scope

	c *Container
}

func NewRuntime(sc Scope) *Runtime {
	return &Runtime{Scope: sc, c: newContainer()}
}

func (r *Runtime) Container() *Container { return r.c }

// Registry stores one Runtime per scope. Runtimes are built and populated by the
// caller, then stored via Set. Use NewRegistry; the zero value is not usable.
type Registry struct {
	mu       sync.RWMutex
	runtimes map[Scope]*Runtime
}

func NewRegistry() *Registry {
	return &Registry{runtimes: make(map[Scope]*Runtime)}
}

// Default is the process-wide scope registry. Everything in the system reads
// and writes its per-scope runtimes through it.
var Default = NewRegistry()

// Get returns the stored runtime for sc, or false when none exists.
func (r *Registry) Get(sc Scope) (*Runtime, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, ok := r.runtimes[sc]
	return rt, ok
}

// Set stores rt under its scope, replacing any existing runtime for that scope.
func (r *Registry) Set(rt *Runtime) {
	r.mu.Lock()
	r.runtimes[rt.Scope] = rt
	r.mu.Unlock()
}

// Delete drops the runtime for sc. Teardown of the sub-services it held is the
// caller's concern for now.
func (r *Registry) Delete(sc Scope) {
	r.mu.Lock()
	delete(r.runtimes, sc)
	r.mu.Unlock()
}
