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
- `TriggerNode.vue` — trigger node: enabled/disabled styling, play/stop actions emitting `test`/`cancel` for dry runs, at most one outgoing edge. Stop shows whenever this node's dry run still holds a session, paused included — the spinner tracks the poll, the stop control tracks the session.
- `TerminationNode.vue` — end node: target handles only, single inbound.
- `VisualNode.vue` — swimlane and rich-text content visuals; NodeResizer when selected; swimlane paints a drop-target highlight while a step is dragged over it.
- `NodePreview.vue` — body-teleported hover popup (title/description/rows) positioned by `useNodePreview` to escape VueFlow's transformed stacking context.

## When changing this

- Handle IDs (`source-bottom`, `target-top-left`, …) are the codec's anchor vocabulary (`HANDLE_TO_MX` in `lib/codec.js`) — renaming or moving one breaks edge-anchor round-trips with server-stored workflows.
- Handle availability is driven from outside: WorkflowEditor passes used-source/target-handle sets and in/out counts from its `nodeConnections` map; a used handle is rendered non-connectable but kept in the DOM so its edge still anchors.
- Run `WorkflowNode.test.js`. The handle sets differ on purpose: `WorkflowNode` renders 5 of the 8 target and 5 of the 8 source positions the codec can persist, `TerminationNode` all 8 targets, `TriggerNode` all 8 sources. Do not align them.
- `ConnectionMode.Loose` lets an edge end on another node's source handle; `handlesToMxStyle` then persists `entryX=0.5;entryY=0`, which decodes to `target-top`, so the edge relocates on reload. Accepted as is — a cosmetic relocation, not corruption.
