---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/taq/utils/flow-constants.ts
  - client/web/unify/src/sections/taq/utils/taq-parser.ts
touched-by:
  - client/web/unify/src/sections/taq/views/Builder.vue
tests: []
---

# TAQ flow canvas renderers

## Intention

Custom VueFlow node and edge components for the builder canvas. They render
graph structure and state only — all mutation goes back to the Builder via
events (`delete`, `replace`, edge `add`); nodes read the freshest data from
the passed `nodes` array, never from VueFlow's (possibly stale) props.

## Map

- `TriggerNode.vue` — flow root; shows trigger label/icon from the triggers catalog; context menu (replace/delete).
- `StepNode.vue` — function step card; hover (or "show all") opens `StepPreviewPopover`; trace tint when executed.
- `BranchNode.vue` — gateway (exclusive/inclusive) with N ordered outputs; per-path condition labels on edges.
- `IteratorNode.vue` — loop-over-collection step with two outputs: body (ends in a loop node) and done.
- `EndNode.vue` — non-selectable leaf; persisted as a termination step.
- `LoopNode.vue` — non-selectable marker closing an iterator body.
- `AddableEdge.vue` — every edge type; hoverable "+" emits `add {edgeId, source, target}` for mid-edge insertion; renders branch/iterator path labels, ancestor-highlight and trace-traversed styling.
- `StepPreviewPopover.vue` — read-only config summary of a node: resolves the catalog definition and renders argument values via `form/viewers/registry` (references shown as chips, conditions via `conditionToSegments`).

## When changing this

- Trace props (`traceFrame`, `traceActive`) are display-only; executed-state
  matching by handle/stepID lives in Builder.vue.
- Node dimensions must stay in sync with `NODE_DIMENSIONS` — dagre layout and
  viewport centering assume them.
