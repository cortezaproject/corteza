# Scope Registry Refactor — Interface-Keyed, Pluggable Per-Scope Services

**Status:** ready to implement. Self-contained handoff doc.

## Goal

Refactor scope registry: one central registry keyed by `(tenantID, projectID)`, owns one `Runtime` per scope. Each `Runtime` holds arbitrary type-safe set of sub-services (NG automation engine, RBAC, future plugins) in **interface-keyed container**, not today's fixed `Services` struct.

Drivers:
- Fixed `Services` struct can't represent runtime-added services
  (system moving toward small core + runtime-pluggable interfaces).
  Provider registration mutex-guarded, allowed after boot;
  late-registered providers apply to runtimes built afterwards (see
  "Provider registration semantics" below).
- Today `AutomationRunner` typed `any`, asserted at every call site.
- No teardown hook; `GC`/`Evict` just `delete()` map entry, leak
  goroutines/state held by per-scope services.

## Architecture decision (settled)

**One central registry, scope is aggregate.** Key = scope; value =
`Runtime` box holding every service for that scope.

```
ScopeRegistry
 ├─ (tenant 5, project 0)  → Runtime{ Container{ ExecutionEngine, AccessControl, ... } }
 ├─ (tenant 5, project 12) → Runtime{ Container{ ExecutionEngine, AccessControl, ... } }
 └─ (tenant 9, project 3)  → Runtime{ Container{ ... } }
```

Rejected alternative: one registry per service type. Scatters single
scope's state across N trees — "evict tenant 5" becomes N coordinated deletes
that can half-fail. Scope — not service — is lifecycle unit, so
it is top-level key.

**On `any`:** runtime plugin boundary is interface seam; open
service set needs top type to live in one container. `any`
quarantined inside `Container`, recovered at boundary by generic free
function `Get[T]` (Go methods can't be generic — type recovery must be free
function). Call sites stay fully statically typed; `any` never leaks past
container.

## Current state (pre-refactor)

Package `server/pkg/scope`:
- `registry.go` — generic `Runtime[S any]`, `ScopeRegistry[S any]`. Two-level
  nested index (tenant → projects), lazy double-checked init, `Tenant`,
  `Project`, `Projects`, `EvictProject`, `EvictTenant`, `GC`, `tenantEntry`.
  Built via two factory closures `newTenant`/`newProject`.
- `services.go` — `Services` struct bag (`RBAC`, `EventBus`, `Corredor`,
  `AgenticRunner`, `LLM`, `DAL`, `AutomationRunner any`, `WorkflowRunner`,
  `Scheduler`) + interface definitions. Only `AutomationRunner` ever
  populated.
- `scope.go` — `Scope{TenantID, ProjectID}`, ctx get/set. **Unchanged by this refactor.**
- `middleware.go` — `TenantScopeMiddleware[S]` / `ProjectScopeMiddleware[S]`,
  `Capabilities`, resolver/validator interfaces.

Consumers (entire blast radius — grep-verified):
- `automation/service/service.go` — builds `DefaultAutomationRegistry =
  scope.NewScopeRegistry[scope.Services](...)`; project factory builds
  engine via `runnerSvc.AutomationService(...)`, stores
  `Services{AutomationRunner: engine}`.
- `automation/service/ng_automation.go` — holds `registry
  *scope.ScopeRegistry[scope.Services]`; `engine(ctx)` resolves via
  `rt.Services.AutomationRunner.(executionEngine)`.

Facts that make this low-risk:
- `*ScopeMiddleware` **not wired into any router yet** (only referenced in
  comment in `system/service/project_scope_resolver.go`).
- Registry `GC` **never called** anywhere.
- **No `pkg/scope` tests** to update.
- Only `.Services.` field access in codebase:
  `ng_automation.go:128`.
- `Projects()` has **no callers** — return-type change (`[]*Runtime[S]` →
  `[]*Runtime`) breaks nothing (grep-verified).

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

- `RegisterTenant`/`RegisterProject` take `r.mu`; registration safe any
  time, including after boot (this makes registry
  runtime-pluggable, per first driver).
- Late-registered provider applies only to runtimes built **after**
  registration. Existing runtimes immutable in composition; to pick up
  new provider, evict scope (next access rebuilds with full set).
- Builders read provider slice via snapshot taken under `r.mu.RLock()`
  (copy slice header). Race-free: `append` in Register never
  writes below snapshot's length, so held snapshot never observes
  partially-written element.

#### Method changes (keep existing locking/double-check structure)

- `Tenant(tid)`, `Project(tid,pid)`, `Projects(tid)`, `tenantEntry(tid)`:
  same control flow, but build via `buildRuntime` (passing provider snapshot
  taken under `r.mu.RLock()`) instead of old factory closures, return
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

- `EvictProject(tid, pid)`: under `entry.mu`, remove runtime from
  `entry.projects`; release lock; **then** call `rt.Close()`. Never close
  while holding registry/entry lock — `Close` can block on worker shutdown,
  would stall every lookup behind that lock.
- `EvictTenant(tid)`: under `r.mu`, remove whole `tenantEntry` from
  `r.tenants`; release lock; then close **all project runtimes first, then
  tenant runtime** (children before parent).
- `GC(maxIdle)`: walk under existing locks but only **collect** idle
  victims, delete map entries; release all locks; then `Close()` each
  victim — project runtimes before their tenant runtime. (Closing is leak
  fix — today GC just deletes. Closing outside locks keeps slow teardown
  from blocking all scope lookups registry-wide.)

#### Eviction vs in-flight work

`Close()` cancels per-runtime ctx — **actively kills in-flight
executions** (old delete-only eviction merely leaked them). `touch()` runs
at acquisition only, so runtime serving one long `ExecuteAndWait` can look
idle whole duration. Until refcount exists, GC `maxIdle` MUST be
well above maximum plausible execution duration. Future in-use
refcount on `Runtime` (GC skips refcount > 0) is proper fix — noted in
follow-ups.

### `pkg/scope/services.go`

Delete `Services` struct and `AutomationRunner any` field.
Per-service interfaces (`executionEngine`, RBAC access interface, …) live in
owning packages, serve as container keys. Keep `services.go` only
if interface genuinely shared across packages; otherwise delete file.

### `pkg/scope/middleware.go`

Drop `[S]` type parameter from both middleware constructors:
`TenantScopeMiddleware(reg *ScopeRegistry, validator TenantValidator)` and
`ProjectScopeMiddleware(reg *ScopeRegistry, resolver ProjectResolver, urlParam string)`.
Bodies unchanged — `reg.Tenant(...)` / `reg.Project(...)` keep same
signatures.

## Consumer changes

### `automation/service/service.go`

Replace `DefaultAutomationRegistry` var type + construction:

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
- No tenant-level provider registered (NG automation project-scoped).
  Tenant runtime has empty container — fine.
- `executionEngine` unexported but both `Set` (here) and `Get`
  (ng_automation.go) in package `service`, so key type matches. Note
  explicit `Set[executionEngine]` type argument — required, see `Set` doc.
- `ctx` passed to `AutomationService` now **per-runtime cancelable
  ctx** (from `buildRuntime`), so engine's workers (governor watcher,
  runtime-manager goroutines all select on `ctx.Done()`) stop when scope
  evicted — improvement over current boot-ctx wiring.

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
   - `Get`/`Set` round-trip + type isolation (two interfaces, one container).
   - `Set` keyed by explicit interface type arg; `Get` under different
     interface misses (guards inference footgun).
   - Lazy build: provider runs once per scope; second access reuses.
   - Late registration: provider registered after runtime built
     appears only in runtimes built afterwards; eviction + re-access rebuilds
     with full provider set.
   - `Close` ordering: reverse insertion order; errors joined.
   - Container close re-entry: `Teardownable` whose `Close` calls
     `Get`/`Set` on same container neither deadlocks nor resurrects state.
   - `GC`/`Evict` call `Close`, remove entry; tenant evicted only after
     projects gone; `Close` runs after locks released (test with
     `Teardownable` that calls back into registry during `Close`).
   - `buildRuntime` rollback: provider returning error closes partial
     runtime, surfaces error.

## Follow-ups (out of scope for this refactor, flag to owner)

- Make automation engine implement `Teardownable.Close` so `GC` actually
  frees `runtime_manager` workers (otherwise teardown no-op for it —
  though ctx cancellation already stops worker goroutines).
- Wire `ScopeRegistry.GC` into background goroutine (nothing calls it today →
  unbounded growth once scopes created per project). **Warning
  from "Eviction vs in-flight work": pick `maxIdle` well above longest
  plausible execution, or land refcount first — GC now kills in-flight
  work, no longer merely leaks.**
- Add in-use refcount to `Runtime` (acquire/release around engine use; GC
  skips refcount > 0) — makes eviction safe regardless of `maxIdle`.
- Add RBAC as next provider. Decision recorded: RBAC at **both** levels —
  tenant runtime holds tenant ruleset; project runtime overlays
  project-scoped rules; resolver walks project → tenant. RBAC currently
  global singleton (`rbac.Global()` / `SetGlobal`) caching rules in memory with
  `Watch` reload goroutine — per-scope instances each need
  tenant/project-filtered `rbacRulesStore`, must implement `Teardownable` to
  stop watcher. Keep `Global()` as fallback for unscoped/boot calls.