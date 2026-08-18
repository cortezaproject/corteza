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
tests:
  - client/web/unify/src/sections/taq/utils/trace-path.test.ts
---

# TAQ utils

## Intention

Pure helpers translating between the ng-automation API shape — the graph and
its execution traces — and the VueFlow canvas, plus shared visual constants.

## Map

- `taq-parser.ts` — `automationToVueFlow` (triggers/steps/paths → typed nodes and `addable` edges; resolves icons from the catalog; defines the `FlowNodeData` contract every node/panel reads), `applyDagreLayout` (top-down auto-layout; the only source of node positions), `conditionToShort` / `conditionToSegments` (human-readable branch condition summaries for edges and previews).
- `trace-path.ts` — resolves an execution trace onto the canvas: `traceFrameFor` and `nodeHasTrace` (node → its frame, matched on handle then stepID), and `isEdgeTraversed` (did control flow along this edge).
- `flow-constants.ts` — `EVENT_ICONS`/`getTriggerMeta` (trigger eventType → icon fallbacks), default trigger/action icons, `NODE_DIMENSIONS` (node size and dagre separations), `getNodeCenterOffset` for viewport centering.

## When changing this

- `FlowNodeData.ref` mirrors the backend handle (`trigger_N`/`step_N`) —
  trace matching and upstream references depend on it.
- Node IDs are type-prefixed (`trigger_`, `step_`, `end_`) to stay unique;
  the parser keeps an id→vueFlowId map for path resolution — keep prefixes
  consistent with `useFlowEditor`'s minting.
- An edge counts as traversed only when its source both ran and did not error
  there. End and loop markers the parser invents to pad a dangling arm carry no
  `stepID` and can never have a frame, so reaching one is inferred from its
  source — never off a `branch`, where the trace does not record which arm ran.
- Gateway steps (`kind` starts with `gateway`) map to `branch` nodes,
  `iterator` to iterator nodes with body/done outputs; edge order out of a
  branch encodes If/Else-If/Else ordering.
