# Observability / Lineage — Design Handoff

**Status:** design agreed, not started. No code written.
**Date:** 2026-06-16. Branch: `feature-projects` (corteza server, Go).

## End goal

Not "trace one request." Build a **durable, resource-indexed causation graph**. Must trivially answer:

> "What path and automations were involved in updating record R to its current value?"

…and the generic form: "show every change to resource R over all time, and for each, the full path + workflows that caused it, from the original entry point (request / background tick / event)."

Two distinct questions, two mechanisms:

| | Q1: one operation, full path across all workflows | Q2: one resource, every change + path, all time |
|---|---|---|
| Mechanism | trace-context propagation (one trace_id across hops) | durable resource-indexed lineage |
| Store | Tempo (recent) + durable tier | durable tier only (Tempo expires) |

## Current state (grounded, verified in code)

| Concern | State | Location |
|---|---|---|
| Logging | zap + zapfilter. prod=JSON, dev=console. Global default + ctx-based propagation. | `pkg/logger/logger.go:35`, `pkg/logger/context.go:14` |
| RequestID | chi `middleware.RequestID` → `requestID` log field + `X-Request-Id` header. **REST only.** Format = `host/counter`, NOT UUID. | `pkg/api/server/middleware.go:24`, `pkg/api/server/logger.go:19-25` |
| Action log | durable event ledger. codegen from `*_actions.yaml` (54 files) → `*_actions.gen.go`. RDBMS table `actionlog`. Fields: actor, IP, requestID, origin, resource, action, error, severity, Meta(JSON). REST `GET /actionlog/` perm-gated. | `pkg/actionlog/`, enrich at `service.go:118` |
| Action log policy | prod persists: origin=automation OR severity∉{Debug,Info}. create/update/delete = **Notice** → persisted ✅. search/lookup/report = Info → dropped. | `pkg/actionlog/canned_policies.go`, `compose/service/record_actions.gen.go:484+` |
| Record revisions | **EXISTS.** field-level value history per record, via DAL, searchable by record ID. This is the value-lineage layer. | `compose/service/record_revisions.go`, `pkg/revisions/` |
| OTEL tracing | SDK v1.40 + OTLP HTTP exporter **wired but DEAD** — nothing makes spans. Conditional on `OBSERVABILITY_OTEL_EXPORTER_OTLP_ENDPOINT`. | `app/boot_levels.go:391` |
| Observability bus | `NewOtelDispatcher` + `NewLangfuseDispatcher` — **agentic/AI only**, not general app. | `app/boot_levels.go:384-397` |
| Prometheus | `/metrics`, basic-auth, `HTTP_METRICS`. HTTP request metrics only (chi-prometheus). No custom app metrics. | `pkg/api/server/metrics.go` |
| Sentry | full, panic recovery widespread, needs DSN. | `pkg/sentry/`, `app/boot_levels.go:129` |
| Health | `/healthcheck` (Scheduler/Mail/Corredor). No /readyz split. | `pkg/api/server/handlers.go:218` |
| pprof | `/debug/pprof` if `HTTP_ENABLE_DEBUG_ROUTE`. | `handlers.go:186` |
| Runtime monitor | MemStats/GC/goroutines → **logs only**, not scrapeable. | `pkg/monitor/monitor.go` |
| DB instrumentation | sqlmw debug wrapper, **logs only**, timing not exported. | `store/adapters/rdbms/instrumentation/debug.go` |
| Absent | StatsD, Datadog, NewRelic, Elastic APM, W3C trace-context propagation. | — |

### The smoking gun (why chains fragment today)
Workflow execution rebuilds context from scratch — **trace context severed at the workflow boundary**:
- `automation/service/session.go:254` → `auth.SetIdentityToContext(context.Background(), i)`
- `automation/service/session.go:355` → `var execCtx = context.Background()`

eventbus: sync `WaitFor` keeps ctx; async `Dispatch` (`pkg/eventbus/eventbus.go:91`) hands request ctx into a goroutine (lifecycle mismatch). Scheduler tick = fresh root.

Net: request → event → workflow1 → event → workflow2 fragments into disconnected traces today.

## The model (two tiers)

| Tier | What | Retention | Sampling | Indexed by |
|---|---|---|---|---|
| **Durable lineage** = enriched action log | business-element touches as span-tree nodes | forever | **never sample** | resource, trace_id, parent |
| **Technical trace** = Tempo | full call tree incl DB/internal | days–weeks | sample OK | trace_id |

**Rule:** sample Tempo, NEVER the lineage tier. Durable answer survives even when the technical trace was sampled away. Keep lineage at business granularity (records/modules/workflows/automations), not every DB row.

### Paths reconstructed from the **action log** (the spine)
Each enriched row = one node. Edges live in new columns:
```
trace_id        ← which operation
span_id         ← this node
parent_span_id  ← its parent (the edge)
causation_id    ← cross-trace link (async re-entry)
entry_point     ← root label (rest:POST /records, scheduler:job, eventbus:evt, automation:trigger:id...)
```
Reconstruction = adjacency walk:
1. `WHERE trace_id=X` → all nodes
2. build parent→children map in Go
3. root = entry_point row → walk down = the tree

Resource view: `WHERE resource=R ORDER BY ts` → each change carries its trace_id → fetch + assemble.
Values at each step come from **record revisions** (join by record + trace_id).
Tempo holds its own tree (recent, technical), joined by trace_id — NOT rebuilt from action log.

## Storage decision: RDBMS, NOT graph DB

Data = forest of shallow trees (one per trace_id), not a general graph. Retrieval is "one trace" or "one resource" — filter then assemble. Graph DBs only win at arbitrary-path / global-pattern queries you don't have.

- Adjacency list (`parent_span_id`, `causation_id` columns) + indexes on `(trace_id)`, `(resource, ts)`, `(parent_span_id)`.
- Fetch flat by trace_id, build tree in Go (one trace = tens/hundreds rows).
- `WITH RECURSIVE` escape hatch if ever need in-DB walk (pg/mysql8/mssql/sqlite all support).
- Zero new infra — action log already a table, Corteza already multi-driver RDBMS.
- Reconsider graph/columnar only for future cross-resource impact-at-scale. Not launch.
- **Ops watch:** hot-resource write volume → plan retention/partitioning on lineage table.

## What to build (sequenced)

**Phase 0 — Foundation (cheap, unblocks all)**
- Activate OTEL TracerProvider on boot (already wired). Point at Tempo.
- W3C `traceparent` propagator + tracer helper in pkg.

**Phase 1 — Spans on sync path (low/med)**
- `otelhttp` middleware in BaseMiddleware → root span per REST req, parse incoming traceparent.
- Inject `trace_id`+`span_id` into request-scoped zap logger (`contextLogger`).
- Return `trace_id` in response header.
- `otelsql` wrap DB driver (replace/augment sqlmw debug wrapper).

**Phase 2 — Action log = durable span tree (med)**
- Add columns: `trace_id`, `span_id`, `parent_span_id`, `causation_id`, `entry_point`.
- `enrich()` stamps them from ctx (`pkg/actionlog/service.go:118`).
- Migration + store adapter + REST filter (query by trace_id / resource).
- Stamp `trace_id` onto record revision row too → exact value↔event join.

**Phase 3 — Entry-point coverage (med) — kill orphans**
- Start root span + tag `entry_point` at: scheduler tick, eventbus async Dispatch, gRPC in, CLI/provision.
- eventbus async: serialize `traceparent` INTO the event payload, restore in handler (don't lean on dying req ctx).

**Phase 4 — Async/workflow propagation (HIGH, the spine)**
- Workflow session: stop `context.Background()` (session.go:254,355) — inherit + continue trace from triggering event.
- Emit span per workflow step; write action-log node per touched element with parent_span_id.
- Add automation identity (workflowID/triggerID/sessionID) to action log Meta.
- corredor: forward `traceparent` over its gRPC.
- Suspend/resume: workflows live for days — persist lineage IDs ON the session row; each resume = new span linked via causation_id. Don't model with one long-lived span.

**Phase 5 — Payoff + cleanup (low)**
- Build "change path over resource" query (walk parent→child by trace_id; join revisions⨝actionlog).
- Cut redundant once spans land: LogRequest/LogResponse mw, sqlmw debug logging, Monitor pkg.
- Sampling policy: Tempo sampled, action-log lineage never.

**Dependency order:** 0→1→2 (parallel-ish) → 3 → 4 (big lift) → 5.

## Keep / cut

| Thing | Verdict |
|---|---|
| Action log | KEEP — make it the hub |
| Record revisions | KEEP — value lineage |
| RequestID | KEEP, upgrade to carry trace_id; "adding requestIDs" is redundant, they exist |
| Tempo/traces | KEEP — recent technical drill-down |
| LogRequest/Response mw | CUT/flag once otelhttp lands |
| sqlmw debug logging | CUT/flag once otelsql lands |
| Monitor pkg (runtime→logs) | CUT once Prom runtime metrics wired |
| Web console ring buffer | keep, orthogonal |

## Hard parts (ranked)
1. **Workflow engine instrumentation** — heaviest. Continue trace + emit span/node per step. corredor is a separate runtime over gRPC.
2. **Suspend/resume** — workflows live days; one open span can't model it → persist lineage IDs on session, link resumes by causation_id.
3. **Async Dispatch context** — serialize traceparent into event payload; request ctx dies.
4. **Entry-point coverage** — REST done; scheduler/eventbus/automation/gRPC/CLI/provision need roots.

## Open questions for next session
- Does record revision row already store requestID/any correlation? (Verify — needed for exact value↔event join. Checked: searchable by rel_resource=record ID, but correlation field not confirmed.)
- What does automation action-log Meta currently capture re: workflow/trigger/session identity? (Verify before Phase 4.)
- Trace sampling strategy + Tempo retention target.
- Lineage table: extend `actionlog` in place vs derived/separate table (volume/retention concerns).
- Policy: propagate-into-one-trace vs new-root-with-span-link for periodic/independent ticks.

## One-line summary
Have the durable store (action log) + value layer (revisions) + OTEL SDK. Missing spine = **trace-context propagation across async/workflow boundaries + the span skeleton (trace_id/parent_span_id/causation_id/entry_point) on every durable action-log row.** Reconstruct paths by walking the action log; Tempo = recent technical zoom joined by trace_id. RDBMS adjacency list, no graph DB.
