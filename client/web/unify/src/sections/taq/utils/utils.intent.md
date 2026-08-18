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
  - client/web/unify/src/sections/taq/utils/reference-binding.test.ts
---

# TAQ utils

## Intention

Pure helpers translating between the ng-automation API shape — the graph and
its execution traces — and the VueFlow canvas, plus shared visual constants.

## Map

- `taq-parser.ts` — `automationToVueFlow` (triggers/steps/paths → typed nodes and `addable` edges; resolves icons from the catalog; defines the `FlowNodeData` contract every node/panel reads), `applyDagreLayout` (top-down auto-layout; the only source of node positions), `conditionToShort` / `conditionToSegments` (human-readable branch condition summaries for edges and previews).
- `trace-path.ts` — resolves an execution trace onto the canvas: `traceFrameFor` and `nodeHasTrace` (node → its frame, matched on handle then stepID), and `isEdgeTraversed` (did control flow along this edge).
- `reference-binding.ts` — `toWireArgument`/`fromWireArgument` and `toWireCondition`/`fromWireCondition`: the identity scopes the reference panel offers (`invoker`, `runner`) are fields of the global execution scope, not step outputs, so they cross the API boundary spelled the way the runtime resolves them and come back as a scope.
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
- `invoker`/`runner` are not scopes the runtime can resolve. An argument sent
  that way is rejected with `scope.unknown` and the automation never registers;
  a path condition is not scope-checked at all and instead fails the run with
  `eval: unknown scope`. An argument spells them as a dotted **expression** on
  the empty scope, a condition as a dotted **symbol** on `global` — the two
  evaluators differ. The builder
  still addresses them as scopes — that is how the panel groups and highlights
  them — and `reference-binding` translates at the API boundary, so a stored
  binding must come back with its scope intact or every "is this a reference"
  check (all truthiness tests on `scope`) reads it as a literal.
- Gateway steps (`kind` starts with `gateway`) map to `branch` nodes,
  `iterator` to iterator nodes with body/done outputs; edge order out of a
  branch encodes If/Else-If/Else ordering.
