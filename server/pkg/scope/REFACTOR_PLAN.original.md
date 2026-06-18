# Scope Registry Refactor — Interface-Keyed, Pluggable Per-Scope Services

**Status:** ready to implement. Self-contained handoff doc.

## Goal

Refactor the existing scope registry so a single central registry, keyed by
`(tenantID, projectID)`, owns one `Runtime` per scope. Each `Runtime` holds an
arbitrary, type-safe set of sub-services (NG automation engine, RBAC, future
plugins) in an **interface-keyed container** instead of today's fixed `Services`
struct.

Drivers:
- The fixed `Services` struct can't represent services added at runtime
  (the system is moving toward a small core + runtime-pluggable interfaces).
  Provider registration is therefore mutex-guarded and allowed after boot;
  late-registered providers apply to runtimes built afterwards (see
  "Provider registration semantics" below).
- Today `AutomationRunner` is typed `any` and asserted at every call site.
- No teardown hook exists; `GC`/`Evict` just `delete()` the map entry, leaking
  goroutines/state held by per-scope services.

## Architecture decision (settled)

**One central registry, scope is the aggregate.** Key = scope; value = a
`Runtime` box holding every service for that scope.

```
ScopeRegistry
 ├─ (tenant 5, project 0)  → Runtime{ Container{ ExecutionEngine, AccessControl, ... } }
 ├─ (tenant 5, project 12) → Runtime{ Container{ ExecutionEngine, AccessControl, ... } }
 └─ (tenant 9, project 3)  → Runtime{ Container{ ... } }
```

Rejected alternative: one registry per service type. It scatters a single
scope's state across N trees, so "evict tenant 5" becomes N coordinated deletes
that can half-fail. The scope — not the service — is the unit of lifecycle, so
it is the top-level key.

**On `any`:** a runtime plugin boundary is inherently an interface seam; an open
service set cannot be held in one container without a top type. The `any` is
quarantined inside `Container` and recovered at the boundary by a generic free
function `Get[T]` (Go methods can't be generic, so type recovery must be a free
function). Call sites stay fully statically typed; `any` never leaks past the
container.

## Current state (pre-refactor)

Package `server/pkg/scope`:
- `registry.go` — generic `Runtime[S any]`, `ScopeRegistry[S any]`. Two-level
  nested index (tenant → projects), lazy double-checked init, `Tenant`,
  `Project`, `Projects`, `EvictProject`, `EvictTenant`, `GC`, `tenantEntry`.
  Built via two factory closures `newTenant`/`newProject`.
- `services.go` — `Services` struct bag (`RBAC`, `EventBus`, `Corredor`,
  `AgenticRunner`, `LLM`, `DAL`, `AutomationRunner any`, `WorkflowRunner`,
  `Scheduler`) + interface definitions. Only `AutomationRunner` is ever
  populated.
- `scope.go` — `Scope{TenantID, ProjectID}`, ctx get/set. **Unchanged by this refactor.**
- `middleware.go` — `TenantScopeMiddleware[S]` / `ProjectScopeMiddleware[S]`,
  `Capabilities`, resolver/validator interfaces.

Consumers (entire blast radius — verified by grep):
- `automation/service/service.go` — builds `DefaultAutomationRegistry =
  scope.NewScopeRegistry[scope.Services](...)`; project factory builds the
  engine via `runnerSvc.AutomationService(...)` and stores
  `Services{AutomationRunner: engine}`.
- `automation/service/ng_automation.go` — holds `registry
  *scope.ScopeRegistry[scope.Services]`; `engine(ctx)` resolves via
  `rt.Services.AutomationRunner.(executionEngine)`.

Facts that make this low-risk:
- `*ScopeMiddleware` are **not wired into any router yet** (only referenced in a
  comment in `system/service/project_scope_resolver.go`).
- Registry `GC` is **never called** anywhere.
- There are **no `pkg/scope` tests** to update.
- The only `.Services.` field access in the codebase is
  `ng_automation.go:128`.
- `Projects()` has **no callers**, so its return-type change (`[]*Runtime[S]` →
  `[]*Runtime`) breaks nothing (verified by grep).

## Target API

### `pkg/scope/container.go` (new)

```go
package scope

import (
    "context"
    "errors"
    "reflect"
    "sync"
)

// Teardownable is implemented by components that hold resources (goroutines,
// pools, watchers) needing release when their scope is evicted.
type Teardownable interface {
    Close(ctx context.Context) error
}

// Container holds one instance per interface type for a single scope.
// The any-typed storage is recovered statically via Get[T].
type Container struct {
    mu    sync.RWMutex
    m     map[reflect.Type]any
    order []reflect.Type // insertion order; teardown walks it in reverse
}

func newContainer() *Container {
    return &Container{m: make(map[reflect.Type]any)}
}

func key[T any]() reflect.Type {
    return reflect.TypeOf((*T)(nil)).Elem()
}

// Set stores v under the interface type T. Re-setting T replaces in place
// (does not duplicate the order entry).
//
// ALWAYS pass the type argument explicitly, naming the interface:
// Set[ExecutionEngine](c, eng). If T is left to inference it resolves to v's
// concrete type, the value is keyed under that concrete type, and a later
// Get[SomeInterface] silently misses.
func Set[T any](c *Container, v T) {
    k := key[T]()
    c.mu.Lock()
    defer c.mu.Unlock()
    if _, exists := c.m[k]; !exists {
        c.order = append(c.order, k)
    }
    c.m[k] = v
}

// Get returns the component stored under interface T.
func Get[T any](c *Container) (T, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.m[key[T]()]
    if !ok {
        var zero T
        return zero, false
    }
    return v.(T), true
}

// MustGet panics if T is absent. Use only where a provider guarantees presence.
func MustGet[T any](c *Container) T {
    v, ok := Get[T](c)
    if !ok {
        panic("scope: component not found: " + key[T]().String())
    }
    return v
}

// close tears down every Teardownable in reverse insertion order, joining
// errors. It snapshots and clears the container under the lock, then runs
// teardowns OUTSIDE it: a component's Close may safely call Get on the same
// container (it observes an already-empty container instead of deadlocking
// on the non-reentrant RWMutex).
func (c *Container) close(ctx context.Context) error {
    c.mu.Lock()
    snapshot := make([]any, 0, len(c.order))
    for i := len(c.order) - 1; i >= 0; i-- { // reverse insertion order
        snapshot = append(snapshot, c.m[c.order[i]])
    }
    c.m = make(map[reflect.Type]any)
    c.order = nil
    c.mu.Unlock()

    var errs []error
    for _, v := range snapshot {
        if td, ok := v.(Teardownable); ok {
            if err := td.Close(ctx); err != nil {
                errs = append(errs, err)
            }
        }
    }
    return errors.Join(errs...)
}
```

### `pkg/scope/registry.go` (rewrite — non-generic)

```go
// Provider builds and registers one component into a scope's container.
// Apply must call scope.Set[Iface](c, instance) with the interface it owns.
// Whether a provider is tenant- or project-level is determined solely by
// which Register method it is passed to (no Level field — one source of truth).
type Provider struct {
    Name  string // for logging / error context
    Apply func(ctx context.Context, sc Scope, c *Container) error
}

type Runtime struct {
    Scope Scope

    c        *Container
    cancel   context.CancelFunc // cancels the per-runtime ctx on Close
    lastUsed time.Time
    mu       sync.RWMutex
}

func (r *Runtime) Container() *Container { return r.c }

func (r *Runtime) Close() error {
    if r.cancel != nil {
        r.cancel()
    }
    return r.c.close(context.Background()) // or a short-deadline ctx
}

func (r *Runtime) touch() { /* unchanged: lock, set lastUsed = time.Now() */ }

type ScopeRegistry struct {
    baseCtx context.Context
    tenants map[uint64]*tenantEntry
    mu      sync.RWMutex

    // Guarded by mu. May grow after boot (runtime-pluggable services);
    // see "Provider registration semantics".
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
```

`tenantEntry` becomes non-generic:

```go
type tenantEntry struct {
    runtime  *Runtime
    projects map[uint64]*Runtime
    mu       sync.RWMutex
}
```

#### Provider registration semantics

- `RegisterTenant`/`RegisterProject` take `r.mu`; registration is safe at any
  time, including after boot (this is what makes the registry
  runtime-pluggable, per the first driver).
- A late-registered provider applies only to runtimes built **after**
  registration. Existing runtimes are immutable in composition; to pick up a
  new provider, evict the scope (next access rebuilds with the full set).
- Builders read the provider slice via a snapshot taken under `r.mu.RLock()`
  (copy the slice header). This is race-free: `append` in Register never
  writes below the snapshot's length, so a held snapshot never observes a
  partially-written element.

#### Method changes (keep existing locking/double-check structure)

- `Tenant(tid)`, `Project(tid,pid)`, `Projects(tid)`, `tenantEntry(tid)`:
  same control flow, but build via `buildRuntime` (passing a provider snapshot
  taken under `r.mu.RLock()`) instead of the old factory closures, and return
  `*Runtime` (no type param).
- **new** `buildRuntime`:

```go
func (r *ScopeRegistry) buildRuntime(sc Scope, providers []Provider) (*Runtime, error) {
    ctx, cancel := context.WithCancel(r.baseCtx)
    rt := &Runtime{Scope: sc, c: newContainer(), cancel: cancel}
    for _, p := range providers {
        if err := p.Apply(ctx, sc, rt.c); err != nil {
            _ = rt.Close() // roll back partially-built runtime
            return nil, fmt.Errorf("scope: provider %q failed for %v: %w", p.Name, sc, err)
        }
    }
    rt.touch()
    return rt, nil
}
```

- `EvictProject(tid, pid)`: under `entry.mu`, remove the runtime from
  `entry.projects`; release the lock; **then** call `rt.Close()`. Never close
  while holding a registry/entry lock — `Close` can block on worker shutdown
  and would stall every lookup behind that lock.
- `EvictTenant(tid)`: under `r.mu`, remove the whole `tenantEntry` from
  `r.tenants`; release the lock; then close **all project runtimes first, then
  the tenant runtime** (children before parent).
- `GC(maxIdle)`: walk under the existing locks but only **collect** idle
  victims and delete their map entries; release all locks; then `Close()` each
  victim — project runtimes before their tenant runtime. (Closing is the leak
  fix — today GC just deletes. Closing outside the locks keeps a slow teardown
  from blocking all scope lookups registry-wide.)

#### Eviction vs in-flight work

`Close()` cancels the per-runtime ctx, which **actively kills in-flight
executions** (the old delete-only eviction merely leaked them). `touch()` runs
at acquisition only, so a runtime serving one long `ExecuteAndWait` can look
idle for its whole duration. Until a refcount exists, the GC `maxIdle` MUST be
chosen well above the maximum plausible execution duration. A future in-use
refcount on `Runtime` (GC skips refcount > 0) is the proper fix — noted in
follow-ups.

### `pkg/scope/services.go`

Delete the `Services` struct and the `AutomationRunner any` field. The
per-service interfaces (`executionEngine`, RBAC access interface, …) live in
their owning packages and serve as the container keys. Keep `services.go` only
if an interface is genuinely shared across packages; otherwise delete the file.

### `pkg/scope/middleware.go`

Drop the `[S]` type parameter from both middleware constructors:
`TenantScopeMiddleware(reg *ScopeRegistry, validator TenantValidator)` and
`ProjectScopeMiddleware(reg *ScopeRegistry, resolver ProjectResolver, urlParam string)`.
Bodies are unchanged — `reg.Tenant(...)` / `reg.Project(...)` keep the same
signatures.

## Consumer changes

### `automation/service/service.go`

Replace the `DefaultAutomationRegistry` var type and its construction:

```go
// var block
DefaultAutomationRegistry *scope.ScopeRegistry   // drop [scope.Services]
```

```go
// in Initialize, replacing the newEngine closure + NewScopeRegistry call:
execLog := DefaultLogger.Named("automation-execution")

DefaultAutomationRegistry = scope.NewScopeRegistry(ctx)

DefaultAutomationRegistry.RegisterProject(scope.Provider{
    Name: "automation-engine",
    Apply: func(ctx context.Context, sc scope.Scope, c *scope.Container) error {
        eng, err := runnerSvc.AutomationService(ctx,
            execLog.With(
                zap.Uint64("tenantID", sc.TenantID),
                zap.Uint64("projectID", sc.ProjectID),
            ),
            manager.Config{MaxConcurrent: 10},
        )
        if err != nil {
            return err
        }
        scope.Set[executionEngine](c, eng) // executionEngine is the local interface in this package
        return nil
    },
})

DefaultNgAutomation = NgAutomation(DefaultLogger.Named("ng-automation"), c.Corredor, DefaultAutomationRegistry)
```

Notes:
- No tenant-level provider is registered (NG automation is project-scoped). A
  tenant runtime simply has an empty container — fine.
- `executionEngine` is unexported but both the `Set` (here) and the `Get`
  (ng_automation.go) are in package `service`, so the key type matches. Note
  the explicit `Set[executionEngine]` type argument — required, see `Set` doc.
- The `ctx` passed to `AutomationService` is now the **per-runtime cancelable
  ctx** (from `buildRuntime`), so the engine's workers (governor watcher,
  runtime-manager goroutines all select on `ctx.Done()`) stop when the scope
  is evicted — an improvement over the current boot-ctx wiring.

### `automation/service/ng_automation.go`

```go
// struct field
registry *scope.ScopeRegistry // drop [scope.Services]
```

```go
// NgAutomation constructor signature
func NgAutomation(log *zap.Logger, corredorOpt options.CorredorOpt, registry *scope.ScopeRegistry) *ngAutomation
```

```go
func (svc *ngAutomation) engine(ctx context.Context) (executionEngine, error) {
    sc := scope.GetScopeFromContext(ctx)

    rt, err := svc.registry.Project(sc.TenantID, sc.ProjectID)
    if err != nil {
        return nil, err
    }

    eng, ok := scope.Get[executionEngine](rt.Container())
    if !ok {
        return nil, fmt.Errorf("automation engine missing for tenant %d project %d", sc.TenantID, sc.ProjectID)
    }
    return eng, nil
}
```

## Implementation order (each step keeps the tree compiling once its dependents are updated)

1. Add `pkg/scope/container.go`. (No breakage.)
2. Rewrite `pkg/scope/registry.go` (non-generic + providers + `Close`).
3. Edit `pkg/scope/middleware.go` (drop `[S]`).
4. Gut/delete `pkg/scope/services.go`.
5. Rewrite registry construction in `automation/service/service.go`.
6. Fix resolver in `automation/service/ng_automation.go`.
7. `go build ./...`; then add `pkg/scope` unit tests:
   - `Get`/`Set` round-trip + type isolation (two interfaces in one container).
   - `Set` keyed by explicit interface type arg; `Get` under a different
     interface misses (guards the inference footgun).
   - Lazy build: provider runs once per scope; second access reuses.
   - Late registration: a provider registered after a runtime was built
     appears only in runtimes built afterwards; eviction + re-access rebuilds
     with the full provider set.
   - `Close` ordering: reverse insertion order; errors joined.
   - Container close re-entry: a `Teardownable` whose `Close` calls
     `Get`/`Set` on the same container neither deadlocks nor resurrects state.
   - `GC`/`Evict` call `Close` and remove the entry; tenant evicted only after
     its projects are gone; `Close` runs after locks are released (test with a
     `Teardownable` that calls back into the registry during `Close`).
   - `buildRuntime` rollback: a provider returning an error closes the partial
     runtime and surfaces the error.

## Follow-ups (out of scope for this refactor, flag to owner)

- Make the automation engine implement `Teardownable.Close` so `GC` actually
  frees `runtime_manager` workers (otherwise teardown is a no-op for it —
  though ctx cancellation already stops the worker goroutines).
- Wire `ScopeRegistry.GC` into a background goroutine (nothing calls it today →
  unbounded growth once scopes start being created per project). **Warning
  from "Eviction vs in-flight work": pick `maxIdle` well above the longest
  plausible execution, or land the refcount first — GC now kills in-flight
  work, it no longer merely leaks.**
- Add an in-use refcount to `Runtime` (acquire/release around engine use; GC
  skips refcount > 0) to make eviction safe regardless of `maxIdle`.
- Add RBAC as the next provider. Decision recorded: RBAC at **both** levels —
  tenant runtime holds the tenant ruleset; project runtime overlays
  project-scoped rules; resolver walks project → tenant. RBAC is currently a
  global singleton (`rbac.Global()` / `SetGlobal`) caching rules in memory with
  a `Watch` reload goroutine — the per-scope instances each need a
  tenant/project-filtered `rbacRulesStore` and must implement `Teardownable` to
  stop their watcher. Keep `Global()` as the fallback for unscoped/boot calls.
