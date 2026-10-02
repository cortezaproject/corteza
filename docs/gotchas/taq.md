# TAQ and automation gotchas

Facts about TAQs (Trigger Action Queries) and workflows — triggers, expressions, the scheduler, execution traces and the workflow canvas — that the code does not make obvious.

## A bare gval function name is a variable

gval is forked in-repo at `server/third_party/gval` (`replace` in `server/go.mod`), so a function name claims the parse only when `(` follows; otherwise `(*Parser).parseVariable` reads it as a variable. This holds for every gval language: `pkg/expr` (value expressions, sanitizers, TAQ, `/system/expressions/evaluate` used by visibility conditions), `pkg/dal/runner_gval.go` and `pkg/rbac/func_expr.go`. A bare `now` reads nil (map scope) or errors `no such key` (`*Vars` scope); it does not return the current time.

**How to apply:** fields named `year`, `split`, `concat` and the like are ordinary scope keys; call a function with parens. `pkg/expr/func_names.go` keeps the name list only because the expression editor mirrors it.

## The manual trigger's namespace and module decide nothing at run time

A page's automation button runs a TAQ by ID (`ngAutomationExec`/`workflowExec`) with an `input` built from the page; `ExecuteAndWait` never merges a trigger's constraints, and nothing dispatches a `compose:record`/`onManual` event. The Automation block supplies namespace, page, record and module; record-list selection buttons add `selected` and `filter`. The constraints only scope the builder's record picker and Reference panel (`getTriggerProperties` in `useFlowEditor.ts`).

**Why:** a required namespace/module there looks load-bearing and is not.

**How to apply:** when a manual TAQ lacks its module, read the button component, not the trigger config. The exec `input` envelope rule is in `dev/agent/README.md` under API gotchas.

## The scheduler's interval match window is one second

`scheduler.onInterval` truncates the tick to the second and requires it to equal the cron expression's next occurrence. The watch loop in `server/pkg/scheduler/service.go` recomputes the delay to the next boundary each cycle, because a `time.Ticker` drifts from its previous delivery and, once past a second, every scheduled automation stops with nothing logged.

**How to apply:** keep the tick aligned; never loosen the matcher (a test asserts "not full minute" does not match). A test of `OnInterval` must spin until `time.Now().Second() == 0`. The scheduler logs at Debug, which may not reach the dev log — probe at Warn.

## A step argument's `source` resolves top-level names only

`{scope: "", source: "invoker.email"}` completes and silently yields `""` — `source` does no dotted traversal and raises no error; dotted `expr` resolves. `FunctionForm.onReferenceSelect` writes references as `{scope, expr}` even though `ReferencePanel` emits `source`, so read `expr` when inspecting a stored binding.

**How to apply:** spell any path as `expr`. The invoker/runner scope rules are in `client/web/unify/src/sections/taq/utils/utils.intent.md`.

## A multi-value argument is repeated Exprs sharing one target

For the aggregate `values` map on `composeRecordsCreate`/`Update`, each value is its own Expr with the same `target`; `server/compose/automation/ng_record_handlers.go` appends same-target Exprs (the `auxVals` loop) and writes them with ascending `Place`. The same code reads namespace and module positionally (`args[0]`, `args[1]`).

**Why:** keeping one Expr per target loses all but the last value on open-and-save, silently.

**How to apply:** treat `target` as non-unique when reading an aggregate argument, emit one Expr per element, and keep namespace and module first.

## A TAQ that writes to its own trigger module loops without a guard

A record-update TAQ writing back to the module it triggers on needs a gateway whose condition is false on the second pass. The guard is per-field: one keyed on `quantity != oldRecord.quantity` does not recompute when a price changes, so a total goes stale while every run reports success.

**How to apply:** derive values with a field value expression where possible; it recomputes on every save without an automation.

## An argument's `expr` sees only its own scope

`expr` evaluates in `contexts[e.Scope]` (`server/automation/types/expr.go`); with no scope it uses the trigger's, so inside an iterator body `record` is still the trigger's record. Reach the loop item with `scope: "<iterator handle>"`. One expression sees one scope, so a formula needing both records cannot be written — split it, or use a compose field value expression.

**How to apply:** `record.values.<recordRefField>` is the target's ID as a string; `.recordID` on it yields nothing and an iterator query built from it runs zero times with the execution still `completed`. A record write that fires a failing automation still returns 200.

## Read a trace frame's `args`, `error` and `input`

A step frame holds what its arguments resolved to under `args` and, on failure, the message under `error`, with the execution `status` set to `failed`. A frame's scope is under `input`, not `scope`; the trigger frame's `input` is the whole run scope as typed envelopes. An iterator frame with no body frames after it means the query matched nothing, not that evaluation failed.

**How to apply:** print the whole frame, never just `stepID`.

## What an execution trace records

Every step `executeStep` runs gets a frame, including a reached termination. A gateway's frame does not say which arm ran — read it off the arm's first step or End frame. Parser-invented loop and End markers have no `stepID` and never a frame. An iterator body emits one frame per pass with the same handle; `traceByHandle` in `useFlowEditor.ts` is a `Map`, so only the last pass survives. `parentID` does not form a chain: it can name an unrecorded frame, and the root's is `0`.

**Why:** the builder's traversed-edge tint must be inferred from these gaps.

**How to apply:** "both endpoints have frames" lights untaken arms; see `isEdgeTraversed` in `client/web/unify/src/sections/taq/utils/trace-path.ts`.

## The TAQ catalog is hand-curated, not derived from events.yaml

The builder palette (`/automation/construct-library/triggers` and `/functions`) is populated by `init()` in `server/automation/service/legacy_triggers.go`. Page, module, namespace and page-layout events in `events.yaml` are not catalogued, and only a small slice of the workflow function registry is reachable. Trigger properties come from shared helpers (`recordProperties()`, `userProperties()`); `manualProperties()` is separate so `page` appears only on `onManual`.

**Why:** an advertised property the runtime never supplies resolves to nothing in an expression.

**How to apply:** for a "missing" scope variable, check `legacy_triggers.go` first; adding to `events.yaml` alone does nothing for TAQs.

## A nested expr-type struct needs a typed `SelectGVal`

In `server/*/automation/expr_types.yaml`, a field whose exprType is another expr type (e.g. `meta: AgentMeta`) generates a selector returning the raw Go struct, so `agent.meta.short` fails with `unknown parameter` while `agent.meta` works. `template.meta` has the same hole.

**How to apply:** set `customGValSelector: true` and hand-write `SelectGVal` returning the typed meta for `"meta"`, delegating to the generated selector otherwise (see `Agent` in `server/system/automation/expr_types.go`). A new expr type also needs `Clone()`, `CastTo…` and an entry in `Registry().AddTypes(...)` in `server/system/service/service.go`. Regenerate with `make codegen-legacy`.

## The exec ledger hears completion twice

`runtime.complete`/`fail` (`server/pkg/automation_exec/runtime/runtime.go`) and `runtime_manager.go` both report the same execution's end, so a terminal-transition hook fires twice unless it checks the previous status (`ledger.apply` does). `id.ID.String()` and `Value()` include JSON quotes; use `strconv.FormatUint(v.Num(), 10)` or `v.Str()`.

**How to apply:** a manual exec answering `manager: executable not found` means the TAQ never registered (any validation issue does it); fix its steps or pick another TAQ rather than debugging the engine.

## vue-flow connects by distance, not hit-testing

Starting a drag is a DOM event on the handle, so z-index, overlap and `pointer-events` decide edge versus node drag. Ending one is geometric: `getClosestHandle` in `@vue-flow/core` scans `handleBounds` within `connectionRadius` (default 20) and finds a target even when it is fully covered.

**How to apply:** to fix "hard to start a connection", raise the handle above the node and pad its hit area (a `::before` with negative `inset`). Do not gate `pointer-events` to free an occluded target; nothing blocks it.
