# Multi-Tenancy — Resource Scoping

Companion to `multi-tenancy-design.md`.

Defines how every resource in the system is scoped to a tenant and/or project.

---

## Isolation policy

**No cross-project data sharing.** Resources belong to exactly one project or one tenant — never shared across project boundaries via NULL columns or implicit fallback. Cross-project communication goes through an explicit API or service call, which re-enters the full auth + scope middleware pipeline.

No cross-tenant sharing under any circumstance.

---

## Scoping model

Every resource falls into one of four scope classes:

| Class | tenantID | projectID | Pattern |
|---|---|---|---|
| **Global** | — | — | Single table, no scope columns. Identity shared across all tenants. |
| **Tenant-scoped** | required | — | `tenant_id` column on table. Null `project_id`. |
| **Project-scoped** | required (denorm.) | required | `tenant_id` + `project_id` columns on table. |
| **Pivot-scoped** | — | — | Resource defined globally; membership/join table binds it to a scope. |

`tenant_id` is always denormalised onto project-scoped tables (do not join through project to get tenant). Denorm enables efficient tenant-wide queries without joining.

---

## Scope guard — the core rule

**Every store query must include scope conditions matching the `Scope` from request context.**

No query touches rows outside the caller's scope. This is enforced in the store layer, not in service or handler. Services receive a `ctx` carrying the `Scope`; stores extract it and append conditions automatically.

Pattern per class:

| Class | WHERE clause appended |
|---|---|
| Global | _(none)_ |
| Tenant-scoped | `tenant_id = :tenantID` |
| Project-scoped | `tenant_id = :tenantID AND project_id = :projectID` — both columns mandatory, never NULL |
| Pivot-scoped | filter via join table scoped to tenant/project |

System-admin requests carry `Scope{TenantID: 0}` — store skips scope guard when `TenantID == 0`.

---

## Resource classification

### Global resources

Defined once, shared across all tenants. No scope columns.

| Resource | Notes |
|---|---|
| `User` | Global identity. Scoped via `TenantMembership` and `ProjectMember` pivot tables. |
| `DalConnection` | System-managed. Referenced by tenant config, not owned by tenant. |
| `LlmProvider` (system-level) | System defines base providers. Tenants may override (see tenant-scoped below). |
| `AuthClient` (system-level) | System-level OAuth clients. |
| `UserGroup` (system-level) | System-wide groups where `IsRoot == true`. |

### Tenant-scoped resources

Belong to a tenant. Visible across all projects within that tenant.

| Resource | File | Notes |
|---|---|---|
| `Tenant` | `system/types/tenant.go` | The tenant itself. |
| `TenantMembership` | `system/types/tenant_membership.go` | User ↔ tenant join. |
| `Project` | `system/types/project.go` | Belongs to tenant. |
| `Application` | `system/types/applications.go` | Tenant-level apps. |
| `AuthClient` (tenant) | `system/types/auth_client.go` | Tenant-specific OAuth clients. |
| `LlmProvider` (tenant) | `system/types/llm_provider.go` | Tenant overrides system providers. |
| `DalSensitivityLevel` | `system/types/dal_sensitivity_level.go` | Data classification, tenant-wide. |
| `Queue` | `system/types/queue.go` | Tenant-level event bus. Per-project lanes for intra-project async. |
| `SmtpConfiguration` | `system/types/smtp_configuration.go` | Tenant mail config. |
| `ConfiguredConnection` | `system/types/configured_connection.go` | Shared external connections within tenant. |
| `UserGroup` (tenant) | `system/types/user_group.go` | Tenant-scoped groups (`IsRoot == false`). |
| `ResourceTranslation` | `system/types/resource_translation.go` | Scoped to where the resource lives; tenant-level translations here. |
| `AppSettings` | `system/types/app_settings.go` | Per-tenant settings overrides. |

### Project-scoped resources

Belong to a project. Isolated from other projects within the same tenant.

| Resource | File | Notes |
|---|---|---|
| `ProjectMember` | `system/types/project_member.go` | User ↔ project join. |
| `Agent` | `system/types/agent.go` | AI agents live in a project. |
| `AiConversation` | `system/types/ai_conversation.go` | Per-project conversation history. |
| `ApigwRoute` | `system/types/apigw_route.go` | API gateway routes per project. |
| `ApigwFilter` | `system/types/apigw_filter.go` | Bound to route, inherits project scope. |
| `Chatbot` | `system/types/chatbot.go` | Per-project chatbot definition. |
| `ChatbotSession` | `system/types/chatbot_session.go` | Bound to chatbot, inherits project scope. |
| `KnowledgeBase` | `system/types/knowledge_base.go` | Per-project knowledge store. |
| `Report` | `system/types/report.go` | Per-project reports. |
| `Template` | `system/types/template.go` | Project-level templates (tenant-level templates TBD). |
| `Reminder` | `system/types/reminder.go` | Per-project, per-user reminders. |
| `Attachment` | `system/types/attachment.go` | Scoped to where it was uploaded. |
| `Notification` | `system/types/notification.go` | Per-project notification config. |
| `DataPrivacyRequest` | `system/types/data_privacy.go` | Per-project data privacy requests. |
| `Workflow` (automation) | `automation/types/workflow.go` | Per-project automations. |
| `Trigger` (automation) | `automation/types/trigger.go` | Bound to workflow, inherits scope. |
| `Session` (automation) | `automation/types/session.go` | Runtime, inherits project scope. |

### Compose resources (all project-scoped)

Compose is inherently project-scoped. The existing `Namespace` maps 1:1 to a `Project`.

| Resource | Notes |
|---|---|
| `Namespace` | Becomes the project's compose container. `NamespaceID == ProjectID` or mapped via FK. |
| `Module` | Scoped via namespace → project. |
| `Record` | Scoped via module → namespace → project. |
| `Chart` | Per-project. |
| `Page`, `PageLayout` | Per-project. |

> **Key decision:** whether `Namespace` is replaced by `Project` or remains as a sub-concept within a project is TBD. For now, treat namespace as project-scoped with `project_id` column.

### Federation resources

Federation nodes operate at system or tenant level. Exposed/shared modules are project-scoped.

| Resource | Scope |
|---|---|
| `Node` | Tenant-scoped |
| `ExposedModule`, `SharedModule` | Project-scoped |
| `ModuleMapping` | Project-scoped |

---

## Implementation approach

### Option A — Column addition (default for new resources)

Add `tenant_id BIGINT` and `project_id BIGINT` columns to the table. Store layer reads `Scope` from context and appends WHERE conditions.

Use for: all resources that are native to Corteza and owned by a project or tenant.

### Option B — Pivot table (for globally-defined resources)

Resource table has no scope columns. A separate join table records which tenants/projects the resource is associated with.

Use for: `User` (already global), potentially `Template` and `LlmProvider` when cross-tenant sharing is needed later.

Current pivot tables:
- `TenantMembership` — User ↔ Tenant
- `ProjectMember` — User ↔ Project

### Option C — Reference inheritance (derived scope)

Resource has no scope columns but is reached only through a scoped parent (e.g. `ApigwFilter` → `ApigwRoute` → project). Scope enforced by requiring parent lookup within scope before child access.

Use sparingly — prefer explicit columns for queryability.

---

## Store layer changes

### Scope extraction helper

Add to store base or middleware:

```go
func scopeFromCtx(ctx context.Context) scope.Scope {
    return scope.GetScopeFromContext(ctx)
}
```

### Query guard pattern

Every store method that lists or looks up scoped resources appends conditions derived from `Scope`:

```
if s.TenantID != 0 {
    q = q.Where("tenant_id = ?", s.TenantID)
}
if s.ProjectID != 0 {
    q = q.Where("project_id = ?", s.ProjectID)
}
```

System admin bypass: `TenantID == 0` skips guard (system scope).

### Lookup by ID must still scope-check

`LookupXByID` must include scope conditions. A valid ID from another tenant must not be resolvable. Return `store.ErrNotFound` on mismatch — do not leak existence.

---

## Migration notes

Existing Corteza deployments have no `tenant_id` / `project_id` columns. Migration strategy:

1. Add nullable `tenant_id` and `project_id` columns to all affected tables
2. Provision a default tenant for existing data
3. Backfill `tenant_id` on all rows with the default tenant ID
4. For project-scoped resources, backfill `project_id` using existing `namespace_id` or a default project per tenant
5. Add NOT NULL constraint after backfill
6. Add indexes on `(tenant_id)` and `(tenant_id, project_id)`

---

## Files to create / modify

| File | Action |
|---|---|
| `server/store/scoped_query.go` | New — scope guard helpers used by all store implementations |
| All affected store files under `server/store/rdbms/` | Modify — add scope conditions to every query |
| DB migrations | New — add `tenant_id` + `project_id` columns and indexes to all scoped tables |
| `server/system/types/` new type files | New — `tenant.go`, `tenant_membership.go` (see tenant definition spec) |
| `server/compose/types/namespace.go` | Modify — add `project_id` FK |
| `server/automation/types/workflow.go` | Modify — add `tenant_id` + `project_id` |
