---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/cast.ts
  - lib/js/src/guards.ts
touched-by:
  - client/web/unify
  - lib/vue
tests: []
---

# Automation types

## Intention

Client-side types for the automation subsystem: classic `Workflow`,
`Function`/`Param` (workflow function definitions), `Prompt`, execution
traces, typed value encoding, and the TAQ model. TAQ = **Trigger Action
Query** (never "Task Queue").

## Map

- `types/taq.ts` — `TAQ` + `NgAutomation` classes and the ng-automation graph model: `NgAutomationStep`/`NgAutomationPath`/`NgAutomationTrigger`, step argument/result `Expr`s, `TriggerConstraint`, and `TAQIssue`/`TAQIssueSet` validation results.
- `types/workflow.ts` — classic workflow resource class.
- `types/values.ts` — `Typed`/`Vars` wire format (`{'@type', '@value'}`) and `Encode`/`IsTyped` for exchanging typed values with the automation engine.
- `types/function.ts`, `types/param.ts` — workflow function catalog entries and their parameter definitions.
- `types/trace.ts` — execution results/stack frames/statuses for run inspection UIs.
- `types/prompt.ts`, `types/icon.ts` — workflow prompts; step icon defs (`IconDef`, `normalizeIcon`, `DEFAULT_ICONS`).

## Contract (what apps may rely on)

- `Expr.scope` references another step's handle and `Expr.source` a named result of that step — the mechanism behind step-output wiring in the TAQ editor (e.g. Run Workflow named I/O).
- `Vars` values must be `Typed` (`Encode` wraps plain values); sending unwrapped values to the engine is a bug.
- Shapes in `taq.ts` mirror server-side ng-automation types — they change together with the server, not independently.

## When changing this

- The TAQ editor and run/trace views in unify deserialize directly into these types; keep them aligned with server codegen'd REST payloads.
