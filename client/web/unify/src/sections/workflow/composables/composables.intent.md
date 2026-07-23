---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/workflow/lib/id.js
  - client/web/unify/src/sections/workflow/lib/style.js
touched-by:
  - client/web/unify/src/sections/workflow/components/WorkflowEditor.vue
  - client/web/unify/src/sections/workflow/components/FlowNodes
tests: []
---

# Workflow canvas composables

## Intention

Stateful canvas behaviors factored out of WorkflowEditor. Each operates
directly on the shared reactive `nodes`/`edges` refs and snapshots via the
caller-supplied `saveToHistory`.

## Map

- `useWorkflowHistory.js` — undo/redo stack, delta-compressed with periodic full checkpoints; restore reuses existing VueFlow node objects by id so measured internals (dimensions, handleBounds) survive, then remeasures changed nodes via `updateNodeInternals`.
- `useWorkflowDnD.js` — toolbar→canvas drag: node-sized dashed drag ghost, drop creates a node at VueFlow-projected coords with `nextId`.
- `useWorkflowClipboard.js` — copy/cut/paste with mxGraph parity: selecting a swimlane implicitly includes nested children, edges crossing the selection boundary are dropped, paste remaps IDs and offsets +160px; payload also written to the system clipboard as prefixed JSON so paste works across tabs.
- `useWorkflowHighlight.js` — Ctrl+click toggle of `data.highlighted` on a node and its connected edges/nodes.
- `useNodePreview.js` — fixed-position coords for hover popups teleported out of VueFlow's transformed wrapper.

## When changing this

- History restore vs. VueFlow internals is the fragile spot: replacing node objects wholesale loses measurements and detaches edges — preserve the reuse-by-id behavior.
- The clipboard's magic prefix + JSON shape is a compatibility contract with itself across tabs/sessions; change it only with fallback parsing.
