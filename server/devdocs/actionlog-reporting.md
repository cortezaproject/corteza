# Action Log Reporting & Standalone Service

## Reporting — locked decisions

Aggregated (group-by + metric) queries over the action log, served at
`GET /actionlog/report`. Turns the flat event stream into counts sliced by
resource, action, actor, origin, severity or day — without pulling raw rows to
the client. **Built and integration-tested.**

| #   | Decision                                                                          |
| --- | --------------------------------------------------------------------------------- |
| A   | Report shape = dimensions (group-by) × metrics (aggregations)                     |
| B   | Dimension/metric keys live in registries in `pkg/actionlog` — single source       |
| C   | rdbms adapter maps registry *kinds* → goqu; dialect-portable, no per-key SQL       |
| D   | Guarded by `action-log.read` RBAC + scope (tenant membership + `CapRead`)          |
| E   | Both `from` and `to` timestamps required on the request                           |
| F   | Actor grouped by ID in SQL; human labels resolved separately, post-query          |
| G   | Read-only feature — no report definitions persisted                               |
| H   | Metrics always `float64` in the payload; normalized across drivers                |
| I   | Live-query only — results are never cached                                         |
| J   | `severity` exposed as its name (`info`…), mapped at the REST layer — not the raw int |

## Standalone service — locked decisions

The action log runs on its *own* `store.Storer` connection (own pool, optionally
its own database), built from a dedicated DSN, independent of the app DAL.
Reporting is the "efficient pull" extension on that dedicated connection.
**Prototyped in tests; not yet productionised.**

| #   | Decision                                                                                     |
| --- | -------------------------------------------------------------------------------------------- |
| S1  | actionlog runs on its **own `store.Storer`** connection (own pool)                            |
| S2  | Connection built from a dedicated DSN; **falls back to main `DB_DSN`** when unset             |
| S3  | **Keep the generated store layer** — actionlog reuses it, does not fork or remove it          |
| S4  | Service receives a DSN **or** a ready `store.Storer` (either constructs the same port)         |
| S5  | actionlog store ops are **pure SQL** (`s.Exec`/`s.Query` + goqu) — no DAL / DataDefiner        |
| S6  | Schema managed by existing `store.UpgradeActionlog` — actionlog table only                     |
| S7  | Reporting = **custom (non-generated)** `ActionlogReport` method on the rdbms `Store`           |
| S8  | Separate DB allowed: cross-DB joins avoided — actor labels resolved post-query via main store  |

---

## Current state

Fully wired end-to-end (REST → service → recorder → store), integration-tested
against the real store ([report_integration_test.go](../pkg/actionlog/report_integration_test.go)).

### Worked example — "actions & errors per day for compose records, last week"

One request traced end to end. Scenario: an admin on the Action Log dashboard
wants a per-day bar chart of how many actions ran against `corteza::compose:record`
and how many of them errored, over `2026-07-01 → 2026-07-08`.

That maps to `dimensions=[day]`, `metrics=[count, errors]`,
`resource=corteza::compose:record`, `from/to` = the week.

Compact path:

```
FE dashboard → SystemAPI.actionlogReport()
  → GET /actionlog/report   (dimensions/metrics/resource/from/to as query args — see curl below)
  ▼ system/rest/request/actionlog.go   Fill()     query → request.ActionlogReport
  ▼ system/rest/actionlog.go           Report()   request → actionlog.ReportRequest
  ▼ system/service/actionlog.go        Run()       scope + action-log.read, from+to, Normalize(), actor labels
  ▼ pkg/actionlog/service.go           Report()    defensive Normalize() → store
  ▼ store/.../custom_actionlog.go      ActionlogReport()   kind → goqu, group/order, scan
  ▲ ReportResult JSON → FE renders chart
```

#### Step 1 — build & send the request

- **Front-end.** Dashboard component collects the picker state (date range,
  resource, ticked metrics) and issues the request. The generated client
  ([system.ts:6019](../../lib/js/src/api-clients/system.ts#L6019)) rejects
  locally if `from`/`to` are missing (before any network call) and serialises the
  arrays as repeated `dimensions[]=`/`metrics[]=` params. The wire call:
  ```sh
  curl -G 'https://HOST/api/system/actionlog/report' \
    -H 'Authorization: Bearer $JWT' \
    --data-urlencode 'dimensions[]=day' \
    --data-urlencode 'metrics[]=count' \
    --data-urlencode 'metrics[]=errors' \
    --data-urlencode 'resource=corteza::compose:record' \
    --data-urlencode 'from=2026-07-01T00:00:00Z' \
    --data-urlencode 'to=2026-07-08T00:00:00Z'
  ```
- **Back-end.** Nothing yet — request in flight.

#### Step 2 — bind query params

- **Front-end.** n/a (already on the wire). Only contract to honour: array
  params carry the `[]` suffix.
- **Back-end.** `request.ActionlogReport.Fill`
  ([request/actionlog.go:317](../system/rest/request/actionlog.go#L317)) reads
  `dimensions[]` (falls back to `dimensions`), `metrics[]`, parses `from`/`to`
  with `ParseISODatePtrWithErr`, copies `resource`/`action`/`actorID`/`origin`/
  `limit`. Result: a typed `request.ActionlogReport`.

#### Step 3 — map to the domain request

- **Front-end.** n/a.
- **Back-end.** `Actionlog.Report`
  ([system/rest/actionlog.go:81](../system/rest/actionlog.go#L81)) copies the
  bound request into `actionlog.ReportRequest{Dimensions, Metrics, Filter{…}}`
  and calls `reportSvc.Run`. The REST layer holds no logic beyond this mapping.

#### Step 4 — access control + preprocessing (the standalone service)

- **Front-end.** Must handle the failure codes this step produces: `401/403`
  (scope or `action-log.read`), `400` (missing/inverted range, unknown
  dimension/metric — the error names the allowed keys). Surface them in the
  dashboard, don't render an empty chart.
- **Back-end.** `actionlogReport.Run`
  ([system/service/actionlog.go](../system/service/actionlog.go)) in order:
  `checkScope(ctx, CapRead)` (tenant membership + project capability) →
  `CanReadActionLog` RBAC → require both `from`/`to` and reject an inverted
  range → `rr.Normalize()` (dedupe keys, default `count`, reject unknown). This
  is the "preproc" half of the stand-alone service.

#### Step 5 — cross the port into pkg/actionlog

- **Front-end.** n/a.
- **Back-end.** `svc.actionlog.Report(ctx, *rr)` calls the `Recorder.Report`
  port ([pkg/actionlog/service.go:145](../pkg/actionlog/service.go#L145)), which
  runs a defensive `Normalize()` and delegates to `store.ActionlogReport`. pkg
  owns the *data pull*; it never does access control or label enrichment.

#### Step 6 — build & run the SQL

- **Front-end.** n/a.
- **Back-end.** `custom_actionlog.go ActionlogReport` turns registry *kinds*
  into goqu: `day` → dialect date-trunc of `ts`; `count` → `COUNT(*)`; `errors`
  → `SUM(CASE WHEN error != '' …)`; `WHERE resource = 'corteza::compose:record'
  AND ts BETWEEN from AND to`; `GROUP BY`/`ORDER BY` the day expr; scans into
  `ReportRowSet`. Roughly:
  ```sql
  SELECT ts::DATE AS day,
         COUNT(*) AS count,
         SUM(CASE WHEN error != '' THEN 1 ELSE 0 END) AS errors
    FROM actionlog
   WHERE resource = 'corteza::compose:record'
     AND ts BETWEEN '2026-07-01' AND '2026-07-08'
   GROUP BY ts::DATE
   ORDER BY ts::DATE;
  ```

#### Step 7 — postprocessing (actor labels)

- **Front-end.** When a report *does* group by `actor`, read `actorLabels`
  (`actorID → name/email`) to label the axis/legend; the `set` rows carry only
  the id. In this example there is no `actor` dimension, so `actorLabels` is
  absent.
- **Back-end.** `Run` calls `actorLabels`: scans rows for an `actor` dimension,
  and only then issues a second query to the **main** store (`SearchUsers`) to
  resolve ids → names. No `actor` dim here → skipped. Assembles `ReportResult`.

#### Step 8 — response & render

- **Front-end.** Receives the payload and maps `set` → chart series (x = `day`,
  y = `count`/`errors`):
  ```json
  {
    "dimensions": ["day"],
    "metrics": ["count", "errors"],
    "set": [
      { "dimensions": { "day": "2026-07-01" }, "metrics": { "count": 1200, "errors": 14 } },
      { "dimensions": { "day": "2026-07-02" }, "metrics": { "count": 980,  "errors": 3  } }
    ]
  }
  ```
  Metrics are always `float64`, so no per-driver number coercion on the client.
- **Back-end.** Marshals `ReportResult` to JSON and returns; the controller adds
  nothing further.

### Pieces

| Layer   | File                                                                     | Responsibility                                                      |
| ------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------ |
| Types   | [pkg/actionlog/report.go](../pkg/actionlog/report.go)                    | request/result types, dimension + metric registries, `Normalize()` |
| Port    | [pkg/actionlog/service.go](../pkg/actionlog/service.go)                  | `Recorder.Report`; delegates to `actionlogStore`                   |
| Adapter | [store/adapters/rdbms/custom_actionlog.go](../store/adapters/rdbms/custom_actionlog.go) | kind→SQL, query, row/value normalization             |
| Service | [system/service/actionlog.go](../system/service/actionlog.go)           | RBAC/scope, from+to validation, actor-label enrichment             |
| REST    | [system/rest/actionlog.go](../system/rest/actionlog.go)                  | query-param → `ReportRequest`                                       |

### Connection & wiring (today)

Reporting works, but the service is **not yet standalone**. All four component
inits build the recorder on the shared app store:

```go
DefaultActionlog = actionlog.NewService(DefaultStore, log, tee, policy)
```
in `compose`, `system`, `federation`, `automation` `service.go`. `DefaultStore`
is the one app-wide connection ([boot_levels.go:210](../app/boot_levels.go#L210)).

The standalone path already exists and is exercised only by the integration
tests ([service_integration_test.go:47](../pkg/actionlog/service_integration_test.go#L47)):

```go
s, _ := store.Connect(ctx, log, dsn, false)   // dedicated DSN → own Storer + pool
store.UpgradeActionlog(ctx, log, s)           // creates ONLY the actionlog table
svc := actionlog.NewService(s, …)             // recorder on the dedicated store
```

Everything needed for a standalone service is in place — it just isn't wired
into boot behind an option.

---

## Request & response schema

`GET /actionlog/report` — guarded by `action-log.read` RBAC + scope (tenant
membership + `CapRead`).

### Request (query params)

Arrays repeat with a `[]` suffix (`dimensions[]=day&dimensions[]=actor`); the
bare form is also accepted. `from`/`to` are ISO-8601 datetimes.

| Param          | Type       | Required | Notes                                                                    |
| -------------- | ---------- | -------- | ------------------------------------------------------------------------ |
| `dimensions[]` | `[]string` | no       | group-by keys: `resource` `action` `origin` `severity` `actor` `day`     |
| `metrics[]`    | `[]string` | no       | aggregations: `count` `actors` `errors`; defaults to `count` when omitted |
| `from`         | datetime   | **yes**  | inclusive lower bound on `ts`                                            |
| `to`           | datetime   | **yes**  | inclusive upper bound; must not precede `from`                          |
| `resource`     | `string`   | no       | exact resource-type match (e.g. `corteza::compose:record`)              |
| `action`       | `string`   | no       | exact action match                                                       |
| `actorID[]`    | `[]string` | no       | filter to one or more actor ids                                         |
| `origin`       | `string`   | no       | request origin (maps to `request_origin` column)                        |
| `limit`        | `uint`     | no       | caps returned **groups** (server `MaxLimit` 1000)                        |

Unknown `dimensions`/`metrics` keys are rejected with `400`; the error lists the
allowed set. Missing `from`/`to` or an inverted range also returns `400`.

### Response (`200`, `ReportResult`)

```json
{
  "dimensions": ["day"],
  "metrics": ["count", "errors"],
  "actorLabels": { "123": "Jane Doe" },
  "set": [
    { "dimensions": { "day": "2026-07-01" }, "metrics": { "count": 1200, "errors": 14 } },
    { "dimensions": { "day": "2026-07-02" }, "metrics": { "count": 980,  "errors": 3  } }
  ]
}
```

| Field             | Type                 | Notes                                                                        |
| ----------------- | -------------------- | ---------------------------------------------------------------------------- |
| `dimensions`      | `[]string`           | requested dimension keys, echoed in request order                            |
| `metrics`         | `[]string`           | effective metric keys after `Normalize` (includes the `count` default)       |
| `actorLabels`     | `map[string]string`  | `actorID → name/email`; **omitted** unless the `actor` dimension is requested |
| `set[]`           | `[]ReportRow`        | one row per group (capped at `limit`)                                        |
| `set[].dimensions`| `map[string]any`     | dimension key → group value; `actor` is a string id, `day` is `YYYY-MM-DD`   |
| `set[].metrics`   | `map[string]float64` | metric key → value; always `float64`, normalized across drivers              |

---

## Reference

### Dimensions (group-by)

Registry: `reportDimensions` in [report.go](../pkg/actionlog/report.go).

| Key        | Kind     | Column           | Notes                              |
| ---------- | -------- | ---------------- | ---------------------------------- |
| `resource` | column   | `resource`       | raw value                          |
| `action`   | column   | `action`         | raw value                          |
| `origin`   | column   | `request_origin` | raw value                          |
| `severity` | column   | `severity`       | raw int in SQL; mapped to name at REST (see gap) |
| `actor`    | ref      | `actor_id`       | serialized as string in payload    |
| `day`      | date     | `ts`             | timestamp truncated to `YYYY-MM-DD` |

### Metrics (aggregations)

Registry: `reportMetrics` in [report.go](../pkg/actionlog/report.go). Default is
`count` when none requested.

| Key      | Kind            | Column     | SQL                                         |
| -------- | --------------- | ---------- | ------------------------------------------- |
| `count`  | count           | —          | `COUNT(*)`                                  |
| `actors` | countDistinct   | `actor_id` | `COUNT(DISTINCT actor_id)`                  |
| `errors` | countNotEmpty   | `error`    | `SUM(CASE WHEN error != '' THEN 1 ELSE 0)`  |

### Kind → SQL (the portable seam)

Dimensions and metrics are declared as *kinds*; the rdbms adapter owns the only
SQL. Adding a key over an existing kind = one registry entry, no adapter change.
A new kind needs an expression builder in the adapter.

| Kind                                | Built in adapter                                          |
| ----------------------------------- | -------------------------------------------------------- |
| `ReportDimensionColumn` / `…Ref`    | `goqu.C(column)`                                          |
| `ReportDimensionDate`               | dialect date-trunc via `ql` `date` ref (pg/mysql/sqlite/mssql) |
| `ReportMetricCount`                 | `COUNT(*)`                                                |
| `ReportMetricCountDistinct`         | `COUNT(DISTINCT column)`                                  |
| `ReportMetricCountNotEmpty`         | `SUM(CASE WHEN column != '' …)`                           |

### Types

```go
ReportRequest struct {
    Dimensions []string
    Metrics    []string
    Filter            // reuses actionlog.Filter (from/to, actor, resource, action, origin, limit)
}

ReportRow struct {
    Dimensions map[string]any     `json:"dimensions"`
    Metrics    map[string]float64 `json:"metrics"`
}

ReportResult struct {
    Dimensions  []string          `json:"dimensions"`
    Metrics     []string          `json:"metrics"`
    ActorLabels map[string]string `json:"actorLabels,omitempty"` // actorID → name/email
    Set         ReportRowSet      `json:"set"`
}
```

### `Normalize()` contract

Runs before the query; lower layers trust the request afterwards.

- Dedupes dimension and metric keys (drops empties).
- Applies default metric `count` when none given.
- Rejects unknown keys with the available set in the error.

### REST API

| Method | Path                | Query params                                                                | Guard             |
| ------ | ------------------- | --------------------------------------------------------------------------- | ----------------- |
| `GET`  | `/actionlog`        | `from, to, beforeActionID, resource, action, actorID, origin, limit`        | `action-log.read` |
| `GET`  | `/actionlog/report` | `dimensions[], metrics[], from, to, resource, action, actorID, origin, limit` | `action-log.read` |

`Run()` additionally requires both `from` and `to`, rejects an inverted range,
and enriches the result with `ActorLabels` (resolved via `SearchUsers`, safe
because the endpoint is already behind `action-log.read`).

---

## Goal

### Reporting completeness

| Gap                     | State  | Note                                                                              |
| ----------------------- | ------ | --------------------------------------------------------------------------------- |
| `resourceID` filter     | open   | `Filter.Resource` matches by type string only; per-instance history needs an ID column (see `service.History`, Gap 1) |
| Additional dimensions   | open    | e.g. request origin buckets — registry-only additions                             |
| `severity` name mapping | decided | map raw int → severity name at the REST layer (decision J); not yet implemented   |
| Additional metrics      | open   | e.g. p50/p95 latency once timing is captured — needs a new kind + adapter builder |
| Front-end consumer      | open   | admin webapp dashboards/charts over the report endpoint                           |
| Pagination / row cap    | open   | only `limit` today; large group sets unbounded ordering cost                      |

### Structural goal — standalone service on its own connection

Give the action log its own `store.Storer` connection, built from a dedicated
DSN, so it can point at a separate pool or a separate database entirely. **Keep**
the generated store layer — it is reused, not forked. This is deliberately *not*
the "rip actionlog out of the codegen" direction: the generated CRUD + the custom
`ActionlogReport` extension are exactly what the dedicated connection runs.

Why this is cheap: the write path is already tx-decoupled (`Record` runs on
`context.Background()`, fire-and-forget), actionlog store ops are pure SQL (never
touch `s.DAL`/`DataDefiner`), and `store.UpgradeActionlog` already builds just the
actionlog table. The standalone flow is already proven in the integration tests.

#### Wiring steps

1. **Option** — add `ActionLog.DB.DSN` (env e.g. `ACTIONLOG_DB_DSN`). Empty →
   reuse the main `DB_DSN`, so default behaviour is unchanged (same database, but
   still its own pool if we connect separately, or literally the same `Storer` if
   we choose to share when unset).
2. **Boot** — in `InitStore` (or a new `InitActionlogStore`), when a dedicated DSN
   is set: `store.Connect(dsn)` → `store.UpgradeActionlog` → keep the handle.
3. **Inject** — pass that handle into the four `actionlog.NewService(…)` sites
   instead of `DefaultStore` (compose/system/federation/automation `service.go`).
   The service already takes the 3-method port, so this is a constructor-arg swap.
4. **Constructor variants (S4)** — support both "give me a DSN" and "give me a
   `Storer`". Keep `pkg/actionlog` store-package-free: put the DSN→Storer helper
   in boot (or a thin `pkg/actionlog/standalone` subpackage), so the core service
   keeps depending only on its own port.
5. **Pool tuning** — a dedicated connection defaults to 256 open / 32 idle
   ([config.go](../store/adapters/rdbms/config.go)); size it down for a
   write-mostly log workload.

#### Blockers

| Blocker | Detail | Fix |
| ------- | ------ | --- |
| ~~`UpgradeActionlog` drops the table first~~ | data loss if pointed at a populated DB | **fixed** — drop removed; `UpgradeActionlog` now only `createTablesFromModels` (create-if-not-exists), rows preserved |
| ~~Delta/OldState columns missed on a separate DB~~ | the two `addColumn` fixes ran only in the *main* `Upgrade` | **fixed** — extracted `actionlogFixes` (single source), run by both `Upgrade`'s `fixesPost` and `UpgradeActionlog` |
| Actor labels need the main store | report enrichment queries `users` via `DefaultStore` ([service/actionlog.go:96](../system/service/actionlog.go#L96)); with a separate actionlog DB there is no cross-DB join | already correct — labels are a second query to the main store, not a join. Keep the two stores distinct in the report service |
| ID allocation across processes | a **separate DB** is not the issue — snowflake ids stay globally unique on one allocator. Collisions arise only if actionlog runs in its **own process** with a second `id.Next()` reusing the same node id | in-process standalone store: share the allocator (safe). Separate process: give it a **distinct snowflake node id** at boot |
