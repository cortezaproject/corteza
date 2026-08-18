---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/taq/utils/taq-parser.ts
  - lib/vue/src/stores/useAutomationStore.js
  - lib/js/src/automation/types/trace.ts
touched-by:
  - client/web/unify/src/sections/taq/views/Builder.vue
  - client/web/unify/src/sections/taq/components/builder/form/FunctionForm.vue
  - client/web/unify/src/sections/taq/components/builder/form/TriggerForm.vue
tests: []
---

# TAQ composables

## Intention

The builder's logic layer, kept out of the components.

## Map

- `useFlowEditor.ts` — the entire editor model: load/save/exec a TAQ, node/edge mutation (add, replace, delete-with-reconnect, branch outputs, gateway type, edge conditions), snapshot undo/redo, dirty tracking, trace state, `getUpstreamResults` (reference panel sources incl. always-present invoker/runner system-user fields), `getTriggerProperties` (RunModal inputs).
- `useSegmentForm.ts` — shared segment-form processing for FunctionForm/TriggerForm: flattens catalog segments into renderable inputs, resolves `context.dependsOn` (disabled placeholders, cascade-clears dependents on change, design-time reference resolution), detects aggregate types (FieldValueMap/WorkflowInputMap/Array).

## When changing this

- Save serializes nodes/edges → `triggers`/`steps`/`paths`: existing backend
  IDs are preserved and new IDs minted above the current max so stored scope
  references (`step_<id>`) never break; end nodes become `termination` steps;
  gateway edge order is meaningful (last = Else) and conditions ride on paths.
  Arguments and path conditions go through `utils/reference-binding` on the way
  out, because the `invoker`/`runner` scopes the panel offers are not scopes the
  runtime can resolve.
- After save the flow is re-parsed from the server response — local-only node
  state does not survive; a new TAQ rewrites the URL to its id.
- Delete keeps every parent a non-leaf: terminations are re-added and orphans
  (no incoming edge, non-trigger) are recursively removed.
- History is JSON snapshots capped at 100; `isDirty` compares history index to
  last-saved index — bump `lastSavedHistoryIndex` on any new save path.
- Exec runs with `trace: true` then fetches frames; frames match nodes by
  `handle` first, `stepID` fallback.
