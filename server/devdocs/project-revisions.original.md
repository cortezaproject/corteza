# Project Revisions, Deployment & Migrations

## Locked decisions

| # | Decision |
|---|---|
| A | Each revision gets its own Compose namespace + full resource clone |
| B | Deploying a revision soft-deletes the old namespace (and its resources) |
| C | Records live in the namespace — migration is always required on activation |

---

## Concepts

**Original project** — `RevisionOf = 0`, `RevisionNumber = 0`, status `active`.

**Revision** — a full clone of a project: new Project row + new Compose namespace + cloned
modules/fields/pages/page_layouts/charts/workflows/triggers/ng_automations/agents/chatbots/KBs.
Status starts `draft`. Only one draft revision per project at a time (enforced).

**Deployment** — activating a revision:
1. Diff revision schema vs parent schema → classify changes
2. Migrate records from old namespace into new namespace
3. Mark revision `active`, old project `deprecated`, old namespace soft-deleted

**Path A (dangerous)** — schema changes that require data transformation:
field type change, field deletion, module deletion (when records exist).
Record migration runs with per-field conversion/dropping logic.
Frontend warned; user must explicitly confirm.

**Path B (safe)** — additive or non-storage changes only:
new module/field, page/layout/chart/workflow/trigger changes, metadata.
Record migration is a straight copy (no transformation).
Frontend notified; auto-confirmed or single-click.

---

## Data model changes

### `Project` struct — new fields

```
RevisionOf     uint64  // 0 = original; points to parent project ID for revisions
RevisionNumber int     // 0 = original, 1+ = revision sequence
```

Status state machine:

```
draft ──deploy──▶ active ──(new rev activates)──▶ deprecated
```

### `ProjectDeploymentPlan` (new type, not stored)

```go
ProjectDeploymentPlan struct {
    Path             string          // "A" | "B"
    SafeChanges      []ProjectChange
    DangerousChanges []ProjectChange
    // populated only on path A:
    RecordMigrations []RecordMigrationSpec
}

ProjectChange struct {
    Kind       string  // see change-kind table below
    Handle     string  // resource handle (stable across revisions)
    Detail     string
}

RecordMigrationSpec struct {
    OldModuleID uint64
    NewModuleID uint64
    FieldMaps   []FieldMigration
}

FieldMigration struct {
    OldHandle string
    NewHandle string
    // nil = straight copy; non-nil = transform
    Transform *FieldTransform
}

FieldTransform struct {
    Kind string // "typecast" | "drop" | "default"
    // ...params
}
```

### Change kinds

| Kind | Path | Note |
|---|---|---|
| `module.added` | B | Additive |
| `module.deleted` | A (if has records), B (empty) | |
| `module.field.added` | B | Additive |
| `module.field.deleted` | A | Data loss |
| `module.field.typeChanged` | A | Migration needed |
| `module.field.metaChanged` | B | Non-structural |
| `page.any` | B | UI only |
| `pageLayout.any` | B | UI only |
| `chart.any` | B | UI only |
| `workflow.any` | B | Logic only |
| `trigger.any` | B | Logic only |
| `ngAutomation.any` | B | Logic only |
| `agent.any` | B | Config only |
| `chatbot.any` | B | Config only |

---

## Phase 1 — Model additions

Files to touch:

- `server/system/project.cue` — add `revision_of` + `revision_number` attributes
- `server/codegen/` — re-run to pick up new store/dal/types
- `server/system/types/project_deployment_plan.go` — new hand-written types
  (not stored, no codegen)

No new store table. Revisions are Project rows.

---

## Phase 2 — Revision creation (clone)

**API:** `POST /projects/{projectID}/revision`

**Guard:** only one draft revision per project (check `SearchProjects` for `RevisionOf=projectID, Status=draft`).

**Clone sequence:**

```
1. Create new Compose namespace (handle = parent.handle + "-rev-N")
2. Clone all modules:
   a. Create new Module row (new ID, new NamespaceID, same handle)
   b. Clone all ModuleFields (new IDs, new ModuleID)
   c. Build refRewriteMap: oldModuleID → newModuleID, oldFieldID → newFieldID
3. Clone pages (new IDs, new NamespaceID) — rewrite module refs via refRewriteMap
4. Clone page layouts (new IDs, new PageID from step 3)
5. Clone charts (new IDs, new NamespaceID) — rewrite module refs
6. Clone workflows + triggers + ng_automations (new IDs, new ProjectID)
   — rewrite any module/field ID refs in step args via refRewriteMap
7. Clone agents, chatbots, KBs (new IDs, new ProjectID)
8. Create new Project row:
   RevisionOf = projectID
   RevisionNumber = max(existing revisions) + 1
   Status = "draft"
   Config.NamespaceID = new namespace ID
```

**Reference rewriting** — resources that embed module/field IDs:

| Resource | Where refs live |
|---|---|
| Page | `ModuleID` on record-list/record-edit blocks |
| Chart | `Config.Reports[].moduleID` |
| Workflow steps | Step arguments (JSON) — field handles, module IDs |
| NgAutomation | Step arguments |

Rewrite via `resourceRefsExt` walk — already generates `ResourceRefs()`. Use same traversal
to find and swap IDs.

---

## Phase 3 — Deployment plan (diff)

**API:** `GET /projects/{projectID}/deployment-plan`

Reads only, no side effects.

**Algorithm:**

```
parent = project where id = revision.RevisionOf
parentModules = modules where projectID = parent.id   (keyed by handle)
revModules    = modules where projectID = revision.id  (keyed by handle)

for handle in union(parentModules, revModules):
  if missing in rev   → module.deleted  (A if has records, else B)
  if missing in parent → module.added   (B)
  else diff fields by handle:
    missing in rev field   → field.deleted (A)
    missing in parent field → field.added  (B)
    dal.Type differs        → field.typeChanged (A)
    meta only               → field.metaChanged (B)

non-module resources (pages/workflows/...) → always B
```

Plan path = A if any dangerous change exists, else B.

Build `RecordMigrationSpec` for every module pair (both path A and B — B is a straight copy,
A has transforms).

---

## Phase 4 — Activation

**API:** `POST /projects/{projectID}/activate`

Body:
```json
{ "confirm": true }
```

Frontend must call `deployment-plan` first, show summary, then POST activate.

**Sequence (both paths):**

```
1. Re-compute deployment plan (verify nothing changed since plan was fetched)
2. If path A and plan still has dangerous changes — require confirm=true
3. For each RecordMigrationSpec:
   a. Paginate records from old module (old namespace)
   b. Transform fields per FieldMigration (typecast / drop / default)
   c. Bulk-insert into new module (new namespace)
4. Update revision: Status = "active"
5. Update parent:   Status = "deprecated"
6. Soft-delete old namespace (and cascade: modules, pages, etc. under it)
```

**Path A extra step** — before record migration, run DDL on new namespace tables
for any type-changed or new fields. Reuse `dal_schema_alteration` service
(`CreateDalSchemaAlteration` + `ApplyDalSchemaAlteration`).

**Atomicity:** wrap steps 3–6 in a DB transaction where possible. DDL cannot be
transactional in all engines — on failure, surface per-module progress + rollback status.

---

## Phase 5 — REST additions

| Method | Path | Notes |
|---|---|---|
| `POST` | `/projects/{projectID}/revision` | Clone; 409 if draft already exists |
| `GET` | `/projects/{projectID}/revisions` | List all revisions (sorted by RevisionNumber) |
| `GET` | `/projects/{projectID}/deployment-plan` | Diff; read-only |
| `POST` | `/projects/{projectID}/activate` | Execute path A or B; requires confirm for path A |

Existing endpoints unchanged.

---

## Open questions

| # | Question | Leaning |
|---|---|---|
| OQ-1 | Record migration: blocking (sync) or async with progress polling? | Async for large datasets; sync for small. Need a threshold or always async. |
| OQ-2 | What triggers soft-delete cascade on old namespace — explicit cascade in code or DB trigger? | Explicit in service layer |
| OQ-3 | Workflow/trigger step args store module IDs as uint64 strings — can we always safely rewrite them, or do some encode IDs opaquely? | Need audit of step-arg encoding |
| OQ-4 | Can a revision be created from a `deprecated` project (to "un-deprecate")? | Probably yes — same clone logic |
| OQ-5 | Draft revision editing: can user modify revision resources directly (edit cloned modules), or must changes go through another mechanism? | Direct edit of revision resources (they're normal rows with new IDs) |
| OQ-6 | Record migration for deleted modules: drop records or archive them? | Drop (soft-delete with DeletedAt) |
| OQ-7 | Multiple namespaces per project in future — does `Config.NamespaceID` stay a single field? | Out of scope for now |
