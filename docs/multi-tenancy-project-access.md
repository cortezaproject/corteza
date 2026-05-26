# Multi-Tenancy — Project Access & Identification

Companion to `multi-tenancy-design.md` and `multi-tenancy-middleware-and-token.md`.

Covers two concerns:
1. **Project identification** — determining which project a request targets
2. **Project access** — determining whether the user may access it and with what capabilities

---

## Part 1 — Project Identification

### Source of truth

Project is identified by its **handle** (URL-safe slug, unique within a tenant). The `projectID` (uint64) is the internal key; handle is the external-facing identifier.

### Where the handle comes from

URL path segment. All project-scoped API routes carry the project handle:

```
/api/project/:projectHandle/...
```

The middleware extracts `:projectHandle` from the routing context.

### Resolution: handle → projectID

1. Read `tenantID` from `Scope` in request context (set by tenant middleware, always present at this point)
2. Look up project by `(tenantID, handle)` — must be scoped to tenant, never global
3. If not found → 404 (do not leak existence; treat missing and forbidden identically if desired)
4. If found but status is not `active` → 403 (suspended/archived projects are inaccessible)
5. Cache result in-process: key `tenantID:handle`, TTL ~60s. Invalidate on project update/delete.

### Why handle, not ID

Handles are stable, human-readable, and safe to expose in URLs. IDs leak sequence information. Handle resolution adds one cache-warm lookup; acceptable cost.

### Output

After resolution, middleware has:
- `projectID uint64`
- `project.Config.Visibility` — needed for access check below

Set `ProjectID` on the existing `Scope` in context and write back via `SetScopeToContext`.

---

## Part 2 — Project Access

### Access sources

A user may access a project through two paths:

| Source | Condition |
|---|---|
| **Tenant membership** | User has active `TenantMembership` for this tenant AND project visibility is `open` |
| **Project membership** | User has active `ProjectMember` record for this `projectID` |

Either source is sufficient. If both exist, project membership takes precedence for capability resolution (it carries an explicit `RolePreset`).

### Decision logic (in order)

```
1. Load ProjectMember record for (userID, projectID)
   → found and active:
       use ProjectMember.RolePreset for capabilities
       → ACCESS GRANTED

2. Load TenantMembership for (userID, tenantID)
   → found and active AND project.Config.Visibility == "open":
       use default tenant member capabilities (read-only; no write, no approval)
       → ACCESS GRANTED

3. Neither condition met → 403
```

Do not expose which condition failed. 403 in all denial cases.

### Capability resolution

Given `RolePreset`, derive `ProjectCapabilities`:

| RolePreset | canRead | canWrite | canRequestApproval | canGrantApproval |
|---|---|---|---|---|
| `governance-owner` | ✓ | ✓ | ✓ | ✓ |
| `security-owner` | ✓ | ✓ | ✓ | ✓ |
| `developer` | ✓ | ✓ | ✓ | ✗ |
| `junior-developer` | ✓ | ✓ | ✗ | ✗ |
| `member` (fallback) | ✓ | ✗ | ✗ | ✗ |
| _(tenant member, no project record)_ | ✓ | ✗ | ✗ | ✗ |

`ProjectCapabilities` is not stored — derived at runtime from `RolePreset`. Stored on context after resolution.

### Context storage

Define a `projectCapabilitiesCtxKey` (analogous to `scopeCtxKey` in `pkg/scope`). After resolution, write capabilities to context. Handlers and services read from context — no repeat DB lookup per request.

Suggested context value type:

```go
type ProjectCapabilities struct {
    CanRead             bool
    CanWrite            bool
    CanRequestApproval  bool
    CanGrantApproval    bool
}
```

### Capability enforcement

Middleware does **not** enforce capabilities — it only resolves and stores them. Each handler or service call enforces the specific capability it requires by reading from context. Example: a write operation reads `CanWrite` from context and returns 403 if false.

---

## Part 3 — Listing Accessible Projects

Used when a user requests the list of projects they can access (e.g. project switcher UI).

### Query strategy

```
SELECT projects WHERE tenantID = :tenantID AND status = 'active'
AND (
    visibility = 'open'                                     -- tenant members see all open projects
    OR id IN (SELECT projectID FROM project_members         -- explicit membership
              WHERE userID = :userID AND status = 'active')
)
```

Scope to `tenantID` from context. Never return projects from other tenants.

### Result

Return list of projects with the user's resolved `RolePreset` per project (or implicit tenant-member role where no explicit record exists). UI uses this to show role badges and gate actions client-side (server enforces independently).

---

## Caching Strategy Summary

| Data | Cache key | TTL | Invalidate on |
|---|---|---|---|
| Handle → projectID | `tenantID:handle` | 60s | Project update / delete |
| Tenant active status | `tenantID` | 30s | Tenant status change |
| ProjectMember record | `userID:projectID` | 30s | Membership change |

All caches are in-process (per server instance). No distributed cache required at this stage. On membership change, evict the relevant key; do not wait for TTL expiry.

---

## Files affected

| File | Change |
|---|---|
| `server/pkg/scope/middleware.go` | Project identification (handle → ID) and access check logic |
| `server/pkg/scope/scope.go` | Add `ProjectCapabilities` type and context helpers |
| Project store (TBD) | `GetByHandle(tenantID, handle)` method |
| `ProjectMember` store (TBD) | `GetByUserAndProject(userID, projectID)` method |
