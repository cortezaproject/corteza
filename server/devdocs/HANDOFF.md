# Engineering Handoff — Corteza Server

---

## Platform context

### Quick Notes

- Component: compose, system, federation, automation (all larger parents of resources)
- Package management: packages are stored into the `/vendor/` folder; standard golang stuff; do make sure to `go mod vendor` after adding new packages.
- REST, Service, Store layer are strictly separated

---

## Codegen system

Two independent generators cover all mechanical code. Run both with `make codegen` (CUE layers) and `make codegen-legacy` (REST layer via pkg/codegen).

| Generator                           | Driven by                 | Output                                                                            | Command               |
| ----------------------------------- | ------------------------- | --------------------------------------------------------------------------------- | --------------------- |
| CUE codegen (`server/codegen/`)     | `.cue` files per resource | service layer, type structs, resource refs, store interfaces, rbac, envoy, locale | `make codegen`        |
| pkg/codegen (`server/cmd/codegen/`) | `rest.yaml` per component | REST handlers, request parsers, REST controllers                                  | `make codegen-legacy` |

Each resource has a `.cue` file declaring: model (attributes), RBAC ops, service config, type generation config. Loaders (`server.*.cue`) aggregate these and feed templates (`assets/templates/gocode/`).

### Generated vs. hand-written split

Generated files are named `<snake>.gen.go`. The companion `<snake>.go` holds everything the template can't express: non-standard logic, bespoke methods, hooks called by generated code. The split mirrors the Go rule that methods on a type can span files.

The invariant: generated file is always machine-regenerable from `.cue` + template. Companion is human-maintained. Never put generated code in companion or companion logic in generated file.

### Service layer

Generates the full CRUD service skeleton (Create/FindByID/Search/Update/DeleteByID/UndeleteByID) with action-log audit wrapping, RBAC access control, scope enforcement, and before/after hooks per operation. Non-standard methods use `customBodyOps` — the template generates the audit shell and delegates the body to an `on<Op>` handler in the companion.

**CAUTION**: All root service-layer functions must be declared in `.cue` files and must be generated.

**RBAC / access control.** `genAccessController: true` (default) generates a `<Resource>AccessController` interface from the enabled CRUD ops — `CanCreate<R>` (create), `CanSearch<R>s` (search), `CanRead<R>` (lookup/update/delete/undelete). Use `customAccessOps: [...]` to suppress the generated Can check for specific ops when access control is bespoke for that op.

**customFunctions.** For bespoke non-CRUD ops (e.g. Reorder, MarkAsRead, Impersonate) that still need the standard scaffold (scope check, AC guard, action-log), declare them in `customFunctions`. The template generates the public wrapper method; only the `on<Name>` handler stays in the companion. Each entry: `name` (Go ident), `cap` (read/write scope capability), optional `ac` (AC method name), `args`/`results` (Go types), `action` (action-log stem, defaults to title-cased name). Extra import paths for custom signatures go in `customFunctionImports`.

### REST controllers

Generates CRUD controller methods (list/create/read/update/delete/undelete) on existing controller structs. Non-plain param types (JSON blobs, nested structs) are skipped and filled by `beforeCreate`/`beforeUpdate` hooks in the companion. 35 controllers currently generating.

### Type structs

Generates the main resource struct from model attributes. A struct is only generatable when it contains exactly model attributes + optional labels. 45/58 structs currently generating. The remaining 13 have non-model fields (relational sets, computed fields, unexported runtime state) and must stay hand-written.

### Resource refs

Generates `ResourceRefs()` for the dependency graph. Covers 10 source resources. The API is path-less: `Ref{Resource, Label, Reason, Unresolved}` with methods `Kind()/ID()`.

---

## Scope / multitenancy

### What it is

`pkg/scope` is a per-request DI container carrying `TenantID` + `ProjectID` and a `Capabilities` set. Middleware sets the scope from the authenticated request context. Every generated service reads it before any RBAC check.

### Always-on enforcement

Scope guards are unconditional in generated service and store code. Do not add back conditional `tenantScoped` branches — they were removed intentionally because silent skip is worse than always-enforce.

---

## Data Migration Layer (DML)

### What it is

A new resource type (`DmlConnection` + `DmlMapping`) that lets users extract schema from an external RDBMS via the DAL connection pool, stage a table→module / column→field mapping, approve it, create compose modules, and import records.

### Status

Core implementation complete. Server-side phases (extraction, mapping, apply, import, API) all built and adversarially reviewed. Branch: `feature-projects`.

### Key design decisions

DML owns its own types — it does NOT reuse DAL connection shapes. It translates to/from DAL only at the extraction boundary. UX is stage-then-apply, not one-shot. RBAC is coarse for now (uses CanSearchDalConnections as a proxy).

### Deferred gaps

- PK-less external tables have unstable cursors on paged reads
- First-class RBAC resource (`corteza::system:dml-*`) pending separate codegen step
- `SearchExternalData` bypasses DAL op/sensitivity guards (gated at REST boundary instead)

---

## Chatbot / AI

### Two surfaces

| Surface                 | Audience  | Auth                     | Streaming |
| ----------------------- | --------- | ------------------------ | --------- |
| Internal agentic app    | Dev/admin | Full user session        | No (sync) |
| External chatbot widget | End users | Widget key + session JWT | Yes (SSE) |

### Handoff (AI-to-human)

Allows a widget conversation to be transferred to a human operator. Data model: `ChatbotSession` → `ChatbotSessionStep` → `ChatbotSessionHandoff` (status workflow). Server side (types, service, API endpoints) is complete.

What remains: widget SSE event handling + "Talk to human" button; operator console (new admin view: session list, real-time session view, message input); notifications to TargetRoles on handoff request.

Spec: `.claude/chatbot-handoff-spec.md`.

### Scenario automation hooks

Each chatbot scenario step can have `beforeAutomationID`/`afterAutomationID`. Server-side `executeStep()` orchestration designed but not implemented. Widget is currently thick (manages scenario state client-side); the design moves it to dumb (SSE-driven). Not yet started.

---

## Observability / Lineage

### Goal

Durable, resource-indexed causation graph answering "what path of automations set record R to its current value" — spanning multiple async hops and time. Not per-request tracing.

### Architecture

Two tiers: durable lineage (enriched action log, forever, resource-indexed) + technical trace (Tempo, recent, sampled). Joined by trace_id. Paths reconstructed by walking the action log as an adjacency list — not a graph DB. Approval/decision notes are a separate write-once entity (not action-log meta) with their own permissions and retention.

### Status

Design complete (lives in AI session memory; original handoff docs were deleted from repo). Action log improvements (History, DiffResourceState, integration tests) committed. Lineage instrumentation not started — the action log schema still has no trace_id/span_id/parent_span_id columns.

### Blockers

- Ph0: Activate OTEL SDK (currently dead in boot)
- Ph1: Add 5 columns to Action struct + migration
- Ph4 (hardest): workflow automation uses `context.Background()` at session entry points, severing the trace — needs explicit propagation

---

## Action log reporting

### Reporting endpoint

`GET /actionlog/report` with dimensions (resource, action, actor, origin, severity, day) and metrics (count, actors, errors). Backed by registry-driven GROUP BY. Complete and integration-tested. Full spec in `server/devdocs/actionlog-reporting.md`.

### Standalone service

Action log designed to run on its own `store.Storer` from a dedicated DSN (`ACTIONLOG_DB_DSN`). Write path is already tx-decoupled, store ops are pure SQL, `UpgradeActionlog` exists, blockers fixed (drop-first removed, delta/OldState migration folded in).

What remains: wire `ACTIONLOG_DB_DSN` option → boot connect → inject into 4 `actionlog.NewService()` call sites (one per component).

---

## Unmigrated resources

### Record service

Non-CRUD shape (bulk/patch). Stays hand-written. No path to codegen without redesigning the API.

### User service

Bespoke AC-vs-load order, validate-vs-AC order, lifecycle. Scaffold model applied (generated audit shell + `on<Op>` handlers in companion). Won't move to full generation.

### Attachment, settings, permissions controllers

Genuinely non-CRUD. No ops to generate. Stays hand-written.

### Module and page type structs

Relational/computed fields in the struct body make them ungeneratable. Would require modelling those fields as `store:false` attributes, which risks leaking into other gen layers.

---

## Gotchas

### Build-green ≠ behavior-preserved

Service codegen migration required two rounds of adversarial review to catch: wrong AC-vs-load ordering, validate-before-AC leaking errors to unauthorized callers, label update conditional vs unconditional semantics, undelete not resetting `DeletedBy`. Always verify data paths, not just compilation.

### Generate only what the template owns

If a service has non-standard logic in a method, use `customBodyOps` (scaffold the audit shell, delegate body to `on<Op>`). Adding template knobs hits diminishing returns fast.

### CUE optional field access

CUE does not short-circuit `&&`. Accessing `attr.dal.type` when `dal` is optional causes a compile error. Pattern: `[ if attr.dal != _|_ { attr.dal.type }, "" ][0]`.

### Workflow agent return values

When a workflow agent returns a file path string, do not do `SERVICES.find(s => s.file === c.file)` — the array element is the path string itself. End rollouts at build-green; run verify as a separate pass after.

### Scope enforcement is always-on

`pkg/scope` guards are unconditional in generated service and store. Do not add back conditional `tenantScoped` branches.

### DML must not touch the source DB schema

`RegisterModelCache` (cache-only, no DDL) must be used when registering external models for read. `CreateModel` triggers DDL on the target and must never be called on a source DB.

---

## Related repositories

### Architecture overview

```
                        Polar.sh (billing events)
                               │ webhook
                               ▼
                        human-billing ──── Elestio API (provision / suspend / terminate VMs)
                               │ outbound webhook (optional)
                               ▼
                        Corteza Compose (billing records)

Corteza App Store UI ──── REST API ────► human-appstore-api (in-memory catalog from GitHub)

human-sre ────► builds + deploys human-billing + human-appstore-api to Elestio host

human-docs      docs for @Human Connections product (separate, no Corteza coupling)
corteza-docs    docs for Corteza platform (cortezaproject/corteza-docs, versioned)
```

`human-billing` and `human-appstore-api` are standalone Go services. Corteza is the consumer/destination, not the host. `human-sre` is the single deploy entry point for both.

---

## human-billing

### What it is

Middleware connecting Polar.sh (subscription billing) and Elestio (managed VM hosting). Ingests Polar webhook events, manages instance lifecycle (provision → suspend → terminate), enforces grace-period/past-due logic.

### Corteza integration

Outbound webhook only — optionally posts billing events to a Corteza Compose module as records. No inbound dependency on Corteza.

### Stack

Go, Chi, PostgreSQL (sqlx), Svix webhook verification. Deployed via `human-sre`.

### Status

Feature-complete MVP. No active gaps.

---

## human-appstore-api

### What it is

In-memory catalog service indexing Connections (third-party integrations) and their Operations. Backend for the Corteza App Store UI. Reads integration definitions from a GitHub repo (markdown), parses into a typed catalog, exposes over REST. GitHub webhook triggers debounced cache refresh on push.

### Corteza integration

Standalone — Corteza App Store UI consumes its API. No inbound Corteza dependency.

### Stack

Go, Chi, no database (memory-only). API key auth on all `/v1/*` endpoints. Deployed via `human-sre` on Elestio port 2281.

### Status

Feature-complete. Single release commit, stable.

---

## human-sre

### What it is

Single entry point for deploying the Human Brain platform. No application code — just the glue (Docker Compose + Makefile) that ties `human-billing` and `human-appstore-api` together and ships them to production.

### Where it runs

Production is a single Elestio-managed VM. Both services run there as Docker containers alongside a shared PostgreSQL database. All config (API keys, DSNs, webhook secrets) lives in a `.env` file on the host.

### How to deploy

1. Make sure `.env` is up to date with the correct credentials (Polar.sh, Elestio, Corteza JWT).
2. Run `make deploy` — builds linux/amd64 images for both services, pushes them to the host, and restarts the stack.
3. For config-only changes (no code change): `make sync` rsyncs the `.env` and compose file and restarts without rebuilding.

---

## human-docs

### What it is

GitBook documentation site for the @Human Connections product — open-source integration/automation platform (Zapier alternative) with 50+ third-party connectors. Docs cover Actions, Queries, and Triggers per connector.

### Stack

Markdown, GitBook. `SUMMARY.md` drives navigation. Python validation scripts check links.

### Corteza integration

None. Separate product by the same company.

---

## corteza-docs

### What it is

Documentation source for Corteza platform, published at docs.cortezaproject.org. Multi-version site (2023.9.x → 2024.9.x) mirroring server release branches.

### Stack

AsciiDoc, Antora (site generator), Docker. `make changelog` runs a Node script to generate release notes.

### Structure

`src/modules/` — integrator-guide, end-user-guide, developer-guide, devops-guide. `drafts/` for unreleased content.

### Status

Changelog-driven release cycle. Feature branches in progress: `2024.9.x-discovery`, `2024.9.x-feature-user-hierarchy`.

---

## Getting started

1. Read the `.cue` file for the resource you're touching — it is the source of truth for model, RBAC, and codegen config.
2. Regenerate (`make server` or the manual cue command) before editing generated files — stale gen is a common source of confusion.
3. For service work: look at the companion `<snake>.go` for the hooks the template calls; look at `<snake>.gen.go` to see what the template already owns.
4. For scope/multitenancy: `pkg/scope` is injected via context; middleware sets it; generated services read it.
5. `server/devdocs/` has detailed specs for actionlog reporting and project revisions. This file covers everything else.

---

## References

- Actionlog reporting + standalone: `server/devdocs/actionlog-reporting.md`
- Project revisions: `server/devdocs/project-revisions.md`
- Chatbot handoff spec: `.claude/chatbot-handoff-spec.md`
- AI session memory (design decisions, blockers, progress): `.claude/projects/.../memory/`
