# Corteza Multi-Tenancy / Scope Isolation — Handoff

Brainstorm session. No code written. This doc captures the agreed design so a new session can continue.

## Goal

Add multi-tenancy + project-level isolation to Corteza. Each user belongs to exactly one tenant and may belong to multiple projects within that tenant. Isolation must be enforced consistently across HTTP API access and all internal/async execution paths.

## Decided design

### Scope model
- Fixed 2 levels: `Tenant` (required) + `Project` (optional, nested under tenant).
- Arbitrary-depth/env-defined hierarchy rejected as over-engineering for now.
- `Scope` type in code = `{ TenantID, ProjectID? }`. Keep abstraction so future extension stays local.

### Membership
- User ↔ Tenant: 1:1, mandatory.
- User ↔ Project: N:M via membership table, within tenant.
- Project → Tenant: N:1, mandatory.

### Layering
- **Scope sits above RBAC.** Scope = where you are. RBAC = what you can do there.
- Request pipeline: Auth → Scope → RBAC → Handler → Service → Store.
- Scope check fails → 404/403 before RBAC even runs.

### Trust model
- **Scope in ctx is the trusted source of truth.** Set once at entry, immutable, verified at entry, trusted downstream.
- Two entry-point types both populate ctx the same way:
  - **External (HTTP/gRPC):** middleware parses scope hint (subdomain/header/JWT claim), verifies identity ∈ scope membership, injects `Scope` into ctx.
  - **Internal (jobs, scheduler, triggers, workflows, webhooks, scripts):** scope derived from owning resource or job envelope, verified, injected into ctx.
- Downstream services/store never re-derive or re-verify — they trust ctx.

### Async same-scope-only rule
- Background jobs/events/workflows execute within enqueuer's scope.
- Cross-scope async = explicit named API with its own machinery + audit.
- Queues/topics optionally partitioned per tenant for fairness + isolation.
- Envelope on every async payload: `{ tenant_id, project_id, actor_id, trace_id }`.

### Store guard (core defence)
- Single wrapper layer between service and raw DAL.
- Reads scope from ctx; injects filter into every query; stamps scope on every write.
- Defence-in-depth: after read, assert row scope == ctx scope; before write, same.
- Missing ctx scope → refuse query (panic dev, error prod).
- Escape hatch: `store.WithBypass(ctx, reason)` — allowlisted callers only, audited, lint-gated.
- Codegen friendly: emit guard in Corteza's existing generated store templates.

### Schema impact
- Tenant-only resources (users, roles, settings): add `tenant_id` column.
- Project-capable resources (modules, records, workflows, automations): add `tenant_id` + `project_id NULL`.
- `project_id NULL` = tenant-level resource, visible to all projects within that tenant.
- Composite indexes with `tenant_id` leading.
- FKs only within same scope.

### Cross-project sharing
- Explicit cross-project grants rejected — too complex, too much attack surface.
- Tenant-level resources (`project_id NULL`) visible to all projects within tenant — one extra OR in store guard filter.
- Store guard filter: `WHERE project_id = ctx.project_id OR project_id IS NULL`.
- No cross-tenant sharing under any circumstance.

### RBAC changes
- Rule rows gain `tenant_id` + `project_id` columns.
- Roles per tenant.
- Evaluator reads scope from ctx.
- Tenant-admin role implies all projects in tenant.
- `system` sentinel scope for platform admin operations (audited).

---

## Runtime architecture

### Three-layer runtime model

| Layer | What lives here | Keyed by |
|---|---|---|
| Global | DB pool, HTTP clients, watch runners, event bus | — (singleton, no tenant/project data) |
| TenantRuntime | RBAC, membership, settings, quotas | `TenantID` |
| ProjectRuntime | Automation, scheduler, event routing, workflows | `(TenantID, ProjectID)` |

Global layer = pure infrastructure. Zero tenant/project business data. Never touches scoped data directly.

### ScopeRegistry (not "SvcIndex")

Two-level nested index. Not a flat composite-key map.

```
ScopeRegistry
  └── map[TenantID] → TenantRuntime
            ├── RBAC service
            ├── Membership service
            ├── Settings service
            └── map[ProjectID] → ProjectRuntime
                        ├── Automation service
                        ├── Scheduler
                        └── Event routing
```

Access patterns:

| Want | Call |
|---|---|
| Tenant services | `registry.Tenant(tenantID)` |
| Specific project | `registry.Tenant(tenantID).Project(projectID)` |
| All projects in tenant | `registry.Tenant(tenantID).Projects()` → iterate |
| Tenant-admin fan-out | `registry.Tenant(tenantID).Projects()` → iterate |

Rules:
- `Projects()` returns live runtime pointers — callers must not cache them. Registry owns lifecycle.
- Tenant-admin RBAC must span all projects — evaluator uses `registry.Tenant(tenantID).Projects()` fallback.
- Project deleted → evict + drain `ProjectRuntime` from registry cleanly. Tenant deleted → cascades.
- Lazy init: runtime created on first use, GC'd after idle timeout.
- Registry itself is a global singleton — concurrent-safe map required.

### Scheduler / worker model

Single global worker pool serves all project runtimes. No per-project scheduler instances.

- Per-project job queues (or tagged lanes in shared queue).
- Worker dequeues job, reads envelope `{ tenant_id, project_id, actor_id, trace_id }`, builds scoped ctx, executes.
- Per-project concurrency limits enforced at dequeue time.
- Per-job `recover()` — one bad job cannot crash the pool.

Isolation guarantees from this model:

| Isolation type | Status |
|---|---|
| Data (store guard + ctx scope) | Full |
| Fairness (A can't starve B) | Full — per-project queue limits |
| Logical scope (ctx never bleeds) | Full — envelope → fresh ctx per job |
| Fault (A's panic can't crash B) | Partial — requires per-job recover() |
| Resource (memory/CPU hard limits) | Not provided — shared process |

Resource (OS-level) isolation requires separate processes/pods — out of scope for now.

---

## Layer-by-layer change summary

| Layer | Change | New artifact |
|---|---|---|
| `scope` package | NEW — `Scope` type, ctx helpers, middleware | yes |
| Auth middleware | Unchanged | — |
| Scope middleware (HTTP) | NEW after auth | yes |
| Async entry points | Each gains scope-derivation step | wrapper per subsystem |
| Ctx | Carries Identity + Scope | helpers |
| Services | Read scope from ctx, pass to store | — |
| Store/DAL | Auto-filter + stamp via guard | generated |
| DB schema | `tenant_id` (+ `project_id`) per table | migrations |
| RBAC | Scope-aware rules + eval | column adds |
| Membership | NEW user↔tenant, user↔project | service + tables |
| Resource types | Declare scope binding in `.cue` | metadata |
| Job/event envelope | Carries tenant/project/actor | payload extension |
| Queues/topics | Per-project lane partitioning | infra |
| Workflows/automation | Bound to scope at definition | scope column |
| Corredor/scripts | Scope via gRPC metadata | header convention |
| Codegen (.cue) | Filters + RBAC keys include scope | template change |
| Observability | Logs/metrics/traces tagged | tag injection |
| Quotas (optional) | Per-tenant/project caps | service |
| Lifecycle APIs | Tenant + project CRUD/suspend/export/delete | endpoints |
| ScopeRegistry | NEW — two-level nested runtime registry | yes |
| Global worker pool | NEW — single pool, per-project queue lanes | yes |
| Tests | ≥2-tenant fixtures, leak assertions, CI lint | helpers |

---

## Anti-patterns to forbid

- Global singletons holding tenant data.
- Admin endpoints skipping scope filter "for convenience".
- Background jobs scanning all tenants in one ctx (loop per tenant instead).
- Optional tenant params.
- Tenant inferred mid-request from current user.
- Cross-tenant foreign keys.
- Workers reusing ctx across jobs.
- Job payloads with resource IDs but no envelope.
- Callers caching ProjectRuntime pointers from registry (registry owns lifecycle).
- Cross-project direct references — use tenant-level (`project_id NULL`) resources instead.

---

## Open questions still to lock

1. **Scope hint source for HTTP** — subdomain / path / header / JWT claim — pick one canonical (others optional fallbacks). BLOCKS step 1.
2. **Routing scheme** (subdomain vs path vs header) — affects deployment.
3. **Migration path** for existing single-tenant Corteza deploys (backfill to `root` tenant proposed).
4. **Per-tenant schema vs shared schema** as default (tier-based dial: shared schema → own schema → own DB → own pod). Lock early — affects migration design.
5. Quota composition rules (sum/min/override) — defer until needed.
6. ~~Cross-project share within tenant~~ — **decided**: tenant-level resources (`project_id NULL`), no explicit grants.
7. Re-parenting (move project between tenants) — forbid or gate?
8. Tenant suspension semantics (defer jobs / drop / degraded perms).

---

## Recommended landing order

1. `scope` package + ctx + middleware (no enforcement yet).
2. Membership model + Tenant/Project CRUD services + tables.
3. Backfill migration: all existing rows → `root` tenant, `project_id NULL`.
4. Store guard ON: refuse queries without scope; auto-filter (`project_id = ctx OR NULL`); auto-stamp.
5. Defence-in-depth row-scope assertions in dev/staging.
6. RBAC scope-aware (rule rows + evaluator).
7. **ScopeRegistry skeleton** — global infra layer + TenantRuntime + ProjectRuntime stubs.
8. Async entry points one subsystem at a time: **ng_automation first** → workflow → scheduler → corredor → webhooks → eventbus.
9. Global worker pool with per-project queue lanes + per-job recover().
10. Observability tagging.
11. Per-tenant quotas (optional, later).
12. Lifecycle APIs (tenant/project create/suspend/export/delete).

---

## Files touched this session

None — brainstorm only. User had these files open during discussion (context, not edits):
- `server/system/types/dml_migration.go`
- `server/system/dml_migration.cue`
- `server/system/rest.yaml`

Note: branch `2026.3.x` has in-flight DML migration work (uncommitted). Unrelated to this scope/tenancy design — keep separate.

---

## Where to start next session

Read this doc, then:
1. Map ng_automation components into Global / TenantRuntime / ProjectRuntime layers.
2. Propose `server/pkg/scope/` package skeleton (types, ctx helpers, middleware signatures).
3. Propose `ScopeRegistry` skeleton with `Tenant()` / `Project()` access pattern.
4. Membership tables in `.cue` spec.
5. One worked example of store-guard codegen for an existing resource (e.g. `compose_record`) to validate before rolling out broadly.
