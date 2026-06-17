# Handoff: Extract Action Log into a DAL-style Sub-Service

**Repo:** `corteza` (Go server) — branch `feature-projects`
**Status:** Design proposal, not started. No code written yet.
**Related:** `OBSERVABILITY_HANDOFF.md` (lineage/causation-graph project — this sub-service is its natural feed point).

---

## 1. Goal

Wrap the **core of the action log** (storage + retrieval) into a central, pluggable **sub-service**, the same way `DAL` is a central sub-service for data access.

Key constraint: **Corteza still generates the action data.** We only extract *where actions go and how to query them back* — not how they're produced.

End state: one specialized thing owns action ingest → policy → storage (one or many backends) → retrieval → retention. Services keep emitting actions exactly as today, through an interface, unaware of the backend.

---

## 2. Action Log Today (current architecture)

Audit trail. Every service mutation builds a structured `Action` and records it to the `actionlog` RDBMS table. Codegen-heavy, policy-gated, context-enriched.

| Area | Key file | Key types / funcs |
|---|---|---|
| Core types | `server/pkg/actionlog/types.go` | `Action`, `Meta`, `Severity`, `Filter`, `ActionSet` |
| Recorder interface | `server/pkg/actionlog/service.go` | `Recorder` = `Record(ctx, *Action)` + `Find(ctx, Filter)` |
| Policy gating | `server/pkg/actionlog/policy.go`, `canned_policies.go` | `policyMatcher`, `NewPolicyAll/Any/Negate/Match*`, `MakeProductionPolicy/DebugPolicy/DisabledPolicy` |
| Context enrich | `server/pkg/actionlog/context.go`, `service.go` (`enrich`) | `RequestOriginToContext`, origin/reqID/actorID/IP/ts injection |
| Codegen (producer) | `server/pkg/codegen/actions.go` + `*/service/*_actions.yaml` → `*_actions.gen.go` | props structs, action/error constructors, `ToAction()`, `recordAction()` helper |
| Store / persistence | `server/store/interfaces.gen.go`, `store/adapters/rdbms/*.gen.go` | `Actionlogs` iface: `SearchActionlogs`, `CreateActionlog`, `TruncateActionlogs`, … ; `actionlog` table |
| REST exposure | `server/system/rest/actionlog.go` (+ handlers/request) | `Actionlog.List` (list-only; no update/delete over REST) |
| Automation exposure | `server/system/automation/actionlog_handler.go` | `actionlog.search` / `.each` / `.record` |
| Boot / wiring | `server/system/service/service.go` (~170–220) | `DefaultActionlog = actionlog.NewService(store, logger, tee, policy)` |

### Data flow today
```
service.Create()
  └─ recordAction(ctx, props, ActionFn, err)   # generated helper
       ├─ ActionFn(props).ToAction()           # build Action
       ├─ enrich(ctx): origin, reqID, actorID, IP, ts
       ├─ policy.Match(a)?  ── no ──> drop
       └─ store.CreateActionlog(ctx, a)         # detached ctx, survives request cancel
```

### Known gaps (motivation for this work)
- **Flat, not causal** — no span/parent linkage; can't trace call → workflow → downstream. (The lineage project adds that on top.)
- **No retention** — `TruncateActionlogs` exists but no scheduler uses it; table grows unbounded.
- **Single hardcoded sink** — only RDBMS; no eventbus/SIEM fan-out. (Flow "step 7 / eventbus" is aspirational, not wired.)

---

## 3. Reference Pattern: how DAL is structured (mirror this)

DAL is the proven in-repo template for a pluggable central sub-service.

| DAL piece | File | Role |
|---|---|---|
| `service` struct + `FullService` iface | `server/pkg/dal/service.go:15`, `:34` | central service, public contract |
| `New(log, inDev)` | `server/pkg/dal/service.go:86` | constructor |
| `Service()` / `SetGlobal()` | `server/pkg/dal/service.go:109`, `:117` | global singleton accessor + boot setter (panics if unset) |
| `Connection` iface | `server/pkg/dal/driver.go:25` | pluggable backend contract |
| `ConnectorFn`, `Driver` | `server/pkg/dal/driver.go:116`, `:129` | backend factory + capability metadata |
| `RegisterConnector` / `RegisterDriver` | `server/pkg/dal/driver.go:161`, `:167` | registries; adapters self-register in `init()` |
| Boot | `server/app/boot_dal.go`, `boot_levels.go:238` | initialized at `bootLevelStoreInitialized`, after primary store |
| Store bridge | `server/store/dal.go` | primary store wrapped as a DAL `Connection` (`ToDalConn`) |
| Config records | `server/system/types/dal_connection.go` | `types.DalConnection` persists which backends exist |

**Pattern in one line:** singleton service + `init()`-registered pluggable backends behind one interface + booted after store + backends optionally configured via persisted records.

---

## 4. Proposed Design

### 4.1 Split into two halves

| Half | Owns | Change |
|---|---|---|
| **Producer** (stays in Corteza core) | `*_actions.gen.go`, props, `ToAction()`, `recordAction()` — the *generated action data* | **none** |
| **Sub-service** (new central thing) | ingest, policy, enrich, fan-out, storage, retrieval, retention | extract from `pkg/actionlog`, add pluggability |

The producer already depends on the `Recorder` **interface**, not the store. That interface is the seam — most of the inversion is already done.

### 4.2 New abstraction: the Sink

Today there's one hardcoded RDBMS write. Invert it into a registry of **sinks** behind one interface, mirroring DAL's `Connection`/`ConnectorFn`:

```
Record(ctx, *Action)
  ├─ enrich + policy.Match
  └─ fan-out ──> [ rdbms sink ]    primary, queryable (wraps store.CreateActionlog/SearchActionlogs)
                 [ stdout sink ]   debug (today's "tee")
                 [ eventbus sink ] -> lineage / causation graph
                 [ external sink ] -> SIEM / Loki (optional)
Find(ctx, Filter) ──> the designated queryable sink (default: rdbms)
```

- `Sink` interface: at minimum `Write(ctx, ...*Action)`; queryable sinks also `Find(ctx, Filter)`.
- `RegisterSink(fn)` registry; each sink self-registers in `init()` (Postgres/REST-DAL style).
- Store becomes the **default sink** (wrap `store.CreateActionlog` / `SearchActionlogs`), exactly as DAL wraps the primary store.

### 4.3 DAL → actionlog mapping

| DAL | actionlog target | Exists today? |
|---|---|---|
| `dal.Service()` + `SetGlobal` | `actionlog.Service()` singleton | ✗ — promote `DefaultActionlog` to a package global accessor |
| `FullService` iface | `Recorder` (`Record`+`Find`) | ✓ extend |
| `Connection` backend | **`Sink`** iface | ✗ new (store currently hardcoded) |
| `ConnectorFn` + register fns | `RegisterSink(fn)` registry | ✗ new |
| `types.DalConnection` records | `ActionlogSink` config records | ✗ new, **optional** (can hardcode default sink set first) |
| boot at `bootLevelStoreInitialized` | same level, after store | ✗ wire |
| codegen → Models | codegen → Actions | ✓ unchanged |

### 4.4 Reuse for free
- `Recorder` interface — the seam; producer already uses it.
- `policy.go` matchers, `enrich()`, `Filter` / `ActionSet`.
- `store` stays as the default sink implementation (no schema change).

---

## 5. Payoffs

| Win | How |
|---|---|
| Retention | sub-service owns `Truncate` + a scheduled prune; fixes unbounded growth |
| Multi-sink | ship to external SIEM/Loki without touching any service |
| Lineage (`OBSERVABILITY_HANDOFF.md`) | eventbus sink = natural feed for the causation graph |
| Testability | swap in an in-memory sink |

---

## 6. Effort / Risk

**Low–medium.** No business-logic or codegen change — the `Recorder` seam absorbs the inversion. Producers keep calling `recordAction(...)` unchanged.

Risk concentrated in: (a) preserving the **detached-context** write semantics (logs must survive request cancellation), (b) keeping `Find`/REST query behavior identical, (c) boot ordering (must come up after primary store).

---

## 7. Open Decisions (resolve early)

1. **Package location** — extend `server/pkg/actionlog` in place, or split a thinner `pkg/actionlog/service` like DAL? (Lean: extend in place; smaller blast radius.)
2. **Config records now or later?** — start with a hardcoded default sink set (rdbms + stdout) and add `ActionlogSink` persisted config only when a second real backend lands.
3. **Fan-out semantics** — sync vs async per sink; failure policy (must the primary rdbms write block the response while a remote SIEM sink is best-effort/async?).
4. **Single queryable sink** vs merge across sinks for `Find` (recommend: one designated queryable sink = rdbms; others write-only).
5. **Global singleton vs DI** — DAL uses a global `Service()`. Match it for consistency, or thread via system service? (Lean: match DAL.)

---

## 8. Suggested First Steps

1. Read `server/pkg/actionlog/service.go` (current `NewService`, `Record`, `Find`, `enrich`) and `server/pkg/dal/driver.go` (registry pattern) side by side.
2. Define `Sink` interface + `RegisterSink` registry in `pkg/actionlog`.
3. Refactor current store write/search into a `rdbmsSink` implementing `Sink` (no behavior change — byte-for-byte same queries).
4. Make `NewService` take a sink set instead of a bare store; fan-out in `Record`, route `Find` to the queryable sink.
5. Add `actionlog.Service()` / `SetGlobal()`; wire at boot after store init (mirror `boot_dal.go`).
6. Verify: existing `GET /actionlog/` + automation `actionlog.search/record` behave identically; production/debug/disabled policies unchanged; detached-context write preserved.
7. Only then: add eventbus sink (for lineage) and a retention prune job.

---

## 9. Context for the new session

- Caveman output mode may be active in the user's session (terse, fragments, tables; code/docs written normally). Not relevant to the code itself.
- User prefers concept-level / TL;DR answers, table-formatted, no code unless asked.
- Today's date context: 2026-06-17.
- This doc is a **design handoff** — confirm the open decisions in §7 with the user before writing code.
