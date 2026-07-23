---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/automation/types/icon.ts
touched-by:
  - client/web/unify/src/sections/taq/composables/useFlowEditor.ts
  - client/web/unify/src/sections/taq/components/flow/StepPreviewPopover.vue
tests: []
---

# TAQ utils

## Intention

Pure helpers translating between the ng-automation API shape and the VueFlow
canvas, plus shared visual constants.

## Map

- `taq-parser.ts` — `automationToVueFlow` (triggers/steps/paths → typed nodes and `addable` edges; resolves icons from the catalog; defines the `FlowNodeData` contract every node/panel reads), `applyDagreLayout` (top-down auto-layout; the only source of node positions), `conditionToShort` / `conditionToSegments` (human-readable branch condition summaries for edges and previews).
- `flow-constants.ts` — `EVENT_ICONS`/`getTriggerMeta` (trigger eventType → icon fallbacks), default trigger/action icons, `NODE_DIMENSIONS` (node size and dagre separations), `getNodeCenterOffset` for viewport centering.

## When changing this

- `FlowNodeData.ref` mirrors the backend handle (`trigger_N`/`step_N`) —
  trace matching and upstream references depend on it.
- Node IDs are type-prefixed (`trigger_`, `step_`, `end_`) to stay unique;
  the parser keeps an id→vueFlowId map for path resolution — keep prefixes
  consistent with `useFlowEditor`'s minting.
- Gateway steps (`kind` starts with `gateway`) map to `branch` nodes,
  `iterator` to iterator nodes with body/done outputs; edge order out of a
  branch encodes If/Else-If/Else ordering.
