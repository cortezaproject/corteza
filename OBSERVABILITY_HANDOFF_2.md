# Observability / Lineage — Handoff 2 (design decisions)

**Status:** design refined, not started. No code written.
**Date:** 2026-06-17. Branch: `feature-projects` (corteza server, Go).
**Reads with:** [OBSERVABILITY_HANDOFF.md](OBSERVABILITY_HANDOFF.md) — the original (phases, current-state audit, smoking gun). This doc captures decisions made on top of it. Where they differ, **this doc wins**.

## What this session settled

The goal is unchanged: *for any specified operation, reconstruct the full history — from the originating action through everything in between.* This session locked down four things the original left open:

1. The three existing pillars and how they relate.
2. Where human/AI **approval notes** live (NOT in the action log — own entity).
3. The **final action-log type** (existing fields + lineage spine).
4. **Access control** for reading the reconstructed path.

## The three existing pillars (grounded)

| Pillar | Granularity | Structure | Scope | Durable? | Location |
|---|---|---|---|---|---|
| Logs (zap+zapfilter) | line | flat, ctx fields | global, ephemeral | no (unless shipped) | `pkg/logger/` |
| **Action log** | business touch | flat rows + Meta JSON | cross-component | yes (policy-gated) | `pkg/actionlog/` (codegen-renamed `ln`/`n` in gen files) |
| Automation sessions | workflow step (`Frame`) | **already a tree** | one session only | yes (if Trace on, or on error) | `automation/types/session.go`, `pkg/wfexec/session.go:72` |

**Key insight — session is a pre-existing in-workflow span tree.** `wfexec.Frame` already carries `ParentID`, `StepID`, `NextSteps`, `ElapsedTime`, `StepTime`, `Input/Scope/Results`. That's a span tree scoped to ONE session, disconnected at the boundary. Phase 4 should **bridge** these Frames into the global trace (Frame→node, `Frame.ParentID`→`parent_span_id`, session root links up via `causation_id`), not invent step spans from zero.

## Logs: carry the key, not the graph

Stamp **only** `trace_id` + `span_id` on log lines (handoff Phase 1). Do NOT put `parent_span_id`/`causation_id`/`entry_point` on logs.

- Logs are a side-effect stream for humans → not durable, lossy by log-level, no schema guarantee.
- Reconstruction must read from a guaranteed-present, indexed, stable-schema store. That's the action log, not logs.
- Rule: **logs hold the key to find the tree; the tree lives where it's guaranteed to be.** Pivot log→trace_id→action-log; never rebuild the tree from logs.

## Decision: approval notes get their OWN entity, NOT the action log

Question raised: where do "why did someone approve the AI-recommended sale" notes go?

**Verdict: a dedicated `approval` entity. The action log only records that the approval event happened and points to it.**

The note being immutable does NOT make the action log the right home. Immutability was one of five objections; the other four stand even when the note is write-once:

| Action-log property | Note's reality | Conflict (holds even if immutable) |
|---|---|---|
| append-only, write-once | (we make note write-once too) | — resolved |
| policy-gated persistence (can drop by severity/origin) | must ALWAYS persist | a required record can't sit in a droppable store |
| coarse single read permission | needs per-approval access control | one `actionlog.read` = all or nothing |
| never referenced by business logic | sale / compliance / UI point AT it | once business logic depends on a row it's a record, not a log |
| system-observed ("X happened") | human-authored intent ("here's my reasoning") | log observes actors, doesn't speak for them |

**The real benefit of a separate entity (not "safer" — usable):** an approval has identity (referenced), relations (sale→approval→note), queries (filter/report), permissions (sensitive), lifecycle (long retention). Those define a domain entity. The action log has none by design.

Litmus scenario: *"every AI sale Jane approved against low confidence, with her reasoning, last quarter."*
- Entity → one indexed query.
- Action-log Meta → full scan of every record-touch, JSON-dig, hope fields stamped consistently, hope policy didn't drop any, can't permission the result.

**Trust is a property you enforce, not a place you store:**
- store layer = **create + read only**, no Update / no Delete path.
- `prev_hash` chain (or signature) → tamper-evident, stronger than "we promise not to UPDATE."
- **Freeze the AI recommendation**: snapshot (model, prompt, output, confidence) + store its hash on the note. "Approved" must mean "approved *this exact* recommendation," provably — else an immutable note points at mutable content and is worthless.

### Reference direction (get this right or you break immutability)

| Direction | Verdict |
|---|---|
| **note → holds `trace_id` + action-log `span_id`** | **primary/authoritative** — note born knowing its event, written once |
| action log Meta → holds note ID | optional, **only if written in the same tx** as the event row (set-at-write, no later UPDATE) |

Never UPDATE an action-log row after the fact to add a note pointer — that mutates the ledger. Let the note point up. Shared `trace_id`/`span_id` joins both directions anyway.

### Approval entity shape (sketch)

| Field | Why |
|---|---|
| `id` | own PK |
| `trace_id`, `span_id` | stitch to lineage + the approval event |
| `actor_id`, `decided_at` | provenance |
| `decision` (approve/reject) | the act (typed/validated) |
| `note` (text) | the rationale — the "why" |
| `subject_ref` | what was approved (e.g. sale/record ID) |
| `recommendation_snapshot` + `snapshot_hash` | frozen AI output approved |
| `prev_hash` | tamper-evident chain |
| `created_by`, `created_at` | standard |

Layering:
```
AI recommendation  → immutable snapshot (hashed)
        ↑ approves
Approval entity    → write-once, hash-chained, trace-linked, holds snapshot hash
        ↑ pointer (note ID, same-tx only)
Action log row     → event "approved, trace Z, actor X"  ← stays lean walkable ledger
```

## The final action-log type

Existing fields + 5 lineage columns. Edit source [pkg/actionlog/types.go:23](server/pkg/actionlog/types.go#L23) (+ `types.yaml`), regen.

```go
Action struct {
	ID uint64 `json:"actionID,string"`

	// --- existing ---
	Timestamp     time.Time `json:"timestamp"`
	RequestOrigin string    `json:"requestOrigin"` // rest-api, cli, grpc, system
	RequestID     string    `json:"requestID"`
	ActorIPAddr   string    `json:"actorIPAddr"`
	ActorID       uint64    `json:"actorID,string"`
	Resource      string    `json:"resource"`
	Action        string    `json:"action"`
	Error         string    `json:"error"`
	Severity      Severity  `json:"severity"`
	Description   string    `json:"description"`
	Meta          Meta      `json:"meta"`

	// --- new: lineage spine ---
	TraceID      string `json:"traceID"`                // W3C trace-id (32 hex) — which operation
	SpanID       string `json:"spanID"`                 // W3C span-id (16 hex) — this node
	ParentSpanID string `json:"parentSpanID,omitempty"` // parent node — the edge (empty at root)
	CausationID  string `json:"causationID,omitempty"`  // cross-trace link (async re-entry)
	EntryPoint   string `json:"entryPoint,omitempty"`   // root label: rest:POST /records, scheduler:job:x, eventbus:evt, automation:trigger:42
}
```

| Field | Role | Empty when |
|---|---|---|
| `TraceID` | groups all nodes of one operation | never after rollout; old rows empty |
| `SpanID` | unique id of this node | never; **unique** |
| `ParentSpanID` | the edge → reconstruction | root only |
| `CausationID` | links async re-entry across traces | unless re-entry |
| `EntryPoint` | root label (meaningful on root row) | non-root rows |

Stamped in `enrich()` from trace ctx ([pkg/actionlog/service.go:118](server/pkg/actionlog/service.go#L118)). Keep workflow/trigger/session IDs, note pointer, snapshot hash in `Meta` unless queried often → then promote.

### DB schema + indexes (multi-driver: pg / mysql8 / mssql / sqlite)

```
trace_id        CHAR(32)      -- hex, no dashes
span_id         CHAR(16)      -- UNIQUE (enables recursive join)
parent_span_id  CHAR(16)      -- NULL at root
causation_id    CHAR(16)      -- NULL unless re-entry
entry_point     VARCHAR(255)
```
- Hex `CHAR` not `VARCHAR`/UUID → portable across all 4 drivers; swap to `BINARY(16)/(8)` later if hot (interface hides it).
- Indexes: `(trace_id)`, `(resource, timestamp)`, `(parent_span_id)`, `(causation_id)`; `UNIQUE(span_id)`.
- Allow NULL, stamp forward; old rows = no trace (acceptable). No partitioning now — plan time-range partition on hot lineage (ops-watch).
- Multi-driver gotchas: mssql no partial index; mysql8 `CHAR` needs binary-ish collation (`ascii_bin`) for clean hex compare; keep DDL in codegen so all 4 stay in sync.

### In-place vs separate lineage table
**In-place on `actionlog` now** (it's the hub, fewer joins). Split to a derived/separate table only when write volume forces partitioning.

## Reconstruction — full history of an event

One walk over the spine; everything else hangs off by `trace_id`/`span_id`.

```sql
-- have a span_id → find its trace
SELECT trace_id FROM actionlog WHERE span_id = ?;
-- pull whole operation flat (one indexed hit)
SELECT * FROM actionlog WHERE trace_id = ? ORDER BY timestamp;
```
Assemble in Go: build `parent_span_id → children` map; root = row with empty `parent_span_id` (carries `entry_point`); DFS by timestamp → ordered tree.

Hang off the walk:
| Add | Source | Join |
|---|---|---|
| async/resumed continuations | action log | `WHERE causation_id = ?` |
| value at each change | record revisions | `WHERE rel_resource=? AND trace_id=?` |
| the "why" (decisions) | approval entity | `WHERE span_id IN (...)` / `trace_id=?` |
| technical zoom (recent) | Tempo | same `trace_id` (skip if sampled/expired) |

`WITH RECURSIVE` in-DB walk = escape hatch (all 4 drivers support).

## Access control — reading the path

The full path is **inherently cross-cutting** (one op touches many resources/workflows/actors, maybe tenants). To see the whole path you must see nodes the user normally couldn't read → per-resource/entity RBAC can't gate it without leaving holes.

**Decision: add a dedicated `lineage.read` (a.k.a. trace-read) RBAC operation, separate from both:**
- raw `actionlog.read` (exists today: `CanReadActionlog` → op `actionlog.read` on system component, gates `GET /actionlog/`) = flat audit dump, no path.
- entity RBAC = scoped, can't span the path.

`lineage.read` = "walk the causation tree, see every hop" — its own trusted capability. Grant to the auditor/compliance group.

Two-tier if mixed trust:
| Op | Sees |
|---|---|
| `lineage.read` | full unredacted tree |
| `lineage.read.redacted` | full shape, sensitive node **content masked** (sees *that* it happened, not *what*) |

For *scoped operational* visibility (e.g. "approvals for sales in my project") use **entity RBAC on the approval entity** — natural on an entity, painful on the log. Likely you want both: `lineage.read` for auditors, entity perms for operational users.

## Build order (unchanged from original, with additions)

0→1→2 (parallel-ish) → 3 → 4 (big lift) → 5. From the original handoff. Additions this session slots in:

- **Phase 2**: also create the `approval` entity (.cue model + create/read-only store enforcement + hash-chain on insert). Stamp `trace_id` on record-revision rows.
- **Phase 4**: bridge `wfexec.Frame` tree into action-log nodes (don't reinvent).
- **Phase 5**: build the reconstruction endpoint + add `lineage.read` RBAC op gating it.

## Open questions (carried + new)

- Record-revision row: does it carry requestID/any correlation today? (Need it for exact value↔event join — verify.)
- Automation action-log Meta: what workflow/trigger/session identity does it capture now? (Verify before Phase 4.)
- Trace sampling strategy + Tempo retention target.
- Hash-chain scope: per-entity chain vs global? Signature vs plain `prev_hash`?
- `lineage.read.redacted` — needed at launch or later?
- Snapshot storage: inline JSON on approval vs separate immutable blob + hash?

## One-line summary

Spine = action log with `trace_id`/`span_id`/`parent_span_id`/`causation_id`/`entry_point`; reconstruct by walking it. Logs carry only the key. Approvals/notes are their OWN write-once, hash-chained, trace-linked entity (action log points to it, never embeds it). Reading the full cross-cutting path = a dedicated `lineage.read` permission, not raw audit perm. RDBMS adjacency list, no graph DB.
