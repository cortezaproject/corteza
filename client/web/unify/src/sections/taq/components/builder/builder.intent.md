---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/taq/composables/useSegmentForm.ts
  - lib/vue/src/stores/useAutomationStore.js
  - lib/vue/src/components/field
touched-by:
  - client/web/unify/src/sections/taq/views/Builder.vue
tests: []
---

# TAQ builder panels & forms

## Intention

Everything around the canvas: picking nodes, configuring the selected node,
binding upstream references, running and tracing. Forms are catalog-driven:
the server's function/trigger definitions carry `segments → sections →
elements` whose `input.type` picks a component from the input/viewer
registries — new argument types are added by registering a component, not by
editing forms.

## Map

- `NodePicker.vue` — the step palette dialog: triggers and function steps from the catalog, grouped (`groups[0]`) with fixed group icons, searchable; emits the chosen node type.
- `ConfigSidebar.vue` — right drawer for the selected node: inline label/description edit, `TriggerForm`/`FunctionForm` dispatch, agent-trigger `InputSchemaEditor`, gateway type switch, drag-reorderable branch list with per-path `ConditionBuilder`; relays all changes as events and exposes `applyReference()`.
- `ReferencePanel.vue` — accordion of upstream sources (invoker/runner, trigger properties, step results); expandable typed values (e.g. ComposeRecord loads module fields) emit `{scope, source}` selections.
- `RunModal.vue` — pre-run input: renders only supported trigger property types, prefills from constraint values, emits the run scope.
- `TracePanel.vue` — input/scope/output tabs for one executed step's frame, searchable, error display (parses JSON errors).
- `TraceValueTree.vue` / `TraceValueRow.vue` — recursive typed-value tree (unwraps `{@type,@value}`), search auto-expands matches.
- `condition/ConditionBuilder.vue` — edits a branch-path condition AST: single comparison or and/or group of rows.
- `condition/ConditionRow.vue` — one comparison: variable side (reference), operator, value side (literal via `CFieldEditor` or reference); emits toggleReference per side.
- `form/DynamicForm.vue` — renders `useSegmentForm.processedSegments` as sections of `DynamicInput`s.
- `form/DynamicInput.vue` — one argument: resolves editor from `inputs/registry`, reference chip toggle, disabled/dependency placeholders.
- `form/FunctionForm.vue` — step arguments ↔ `Expr[]` (`value` vs `expr`/`source` reference mode, aggregate targets for map/array types).
- `form/TriggerForm.vue` — trigger constraints ↔ catalog constraint defs (`values[0]['@value']` encoding).
- `form/InputSchemaEditor.vue` — agent/manual trigger declared input params (name/type/required/description; duplicate-name validation).
- `form/CReferenceChip.vue` — the "bound to upstream value" chip: click re-opens panel, clear unbinds.
- `form/inputs/registry.ts` — `input.type` → editor component map (lib CInput\* + local aggregates); fallback InputText.
- `form/inputs/CInputArray.vue` — multi-row list argument; per-row value-or-reference (`toggleRowReference` with stable row targets).
- `form/inputs/CInputFieldValueMap.vue` — record field→value rows (module fields + system fields); per-field value-or-reference.
- `form/inputs/CInputWorkflowInputMap.vue` — Run Workflow step: the selected workflow's declared inputs as value-or-reference rows (registered as `WorkflowInputMap`, alongside `Workflow`/`WorkflowSelector` → lib `CInputWorkflow`).
- `form/viewers/registry.ts` — `input.type` → read-only viewer map used by `StepPreviewPopover`; falls back to `CViewText`.
- `form/viewers/CView*.vue` — one viewer per input type: `Text`, `Select` (label lookup), `Namespace`/`Module`/`User`/`Agent` (ID → display name), `Reference` (chip), `FieldValueMap` (row summary).

## When changing this

- Keep `inputs/registry.ts` and `viewers/registry.ts` type keys in sync — an
  editable type without a viewer degrades step previews to raw text.

> **DRIFT:** the two registries are NOT in sync today. Nine of 24 editable
> types have no viewer — `Record`, `RecordSelector`, `Workflow`,
> `WorkflowSelector`, `WorkflowInputMap`, `Array`, `Boolean`, `Cron`,
> `Interval` — so those step previews render raw text.

- Reference bindings must stay `{scope, source}` end-to-end (panel → chip →
  saved `Expr`); condition references address rows as
  `condition:<edgeId>:<side>:<rowIndex>`.
- `activeReferenceArgument` is injected from Builder.vue; inputs must tolerate
  its absence (default null ref) for reuse outside the builder.
