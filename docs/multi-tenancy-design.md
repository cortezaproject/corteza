# Multi-Tenancy Design — Implementation Spec

## Context

This document describes the multi-tenancy extension to Corteza. It is an **implementation spec** for a backend engineer. Design decisions are final unless noted as open questions.

The `server/pkg/scope` package is already implemented and must be used as-is. Do not redesign it.

---

## Architecture Overview

```
System (global)
  └── Tenant (1..n)
        └── Project (1..n)
```

- **System** — global Corteza instance. Sudo-admin level.
- **Tenant** — top-level isolation boundary. Has its own DAL connection. Data from one tenant is never visible to another.
- **Project** — belongs to exactly one tenant. Secondary isolation unit within the tenant.

Users are stored once globally. Scoping is determined entirely by membership records, not by separate user tables.

---

## Existing Foundation (`server/pkg/scope`)

Already implemented. Do not modify unless extending.

| Type | Purpose |
|---|---|
| `Scope` | Carries `TenantID` and `ProjectID` in request context |
| `ScopeRegistry[S]` | Two-level nested runtime index: tenant → project |
| `Runtime[S]` | Holds a `Scope` and its associated service set |

`Scope` is stored in request context via `SetScopeToContext` / `GetScopeFromContext`.

`ScopeRegistry` lazy-initialises `Runtime[S]` instances per tenant and project on first access, with idle-based GC. Tenant runtime is evicted only after all its project runtimes are gone.

Services in `server/pkg/scope/services.go` define the `Services` interface that `Runtime[S]` is parameterised over.

---

## Tenant

### Concept

Top-level isolation boundary. Every resource in the system belongs to a tenant. Users belong to tenants via `TenantMembership` records.

### Fields

| Field | Type | Notes |
|---|---|---|
| ID | uint64 | Primary key |
| Handle | string | URL-safe slug, unique across system |
| Name | string | Display name |
| Status | enum | `active`, `suspended`, `archived` |
| Config | `TenantConfig` | Behavioral settings |
| Meta | `TenantMeta` | Display-oriented data |
| CreatedAt | time.Time | |
| UpdatedAt | *time.Time | |
| SuspendedAt | *time.Time | |
| DeletedAt | *time.Time | Soft delete |

### TenantConfig

| Field | Notes |
|---|---|
| DAL connection reference | Which DB connection/schema this tenant's data lives in |
| Allowed auth providers | Which login methods are permitted (SSO, password, etc.) |
| Feature flags | Which Corteza features are enabled |
| Quotas | Max users, max projects, storage, API rate limits |
| Default locale / timezone | |

### TenantMeta

Display-only. Description, logo, brand color, tags. System does not act on these.

### TenantMembership

Join record between a user and a tenant.

| Field | Notes |
|---|---|
| UserID | |
| TenantID | |
| Role | Tenant-level role (admin, member, etc.) |
| Status | `active`, `suspended`, `invited` (pending acceptance) |
| CreatedAt | When they joined |
| InvitedBy | UserID of the inviter |

**Constraint:** one user belongs to at most one tenant (for now). Enforce at the application layer.

---

## Project

### Concept

Belongs to exactly one tenant. Has its own isolation boundary within the tenant. Users access projects via `ProjectMember` records. A tenant-level member implicitly has access to all projects in that tenant.

### Fields

| Field | Type | Notes |
|---|---|---|
| ID | uint64 | Primary key |
| TenantID | uint64 | FK to tenant |
| Handle | string | URL-safe slug, unique within tenant |
| Name | string | Display name |
| Status | enum | `active`, `archived`, `suspended` |
| CreatedBy | uint64 | UserID of creator |
| Config | `ProjectConfig` | Behavioral settings |
| Meta | `ProjectMeta` | Display-oriented data |
| CreatedAt | time.Time | |
| UpdatedAt | *time.Time | |
| DeletedAt | *time.Time | Soft delete |

### ProjectConfig

| Field | Notes |
|---|---|
| Visibility | `open` (all tenant members) vs `invite-only` |
| DefaultMemberRole | Role preset assigned to new members automatically |
| Feature flags | Project-specific feature overrides |

### ProjectMeta

Display-only. Description, icon, color, tags. System does not act on these.

---

## ProjectMember

Join record between a user and a project.

| Field | Notes |
|---|---|
| UserID | |
| ProjectID | |
| TenantID | Denormalised for query efficiency |
| RolePreset | ID of the `ProjectMemberRole` preset |
| InvitedBy | UserID of inviter |
| JoinedAt | Timestamp |

**Constraint:** one membership record per user per project.

---

## ProjectMemberRole (Role Presets)

Fixed named presets. Not freeform. Stored by ID on the membership record. Capabilities are derived at runtime from the preset definition — not stored per-member.

| Preset | Read | Write | Request Approval | Grant Approval |
|---|---|---|---|---|
| `governance-owner` | ✓ | ✓ | ✓ | ✓ |
| `security-owner` | ✓ | ✓ | ✓ | ✓ |
| `developer` | ✓ | ✓ | ✓ | ✗ |
| `junior-developer` | ✓ | ✓ | ✗ | ✗ |
| `member` (fallback) | ✓ | ✗ | ✗ | ✗ |

**Key design decision:** ownership is not a single owner field. It is a named role on a membership record. Two distinct ownership axes (`governance-owner`, `security-owner`) can be held by different users simultaneously.

---

## ProjectCapabilities

Not a stored type. Derived at runtime.

Given a `ProjectMember`'s `RolePreset`, resolve:
- `canRead`
- `canWrite`
- `canRequestApproval`
- `canGrantApproval`

Only `governance-owner` and `security-owner` can grant. Only `developer` and above can request. `junior-developer` and `member` can neither request nor grant.

---

## User Scoping Model

Users are stored once globally in the existing `users` table. Scope is purely determined by membership records.

| Scope level | Membership type | Access |
|---|---|---|
| System | No membership record — flag on user or role | All tenants, all projects, sudo |
| Tenant | `TenantMembership` | All projects within that tenant |
| Project | `ProjectMember` | That project only |

Tenant-level membership implicitly grants access to all projects in the tenant. Project-level memberships grant access to specific projects only (used when project visibility is `invite-only` or for explicit role assignment).

---

## Authentication & Scope Resolution Flow

### Step 1 — Authenticate

User authenticates via existing Corteza auth. Result: verified `User` identity. No tenant context yet.

### Step 2 — Tenant resolution at login

After successful auth, resolve the user's `TenantMembership`. Since one user = one tenant (current constraint), this yields exactly one `TenantID`.

Stamp `TenantID` into the access token claims alongside existing `userID` and `roles`.

**Subdomain strategy (recommended):** use `tenant-handle.app.example.com` to signal expected tenant at login. After auth, verify user's resolved tenant matches the subdomain tenant. Mismatch → reject.

### Step 3 — Middleware (every request)

A middleware runs before every authenticated request:

1. Extract `tenantID` from token claims
2. Optionally validate tenant is active (cache tenant status — do not DB-hit every request)
3. Call `scope.SetScopeToContext(ctx, scope.Scope{TenantID: tenantID})`
4. If request is project-scoped (URL carries project handle/ID), additionally set `ProjectID` on the `Scope`

Missing or invalid `tenantID` → 401/403.

### Step 4 — Project access check

When a request targets a project:

1. Read `Scope.ProjectID` from context
2. Check `ProjectMember` record exists for this user + project
3. If no direct membership, check if user has tenant-level membership with access to this project (based on project visibility)
4. Resolve `ProjectCapabilities` from the member's `RolePreset`
5. Enforce capability required by the operation

### Step 5 — ScopeRegistry routing

Services use `ScopeRegistry.Project(tenantID, projectID)` or `ScopeRegistry.Tenant(tenantID)` to get the correct `Runtime[S]` for the request scope. This routes to the correct DAL connection and service configuration for that tenant/project.

Do not bypass the registry. Do not cache `*Runtime` pointers outside the registry.

---

## Data Isolation

- **Tenant level:** each tenant has its own DAL connection (see `DalConnection` in `server/system/types/dal_connection.go`). Tenant config references which connection to use.
- **Project level:** projects share the tenant's DAL connection but are isolated by namespace/schema within it, or by row-level tenantID+projectID scoping on all tables — exact mechanism TBD per store implementation.
- **Cross-project communication:** async only, via event bus at the tenant level. Projects emit events; other projects within the same tenant may subscribe. No synchronous cross-project calls.

---

## Open Questions (resolve before implementing)

1. **Project DAL isolation granularity** — separate schema per project, or row-level `projectID` scoping within tenant schema? Separate schema is stronger but operationally heavier.
2. **Tenant-member project override** — can a tenant-admin be restricted to specific projects (project membership as restriction, not just addition)?
3. **Role preset extensibility** — are the five presets hardcoded, or can tenants define custom presets in future?
4. **Invitation flow** — does an invited user need to accept, or is membership active immediately on invite?

---

## Files Already In Place

| File | Status |
|---|---|
| `server/pkg/scope/scope.go` | Done — `Scope`, context helpers |
| `server/pkg/scope/registry.go` | Done — `ScopeRegistry[S]`, `Runtime[S]`, GC |
| `server/pkg/scope/services.go` | Done — `Services` interface definition |

Start implementation from `server/system/types/` for new types, then stores, then REST handlers.

See [multi-tenancy-middleware-and-token.md](multi-tenancy-middleware-and-token.md) for access token changes and middleware spec.
