# Multi-Tenancy — implementation status

Written 2026-08-05 from the code, the five design docs in this folder, and the
engineering handover gist. It records **what is built versus what was
specified**, because the difference is currently only visible as five failing
tests nobody had attributed.

Nothing here is a design proposal. The open questions at the end are for
whoever owns the tenancy work.

## What the design specifies

From `multi-tenancy-design.md` and the handover:

- `System → Tenant → Project`. A tenant is the top-level isolation boundary and
  **has its own DAL connection**; project is the secondary unit within it.
- `pkg/scope` carries `TenantID` + `ProjectID` and a capability set, set by
  middleware from the authenticated request. It is described as already
  implemented and not to be redesigned.
- **"Scope guards are unconditional in generated service and store code. Do not
  add back conditional `tenantScoped` branches — they were removed
  intentionally because silent skip is worse than always-enforce."**
- `multi-tenancy-resource-scoping.md` classifies compose resources as
  project-scoped, with **`Record` scoped by inheritance** — module → namespace →
  project — rather than by a column of its own.

## What is built

| Piece                                     | State                                                                                                                                        |
| ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `pkg/scope` — Scope, registry, middleware | present                                                                                                                                      |
| Scope guards in **generated** services    | present — `RequireCapability`, `RequireTenantMembership` in chart, module, namespace, page, page_layout, and 37 files under `system/service` |
| `rel_tenant` / `rel_project` columns      | present on ~30 tables (`store/adapters/rdbms/upgrade_tenancy.go`)                                                                            |
| Project values                            | **populated** — live namespaces carry a `projectID`                                                                                          |
| Tenant values                             | **never populated** — every live resource reports `tenantID` absent (0)                                                                      |
| Scope in `compose/service/record.go`      | **absent** — hand-written, and the only compose resource with no scope guard                                                                 |
| Scope in the DAL write path               | **absent** — `pkg/dal`, the rdbms adapters and `compose/dalutils` import `pkg/scope` nowhere                                                 |

The migration comment states the intent for the unpopulated case: _"Columns
default to 0, which the store scope guard treats as unset — existing data is
reachable until a tenant is provisioned."_ So 0 is a deliberate placeholder,
not a bug in itself.

## The gap that is visible today

Five tests in `compose/service` fail with:

```
NOT NULL constraint failed: compose_record.rel_tenant
```

`TestRecord_boolFieldPermissionIssueKBR`, `TestRecord_defValueFieldPermissionIssue`,
`TestRecord_refAccessControl`, `TestRecord_searchAccessControl`,
`TestRecord_contextualRolesAccessControl`.

They are **not wrong, and they are not stale**. The chain:

1. Records are written through the DAL (`compose/dalutils` → `dal.Create`), not
   through the generated rdbms store. That path has no scope awareness, so
   nothing ever sets `Record.TenantID`.
2. The model declares the attribute with `HasDefault: true, DefaultValue: 0`
   (`compose/model/models.gen.go`), but that default does not survive into a
   schema built from the models — the insert violates NOT NULL.
3. A **running** instance is unaffected because its column arrived through
   `upgrade_tenancy.go`, which adds it defaulted to 0. Verified: creating a
   record against the local dev server succeeds.

So the tests are the only thing exercising a from-models schema, and what they
report is that a record cannot be inserted into one.

**Unverified:** whether a genuinely fresh deployment builds its tables from the
models and therefore hits this. The test path suggests it would; nobody has
tried it. That is the single most useful thing to check next.

## Open questions for whoever owns this

1. **Should `compose_record` carry `rel_tenant` at all?**
   `multi-tenancy-resource-scoping.md` specifies records as scoped by
   inheritance, but `upgrade_tenancy.go` lists `composeModel.Record` among the
   tables that gained the column. One of the two is out of date.
2. **Where should the tenant be stamped on a record?** The handover says guards
   are unconditional in generated service and store code — records go through
   neither. If tenancy is meant to reach them, the DAL write path needs the
   scope, and it currently has no access to it.
3. **Does a fresh install work?** If tables come from models rather than the
   migration, record creation may fail on day one.
4. **Is `compose/service/record.go` missing its guard deliberately?** It is the
   only compose resource whose service does not check scope. Given "silent skip
   is worse than always-enforce", this looks like an omission rather than a
   decision — but it is hand-written, so it may predate the rule.

## Do not

- Do not make the five tests pass by setting `TenantID` in their fixtures. They
  are the only signal about the from-models path, and silencing them is how
  this went unnoticed for weeks.
- Do not add conditional `tenantScoped` branches. The handover is explicit that
  these were removed on purpose.
