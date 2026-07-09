# Project Revisions, Deployment & Migrations

## Locked decisions

| #   | Decision                                                        |
| --- | --------------------------------------------------------------- |
| A   | Each revision gets own Compose namespace + full resource clone  |
| B   | Deploying revision soft-deletes old namespace (and resources)   |
| C   | Records live in namespace — migration required on publish       |
| D   | Record migration sync (v1); internal switch enables async later |
| E   | Revision from `deprecated` project not allowed (v1)             |
| F   | Deleted module records soft-deleted (`DeletedAt`)               |
| G   | One namespace per project; multi-namespace out of scope         |

---

## Walkthrough

### State 0: Example project

- Front-end creates a new project
- All resources are bound to the given project
- All resource modifications occur directly on the underlying resource

```
Project (id 101)   ProjectID=0  ParentRevisionID=0  Revision=0  status=active
└─ Namespace "crm"
   └─ Module "contacts"   fields: name (text), age (text)
      ├─ Record { name: "Ada",   age: "37" }
      ├─ Record { name: "Alan",  age: "41" }
      └─ Record { name: "Grace", age: "?"  }
```

### Step 1: Revise project

`POST /projects/101/revision`

- Front-end requests a revision on the project
- Back-end produces a clone of the project and all underlying resources
- All modifications on this project revision occur directly on the cloned resources

Notes:

- Only active projects can be revised
- One revision per any parent

```
Project (id 101)   status=active            ← still live
└─ Namespace "crm"
   └─ Module "contacts" (name, age)  + records

Project (id 102)   ProjectID=101  ParentRevisionID=101  Revision=1  status=draft   ← NEW
└─ Namespace "crm-{revision handle}"
   └─ Module "contacts"   fields: name (text), age (text)   ← no records
```

### Step 2: Modify draft project

- All modifications on this project revision occur directly on the cloned resources
- All modifications to the actual project occur directly on the draft project

```
Project (id 102)   status=draft
└─ Namespace "crm-{revision handle}"
   └─ Module "contacts"   fields: name (text), age (NUMBER), phone (text)
```

### Step 3: Prepare for project publish

`GET /projects/102/publish`

- Read the state, diff changes, determine release strategy
- The endpoint determines if changes are safe or dangerous (terminology might change)

```
overall: dangerous
Changes:
  contacts.phone  added        → safe
  contacts.age    text→number  → dangerous
Suggested mapping:
  contacts: name→name (copy), age→age (cast text→number), phone→phone (default "")
```

### Step 4: Publish draft

`POST /projects/102/publish`

- Front-end requests publish
- defines project to publish
- provides all parameters for data migration

```json
{
  "confirm": true,
  "mappings": [
    {
      "module": "contacts",
      "fields": {
        "name": { "op": "copy" },
        "age": { "op": "cast", "to": "number", "onError": "null" },
        "phone": { "op": "default", "value": "" }
      }
    }
  ]
}
```

#### Flow

- Lock draft project
- re-calculate diff
- Validate request with project and diff
- Use DML to migrate records from old to new
- rename the old namespace to indicate an older revision
- soft-delete old namespace and resources
- rename new namespace to match the project
- mark the new project as active

```
Project (id 101)   status=DEPRECATED
└─ Namespace "crm-{version handle}"   soft-deleted (cascades to module + records)

Project (id 102)   Revision=1  status=ACTIVE
└─ Namespace "crm (revision 1)"
   └─ Module "contacts" (name, age: number, phone)
      ├─ Record { name: "Ada",   age: 37, phone: "" }
      ├─ Record { name: "Alan",  age: 41, phone: "" }
      └─ Record { name: "Grace", age: null, phone: "" }
```

---

## Reference

### `Project` — new fields

```
ProjectID        uint64  // root project ID; 0 = original; never changes across chain
ParentRevisionID uint64  // immediate parent; 0 = original
Revision         int     // 0 = original; increments per revision
```

```
draft ──publish──▶ active ──(next rev publishes)──▶ deprecated
```

### `ProjectDeploymentPlan` (transient)

Computed on `GET`, never stored. Publish recomputes independently and validates the submitted mapping against the fresh diff.

```go
ProjectDeploymentPlan struct {
    Path              string          // "safe" | "dangerous"
    Changes           []ProjectChange
    SuggestedMappings []ModuleMapping
}

ProjectChange struct {
    Path string            // e.g. "module.contacts.field.age"
    Risk ProjectChangeRisk // "safe" | "dangerous"
}
```

### Risk

Risk is presentational — the back-end publish path is identical for both.

| Risk        | Examples                                           | Front-end                       |
| ----------- | -------------------------------------------------- | ------------------------------- |
| `safe`      | new module/field, page/chart/workflow/metadata     | auto-confirm                    |
| `dangerous` | field type change, field deletion, module deletion | warn; explicit confirm required |

`confirm=true` is the only backend gate. Front-end sets it automatically for safe, after consent for dangerous.

### Risk classification (per resource)

```go
// generated
func (svc *<resource>Service) ClassifyChanges(changes []ProjectChange) ([]ProjectChangeRisk, error) {
    out := make([]ProjectChangeRisk, len(changes))
    for i := range out { out[i] = ProjectChangeRiskDangerous }
    return svc.classifyChanges(changes, out)
}

// hand-written companion — overrides entries known safe; rest stay dangerous
func (svc *<resource>Service) classifyChanges(changes []ProjectChange, out []ProjectChangeRisk) ([]ProjectChangeRisk, error) {
    return out, nil
}
```

### Embedded references rewritten on clone

Envoy handles ref rewrites automatically during clone.

| Resource    | Field                                        | Handled by envoy |
| ----------- | -------------------------------------------- | ---------------- |
| Page        | `ModuleID` on record-list/record-edit blocks | ✅ yes           |
| Chart       | `Config.Reports[].moduleID`                  | ✅ yes           |
| Workflow    | Step arguments (module, namespace, role, workflow refs) | ✅ yes |
| NgAutomation | Step arguments                              | ⚠️ no envoy support yet (`envoy: { omit: true }` in ng_automation.cue) |

### REST API

| Method | Path                              | Notes                                                                     |
| ------ | --------------------------------- | ------------------------------------------------------------------------- |
| `POST` | `/projects/{projectID}/revision`  | 409 if draft exists                                                       |
| `GET`  | `/projects/{projectID}/revisions` | List by `ProjectID`, sorted by `Revision`                                 |
| `GET`  | `/projects/{projectID}/publish`   | Diff + risk + suggested mapping; nothing stored                           |
| `POST` | `/projects/{projectID}/publish`   | Body: mapping + confirm; locks, validates, applies DDL + DML, flips state |

---

## Internal DML source (OQ-1 resolution)

DML today reads from external connections via `dal.Connection` interface
(`Search`, `Models`). The old namespace is already DAL-connected; we expose it
as a new connection type so the existing import machinery needs no structural
changes — only a new driver.

### New connection type: `corteza::dal:connection:internal`

`DmlConnectionParams.Type` is a polymorphic string. Add:

```
corteza::dal:connection:internal
params: { namespaceID: uint64 }
```

`connectInternal(namespaceID)` returns a `dal.Connection` implementation backed
by the internal record store. Registered alongside the existing `dsn` and `rest`
drivers in [pkg/dal/driver.go](../pkg/dal/driver.go).

### Internal `dal.Connection` implementation

Implement `dal.Connection` interface in `pkg/dal/internal_conn.go`:

| Method                       | Implementation                                                                             |
| ---------------------------- | ------------------------------------------------------------------------------------------ |
| `Models(ctx)`                | load modules for the namespace via compose module store; return as `dal.ModelSet`          |
| `Search(ctx, model, filter)` | read records via compose record store for that namespace + module; adapt to `dal.Iterator` |
| `Create/Update/Delete`       | return `ErrNotSupported` — source is read-only                                             |
| `Can(op)`                    | return true for Read, false for Write                                                      |

No new store tables. No new CUE schema. The connection row is short-lived: created
by publish, deleted after the import run completes (or fails).

### Publish flow changes

In `activateRevision` / publish handler:

1. Create `DmlConnection { type: internal, params: { namespaceID: oldNamespaceID } }` — ephemeral, stored only for the duration of the run.
2. For each module pair in the mapping: create `DmlMapping { connectionID: ^above, sourceIdent: oldModuleHandle, moduleHandle: newModuleHandle, namespaceHandle: newNamespaceHandle, columns: from publish body }`.
3. Run `Importer.RunImport()` per mapping (existing code path — unchanged).
4. On completion/rollback: delete the ephemeral `DmlConnection` + `DmlMapping` rows.

### What does NOT change

- `Importer.runImportWork` — unchanged; already calls `SearchExternalModels` + `SearchExternalData` by connection ID.
- `DmlMapping`, `DmlImportRun` schemas — unchanged.
- `DmlConnection` CUE schema — unchanged (new type is just a new string value).

---

## Open questions

| #   | Question | Leaning                                                                                 |
| --- | -------- | --------------------------------------------------------------------------------------- | ------------ | --- |
|     | OQ-2     | Workflow/trigger step args — module IDs always uint64 strings, or some encode opaquely? | Audit needed |     |
