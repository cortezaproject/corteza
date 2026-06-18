# DML (Data Migration Layer) — Implementation Plan

> Audience: implementing agent. This is prescriptive. Follow phases in order. Build + verify after each phase. Do not skip Phase 1 — it is the critical path.

## What we are building

A **DML** feature that lets a user:
1. Define a **connection** to an external RDBMS (postgres / mysql / sqlite).
2. **Extract that DB's schema** through the DAL → a set of models.
3. Review/edit a **mapping** (source table → compose module, source column → module field).
4. **Apply** the mapping → create compose namespace + modules.
5. **Import** the source rows → compose records.

UX is **stage-then-apply**: extraction proposes a mapping, user edits/approves, then apply + import run as separate steps.

## Locked design decisions

- DML defines **its own types** (`DmlConnection`, `DmlModel`, `DmlAttribute`, `DmlMapping`, `DmlImportRun`). These **translate to/from `dal.*`** via a conversion layer — DML does not reuse `system/types.DalConnection` as its own shape, but it does drive the DAL service to open the live connection.
- v1 introspection: **RDBMS only** (postgres / mysql / sqlite). Use the existing `information_schema` readers.
- Target of creation: **compose modules** (with generated configs) + **compose records** imported into them.
- **Stage-then-apply**, never one-shot.

## Current state (verified)

- Only file in the package: `server/system/service/dml/connection.go`. It is a fixture-backed stub that references **undefined** `types.DmlConnection` / `DmlModel` / `DmlAttribute` / `DmlConnectionFilter`. **The package does not compile.** Nothing imports it.
- `git` history has **no** prior `conv.go` / `dml.go` / `run_store.go` / `dml_migration.go` — greenfield, do not try to restore them.
- Verify before starting:
  ```
  cd server && go build ./system/service/dml/   # currently fails: undefined types.Dml*
  ```

## Two gaps that block everything (Phase 1)

1. **Schema introspection is a stub.** `server/store/adapters/rdbms/dal/connection.go`:
   ```go
   func (c *connection) Models(ctx context.Context) (dal.ModelSet, error) {
       return nil, nil   // <-- must be implemented
   }
   ```
2. **No `ddl.Table → dal.Model` converter.** Only the forward direction exists (`dal.Model`→`ddl.Table`). Must build the inverse, including **SQL-type-name → `dal.Type`** mapping.

---

## Phase 0 — Types + unbreak the package

Goal: package compiles; types exist; fixtures replaced by real (still-empty) wiring.

### 0.1 Create `server/system/types/dml.go`

Define (mirror the field names already used in `connection.go` fixtures so the existing file keeps working):

```go
package types

type (
    DmlConnection struct {
        ID           uint64   `json:"connectionID,string"`
        DalConnectionID uint64 `json:"dalConnectionID,string"` // links to the live DAL connection
        Handle       string   `json:"handle"`
        Type         string   `json:"type"`
        Label        string   `json:"label"`
        Driver       string   `json:"driver"`        // rdbms only in v1
        Capabilities []string `json:"capabilities"`
        ModelIdent   string   `json:"modelIdent,omitempty"`
    }

    DmlConnectionFilter struct {
        ConnectionID []string `json:"connectionID"`
        Handle       string   `json:"handle"`
        Type         string   `json:"type"`
    }

    DmlModel struct {
        ConnectionID uint64          `json:"connectionID,string"`
        Ident        string          `json:"ident"`       // source table name
        Label        string          `json:"label"`
        ResourceType string          `json:"resourceType"`
        Attributes   []*DmlAttribute `json:"attributes"`
    }

    DmlAttribute struct {
        Ident      string `json:"ident"`        // source column name
        Label      string `json:"label"`
        PrimaryKey bool   `json:"primaryKey"`
        Sortable   bool   `json:"sortable"`
        Filterable bool   `json:"filterable"`
        Type       string `json:"type"`         // corteza::dal:attribute-type:* (string form)
        Store      string `json:"store"`        // corteza::dal:attribute-codec:* (string form)
    }
)
```

Add **new** types for mapping + import (own file `server/system/types/dml_mapping.go`):

```go
type (
    DmlMapping struct {
        ID           uint64           `json:"mappingID,string"`
        ConnectionID uint64           `json:"connectionID,string"`
        NamespaceHandle string        `json:"namespaceHandle"` // target compose namespace
        Tables       []*DmlTableMap   `json:"tables"`
    }

    DmlTableMap struct {
        SourceIdent  string         `json:"sourceIdent"`   // source table
        ModuleHandle string         `json:"moduleHandle"`  // target module handle
        ModuleName   string         `json:"moduleName"`
        Skip         bool           `json:"skip"`
        Identifier   string         `json:"identifier"`    // source col used to dedup on import
        Columns      []*DmlColumnMap `json:"columns"`
    }

    DmlColumnMap struct {
        SourceIdent string `json:"sourceIdent"` // source column
        FieldName   string `json:"fieldName"`   // target module field
        FieldKind   string `json:"fieldKind"`   // compose field kind: String/Number/DateTime/Bool/Record/...
        Skip        bool   `json:"skip"`
    }

    DmlImportRun struct {
        ID           uint64 `json:"runID,string"`
        ConnectionID uint64 `json:"connectionID,string"`
        MappingID    uint64 `json:"mappingID,string"`
        Status       string `json:"status"`   // pending|running|completed|failed
        Processed    uint64 `json:"processed"`
        Failed       uint64 `json:"failed"`
        Error        string `json:"error,omitempty"`
        // resume cursor per table, opaque
        Cursor       map[string]string `json:"cursor,omitempty"`
    }
)
```

### 0.2 Rework `service/dml/connection.go`

- Remove the `fixtureConnections` / `fixtureModels` once Phase 1 lands; for Phase 0 keep them so it compiles, but move the type defs out (now in `system/types`).
- Keep the `dalReader` interface but **widen it** to what Phase 1 needs:
  ```go
  dalReader interface {
      GetConnectionByID(connectionID uint64) *dal.ConnectionWrap
      SearchModels(ctx context.Context) (dal.ModelSet, error)
      // NEW (Phase 1): introspect an external connection's live schema
      SearchExternalModels(ctx context.Context, connectionID uint64) (dal.ModelSet, error)
  }
  ```

### 0.3 Verify
```
cd server && go build ./system/service/dml/   # must pass now
```

---

## Phase 1 — Schema extraction (CRITICAL PATH)

Goal: given a live DAL connection ID, return real `dal.ModelSet` read from the external DB.

### 1.1 Implement table enumeration on the driver

`TableLookup(ctx, table, schema, dbname)` returns a single `*ddl.Table`. You need to **list all tables first**.

- Add to each driver's `informationSchema` (postgres/mysql/sqlite) a method:
  ```go
  func (i *informationSchema) TableSet(ctx context.Context, schema, dbname string) ([]*ddl.Table, error)
  ```
  Implementation: reuse `columnSelect()` / `scanColumns()` but **without** the per-table WHERE filter (i.e. select all columns for all base tables in the schema, group by table — `scanColumns` already groups into `[]*ddl.Table`). For postgres filter `table_schema = current_schema()` and exclude system schemas; for sqlite use `sqlite_master`.
- Wire `IndexLookup` per returned table so PK detection works (see 1.3).

### 1.2 Implement `(c *connection) Models(ctx)`

In `server/store/adapters/rdbms/dal/connection.go`, replace the `nil,nil` stub:

```go
func (c *connection) Models(ctx context.Context) (dal.ModelSet, error) {
    tables, err := c.informationSchema().TableSet(ctx, c.schema(), c.dbName())
    if err != nil { return nil, err }

    out := make(dal.ModelSet, 0, len(tables))
    for _, t := range tables {
        m, err := tableToModel(c.cw.ID, t)   // 1.3
        if err != nil { return nil, err }
        out = append(out, m)
    }
    return out, nil
}
```
- Find how `c` reaches the driver's `informationSchema` — look at how `assertAlterations*` (same file, ~line 415) already obtains the definer/IS reader and follow the same pattern. Do **not** invent a new connection field if one exists.

### 1.3 Build `ddl.Table → dal.Model` + the inverse type map

New file `server/store/adapters/rdbms/dal/introspect.go` (or driver-level if type names differ — see note):

```go
func tableToModel(connectionID uint64, t *ddl.Table) (*dal.Model, error) {
    m := &dal.Model{
        ConnectionID: connectionID,
        Ident:        t.Ident,
        Label:        t.Ident,
    }
    pk := primaryKeyColumns(t)   // from t.Indexes (IndexLookup), fallback: column named "id"
    for _, col := range t.Columns {
        typ, err := sqlTypeToDalType(col.Type)   // inverse map below
        if err != nil { return nil, err }
        m.Attributes = append(m.Attributes, &dal.Attribute{
            Ident:      col.Ident,
            Label:      col.Ident,
            PrimaryKey: pk[col.Ident],
            Sortable:   pk[col.Ident],
            Type:       typ,
            Store:      &dal.CodecPlain{},   // external columns are plain
        })
    }
    return m, nil
}
```

**Inverse type map** `sqlTypeToDalType(*ddl.ColumnType) (dal.Type, error)`. The `information_schema` reader normalizes names (see `scanColumns`: `timestamptz`, `timestamp`, `varchar(n)`, `numeric(p,s)`, etc.). Lowercase `col.Type.Name`, strip the `(...)` suffix, then map. This **inverts** the postgres forward map in `drivers/postgres/dialect.go`:

| SQL type (normalized, prefix) | dal.Type | Notes |
|---|---|---|
| `numeric`, `decimal`, `int`, `int2/4/8`, `bigint`, `smallint`, `real`, `double`, `float` | `&dal.TypeNumber{}` | parse `(p,s)` into Precision/Scale when present |
| `varchar`, `text`, `char`, `character`, `bpchar` | `&dal.TypeText{}` | parse `(n)` into Length |
| `timestamp`, `timestamptz` | `&dal.TypeTimestamp{Timezone: hasTZ, Nullable: col.Null}` | |
| `time`, `timetz` | `&dal.TypeTime{Timezone: hasTZ}` | |
| `date` | `&dal.TypeDate{}` | |
| `bool`, `boolean` | `&dal.TypeBoolean{}` | |
| `json`, `jsonb` | `&dal.TypeJSON{}` | |
| `uuid` | `&dal.TypeUUID{}` | |
| `bytea`, `blob` | `&dal.TypeBlob{}` | |
| anything else | `&dal.TypeText{}` + log warning | safe fallback; never error out the whole import |

Set `Nullable` from `col.Type.Null` on each produced type that supports it.

> Note on driver differences: mysql/sqlite return slightly different type names (`tinyint(1)`→bool, `datetime`, `double`, `integer`). Keep the map permissive (prefix match) so all three drivers funnel through it. Put the map in a shared spot reachable by `rdbms/dal`.

### 1.4 Expose via DAL service + DML service

- Add `SearchExternalModels(ctx, connectionID)` to `dal.FullService` (`server/pkg/dal/service.go`, interface ~line 34 + impl): look up `ConnectionWrap` by ID, call `cw.connection.Models(ctx)`, return.
- Implement `service/dml.Connection.FindModels` / `FindModelByIdent` to call `r.dal.SearchExternalModels(...)` and convert `dal.Model` → `types.DmlModel` (this is the **conversion layer** — create `service/dml/conv.go`: `fromDalModel`, `fromDalAttribute`, mapping `dal.Type.Type()` string → the `corteza::dal:attribute-type:*` constant). Drop the fixtures.

### 1.5 Verify
- Unit test `sqlTypeToDalType` for every row above (table-driven).
- Integration: point at a throwaway postgres with 2-3 tables, assert `Models(ctx)` returns expected idents + attribute types.
```
cd server && go test ./store/adapters/rdbms/... ./pkg/dal/... ./system/service/dml/...
```

---

## Phase 2 — Mapping (staged)

Goal: auto-generate an editable `DmlMapping` from extracted models; persist it.

- New service `service/dml/mapping.go`:
  - `GenerateMapping(ctx, connectionID) (*types.DmlMapping, error)` — call `FindModels`, then per model build a `DmlTableMap` (default `ModuleHandle` = sanitized table ident via `handle.Cast`; per column build `DmlColumnMap` with `FieldKind` derived from the dal type — see kind map below).
  - `CRUD` on `DmlMapping` (Create/Update/FindByID/DeleteByID). Persist via `store.Storer` — add a store type + migrations following an existing simple resource (copy the pattern from `system` resource stores, e.g. how `DalConnection` is stored).
- **dal.Type → compose field kind** map (for `DmlColumnMap.FieldKind`):
  | dal.Type | compose field kind |
  |---|---|
  | TypeNumber | `Number` |
  | TypeText / TypeEnum | `String` |
  | TypeTimestamp / TypeDate / TypeTime | `DateTime` |
  | TypeBoolean | `Bool` |
  | TypeRef | `Record` |
  | TypeJSON / TypeBlob / TypeUUID / other | `String` |
- Validation before apply: handles unique within mapping, valid handle format, at least one non-skipped column per non-skipped table, identifier column exists.

---

## Phase 3 — Apply (create compose modules)

Goal: mapping → real compose namespace + modules.

- New `service/dml/apply.go`: `Apply(ctx, mappingID) error`:
  1. Ensure namespace: find by `NamespaceHandle`, else `composeNamespaceSvc.Create`.
  2. Per non-skipped `DmlTableMap`: build a `*compose/types.Module`:
     - `Handle = ModuleHandle`, `Name = ModuleName`, `NamespaceID = ns.ID`.
     - `Fields`: per non-skipped `DmlColumnMap` → `*compose/types.ModuleField{Name, Label, Kind: FieldKind}`.
     - Call `composeModuleSvc.Create(ctx, mod)` — this **auto-runs** `DalModelReplace` → DDL creates the internal table (see `compose/service/module.go:318` and `ModuleToModel:1246`). Do not call DAL directly.
  3. Idempotency: if a module with that handle already exists in the namespace, update instead of create (or skip). Optionally route through envoy (`pkg/envoy/store/compose_module_marshal.go`) for upsert + ref resolution — prefer direct service calls for v1 simplicity.
- Record the created module IDs back onto the mapping (or a side table) so import knows the target.

---

## Phase 4 — Import run

Goal: copy source rows → compose records.

- New `service/dml/import.go`: `RunImport(ctx, mappingID) (*types.DmlImportRun, error)`:
  1. Create `DmlImportRun{Status: "running"}`, persist.
  2. Open source connection: `dalSvc.GetConnectionByID(connectionID)`.
  3. Per non-skipped table: build the **source** `dal.Model` (reuse the extracted one) and `Search(ctx, model, filter)` to get a `dal.Iterator`. Page with the iterator; respect `Cursor` for resume.
  4. Per source row: map columns→field values per `DmlColumnMap`, construct `*compose/types.Record{Values: ...}`, accumulate into batches.
  5. Write batches via `composeRecordSvc.Bulk(ctx, skipFailed, ops...)` (`compose/service/record.go:601`). Use `Identifier` column for dedup (lookup-then-update vs create) if set.
  6. Update `Processed`/`Failed`/`Cursor` as you go; set `Status: completed|failed` at end. Persist progress so a poll endpoint can read it.
- Keep batches bounded (e.g. 500 rows). On batch error with `skipFailed=true`, increment `Failed`, continue.

---

## Phase 5 — API + wiring

- REST endpoints (follow existing system API gen patterns — `server/system/rest/` + `*.yaml` definitions, then `make` codegen):
  - `GET  /dml/connections` , `GET /dml/connections/{id}` — list/show DML connections.
  - `GET  /dml/connections/{id}/models` — **extract** schema (Phase 1).
  - `POST /dml/connections/{id}/mapping` (generate) , `GET/PUT /dml/mappings/{id}` — mapping CRUD (Phase 2).
  - `POST /dml/mappings/{id}/apply` — create modules (Phase 3).
  - `POST /dml/mappings/{id}/import` , `GET /dml/runs/{id}` — start/poll import (Phase 4).
- Wire the `dml.Connection` (+ new services) into app boot / service registry — currently the package is **orphaned** (imported by nothing). Find where `system/service` services are constructed and registered, add DML there, inject `dal.FullService`, `store.Storer`, compose module/namespace/record services.

---

## Gotchas / rules

- **Never** call DAL DDL directly to create internal tables — go through `composeModuleSvc.Create` so module + dal.Model + DDL stay consistent.
- `Models()` must stay tolerant: unknown SQL types → `TypeText` fallback + warn, never hard-fail a whole table.
- The existing `Models()` stub comment says it returns nil to avoid breaking the model-add procedure. Your new impl runs only on the **external/source** connection during extraction — make sure you do not call it on the primary connection's normal model-add flow. Gate via the new `SearchExternalModels` entrypoint, not by changing how the primary connection uses `Models()` internally. If the primary path also calls `Models()`, add a guard/flag so introspection only runs when explicitly requested.
- Build after every phase: `cd server && go build ./... && go vet ./system/service/dml/...`.
- Reuse `handle.Cast` / existing handle validators for module & namespace handles.

## File checklist

Create:
- `server/system/types/dml.go`, `dml_mapping.go`
- `server/system/service/dml/conv.go`, `mapping.go`, `apply.go`, `import.go`
- `server/store/adapters/rdbms/dal/introspect.go` (+ `TableSet` on each driver's `information_schema.go`)
- store + migrations for `DmlMapping`, `DmlImportRun`
- REST handlers + API yaml + codegen

Edit:
- `server/store/adapters/rdbms/dal/connection.go` — real `Models()`
- `server/pkg/dal/service.go` — add `SearchExternalModels`
- `server/system/service/dml/connection.go` — drop fixtures, real DAL calls
- app boot / service registry — wire DML services
