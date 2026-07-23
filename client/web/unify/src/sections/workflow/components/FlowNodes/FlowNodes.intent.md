---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/workflow/lib/style.js
  - client/web/unify/src/sections/workflow/lib/icon.js
  - client/web/unify/src/sections/workflow/composables/useNodePreview.js
touched-by:
  - client/web/unify/src/sections/workflow/components/WorkflowEditor.vue
tests:
  - client/web/unify/src/sections/workflow/components/FlowNodes/WorkflowNode.test.js
---

# Flow nodes

## Intention

Custom VueFlow node renderers for the workflow canvas — the visual vocabulary
of the editor. Each exposes a fixed set of named connection handles that the
codec maps to legacy mxGraph anchor coordinates.

## Map

- `WorkflowNode.vue` — generic step node: kind icon + label, issues badge (opens issues modal via `open-issues`), trace success/error state, hover preview.
- `TriggerNode.vue` — trigger node: enabled/disabled styling, play/stop actions emitting `test`/`cancel` for dry runs, at most one outgoing edge.
- `TerminationNode.vue` — end node: target handles only, single inbound.
- `VisualNode.vue` — swimlane and rich-text content visuals; NodeResizer when selected; swimlane paints a drop-target highlight while a step is dragged over it.
- `NodePreview.vue` — body-teleported hover popup (title/description/rows) positioned by `useNodePreview` to escape VueFlow's transformed stacking context.

## When changing this

- Handle IDs (`source-bottom`, `target-top-left`, …) are the codec's anchor vocabulary (`HANDLE_TO_MX` in `lib/codec.js`) — renaming or moving one breaks edge-anchor round-trips with server-stored workflows.
- Handle availability is driven from outside: WorkflowEditor passes used-source/target-handle sets and in/out counts from its `nodeConnections` map; a used handle is rendered non-connectable but kept in the DOM so its edge still anchors.
- Run `WorkflowNode.test.js`; keep the handle sets of workflow/termination nodes consistent with each other.
