# Workflow Editor: mxGraph → Vue Flow Migration

## Context
The workflow editor (`WorkflowEditor.vue`, ~2900 lines) uses mxGraph 4.2.0, which is incompatible with Vue 3's reactive Proxy system. mxGraph relies on object identity (`===`) for internal state, but Vue 3 wraps objects in deep Proxies, causing edges to permanently disappear on hover and after save. The fix is to replace mxGraph entirely with `@vue-flow/core` — a Vue 3-native library already used in the taq app. All features of the current editor must be preserved.

---

## Package Changes (`workflow/package.json`)

```diff
- "mxgraph": "^4.2.0"
+ "@vue-flow/core": "^1.48.2"
+ "@vue-flow/background": "^1.3.2"
+ "@vue-flow/node-resizer": "^1.4.0"
```
Import `@vue-flow/core/dist/style.css` in `plugins/index.js` or `WorkflowEditor.vue`.

---

## Files to KEEP Unchanged
- `src/components/Configurator/*` — all configurator components (no mxGraph dependency)
- `src/views/Editor.vue` — save/load/API logic (unchanged)
- `src/lib/dry-run.js`, `lib/toolbar.js`, `lib/style.js`, `lib/constraint.js`, `lib/string.js`, `lib/eventBus.js`

## Files to REWRITE
- `src/components/WorkflowEditor.vue` — core canvas (mxGraph → VueFlow)
- `src/lib/codec.js` — serialization (mxGraph model → VueFlow nodes/edges arrays)

## New Files to CREATE
```
src/components/FlowNodes/WorkflowNode.vue       — generic step node
src/components/FlowNodes/TriggerNode.vue        — trigger with test/spinner/cancel controls
src/components/FlowNodes/TerminationNode.vue    — end node
src/components/FlowNodes/VisualNode.vue         — swimlane/content annotation (resizable)
src/components/FlowEdge.vue                     — orthogonal edge with editable label
src/composables/useWorkflowHistory.js           — undo/redo history stack
src/composables/useWorkflowDnD.js               — drag from toolbar to canvas
src/composables/useWorkflowHighlight.js         — ctrl+click path highlighting + trace
src/composables/useWorkflowClipboard.js         — copy/cut/paste with ID remapping
```

---

## Data Model

### VueFlow Node Shape
```js
{
  id: "42",                         // = stepID / trigger visual id
  type: "workflow" | "trigger" | "termination" | "visual",
  position: { x: 200, y: 300 },    // free positioning, 8px grid snap
  data: {
    stepID: "42",
    kind: "function",               // step kind from lib/style.js
    ref: "compose:record:create",
    label: "Create record",         // = cell.value in mxGraph
    description: "",
    arguments: [...],
    results: [...],
    defaultName: false,
    triggers: { ... } | undefined,  // trigger nodes only
    highlighted: false,             // ctrl+click highlight
    traceState: null,               // null | "success" | "error"
    traceLog: null,                 // execution time string
  }
}
```

### VueFlow Edge Shape
```js
{
  id: "p-42-55",
  source: "42",
  target: "55",
  type: "workflow",
  label: "My label",
  data: {
    expr: "",           // gateway condition expression
    parentID: "42",
    childID: "55",
    highlighted: false,
    traceState: null
  }
}
```

---

## `codec.js` — New API

### `decodeWorkflow(workflow, triggers)` → `{ nodes, edges }`
1. Trigger array → trigger nodes: `node.id = trigger.meta.visual.id`, position from `xywh[0,1]`, `data.triggers = { resourceType, eventType, constraints, enabled, stepID }`
2. Steps array → workflow/termination/visual nodes: position from `step.meta.visual.xywh[0,1]`
3. Paths array → workflow edges: `edge.id = path.meta.visual.id`
4. **Trigger→step edges**: reconstruct from `trigger.stepID` as an edge `source=trigger.id, target=trigger.stepID` (NOT from paths array)

### `encodeWorkflow(nodes, edges)` → `{ steps, paths, triggers }`
1. Trigger nodes → triggers: `trigger.stepID = first outgoing edge's target id`; outgoing edges excluded from `paths[]`
2. Non-trigger nodes → steps: position stored in `meta.visual.xywh`
3. Non-trigger edges → paths array

---

## `WorkflowEditor.vue` — New Structure

**Script:** Composition API (`<script setup>`):
```js
const nodes = ref([])
const edges = ref([])
const { saveToHistory, undo, redo, canUndo, canRedo } = useWorkflowHistory(nodes, edges)
const { onToolbarDragStart, onCanvasDrop } = useWorkflowDnD(nodes, saveToHistory)
const { highlightConnected, clearHighlights } = useWorkflowHighlight(nodes, edges)
const { copySelected, cutSelected, pasteClipboard } = useWorkflowClipboard(nodes, edges, saveToHistory)
const { fitView, zoomIn, zoomOut, getZoom } = useVueFlow()
```

**Template:**
```
[toolbar strip] [VueFlow canvas] [sidebar]
```

VueFlow config:
```vue
<VueFlow
  v-model:nodes="nodes" v-model:edges="edges"
  :snap-to-grid="true" :snap-grid="[8,8]"
  :min-zoom="0.1" :max-zoom="3"
  :pan-on-drag="[1,2]" :pan-on-scroll="true"
  :zoom-on-scroll="false" :zoom-on-double-click="false"
  :delete-key-code="null"
  :connection-mode="ConnectionMode.Loose"
  @connect="onConnect"
  @node-click="onNodeClick"
  @edge-click="onEdgeClick"
  @pane-click="onPaneClick"
  @node-drag-stop="onNodeDragStop"
>
  <Background variant="dots" />
  <template #node-workflow="props"><WorkflowNode ... /></template>
  <template #node-trigger="props"><TriggerNode ... /></template>
  <template #node-termination="props"><TerminationNode ... /></template>
  <template #node-visual="props"><VisualNode ... /></template>
  <template #edge-workflow="props"><FlowEdge ... /></template>
</VueFlow>
```

---

## Node Components

### `WorkflowNode.vue`
```
Props: id, data, selected, issues (map), functionTypes, eventTypes
Template:
  ┌─────────────────────────────────┐
  │ [icon] Step Type Label  [id/🔴] │  ← header row
  ├─────────────────────────────────┤
  │ Step label (hover to expand)    │  ← truncate by default, whitespace-normal on hover
  ├─────────────────────────────────┤
  │ Values table (v-if hasValues)   │  ← args/results
  └─────────────────────────────────┘
  Handle(target, top) + Handle(source, bottom)
  CSS: .workflow-node:not(:hover) .vue-flow__handle { opacity: 0 }
```

`valueRows` computed returns `[{key, value}]` array for v-for — replaces mxGraph HTML string generation.

### `TriggerNode.vue`
Same layout as WorkflowNode. Header adds play/spinner/cancel icons. Only source Handle (no target). Opacity 0.7 when disabled.

### `TerminationNode.vue`
Small pill. Only target Handle.

### `VisualNode.vue`
Annotation box. Uses `NodeResizer`. Not connectable. Low z-index.

---

## `FlowEdge.vue`
- Path: `getSmoothStepPath` (orthogonal — replaces mxGraph OrthConnector)
- Label: `EdgeLabelRenderer` with white box + `#A7D0E3` 2px border, double-click to edit
- Stroke: primary when `data.highlighted`, green/red for `data.traceState`

---

## Gateway Edge Auto-Labeling (`onConnect`)
```js
function onConnect(connection) {
  const sourceNode = nodes.value.find(n => n.id === connection.source)
  const outCount = edges.value.filter(e => e.source === connection.source).length
  let label = ''
  if (sourceNode?.data.kind === 'gateway') {
    if (sourceNode.data.ref === 'excl') label = outCount === 0 ? '#1 - If' : `#${outCount+1} - Else (if)`
    else if (sourceNode.data.ref === 'incl') label = 'If'
  } else if (sourceNode?.data.kind === 'error-handler') {
    label = outCount === 0 ? 'Try' : 'Catch'
  } else if (sourceNode?.data.kind === 'iterator') {
    label = outCount === 0 ? 'Body' : 'End'
  }
  edges.value = [...edges.value, { id: `e-${Date.now()}`, ...connection, type:'workflow', label, data:{...} }]
  saveToHistory()
}
```
On excl gateway edge delete: renumber remaining edges.

---

## Sidebar Configurator Compatibility

All existing Configurator components kept **unchanged**. `buildSidebarItem(node)` constructs a compatibility shim:
```js
{
  node: { id, value: data.label, edges: outEdges.map(...), style: data.kind },
  config: { stepID: id, kind, ref, arguments, results, defaultName },
  triggers: data.triggers || undefined
}
```
`sidebarEdges` computed (getter/setter) bridges VueFlow `edges[]` ↔ `{[id]: {node, config}}` map shape the Configurator expects.

---

## Composables Summary

| Composable | Key functions |
|---|---|
| `useWorkflowHistory` | `saveToHistory()`, `undo()`, `redo()` — JSON snapshot, max 100 |
| `useWorkflowDnD` | `onToolbarDragStart()`, `onCanvasDrop()` — screenToFlowPosition + 8px snap |
| `useWorkflowHighlight` | `highlightConnected(id)`, `clearHighlights()` — sets `data.highlighted` |
| `useWorkflowClipboard` | `copySelected()`, `cutSelected()`, `pasteClipboard()` — ID remap + 160px offset |

---

## Keyboard Shortcuts
| Key | Action |
|---|---|
| Ctrl+S | Save |
| Ctrl+Z / Ctrl+Shift+Z | Undo / Redo |
| Ctrl+C / X / V | Copy / Cut / Paste |
| Ctrl+A | Select all |
| Ctrl+Space | Fit view |
| Delete/Backspace | Delete selected (if not in input) |
| Arrow keys | Nudge 40px; Shift = 8px |

---

## Dry-Run Trace
`renderTrace()` sets `node.data.traceState` + `edge.data.traceState` directly. WorkflowNode renders clock badge from `data.traceLog`. FlowEdge applies green/red stroke. Replaces all mxGraph `mxCellHighlight`/`mxCellOverlay` calls.

---

## Implementation Phases

| Phase | Tasks |
|---|---|
| 1. Foundation | Add packages to package.json, write + test new `codec.js` |
| 2. Nodes | Create all 4 node components (display only, no interactions) |
| 3. Canvas | Rewrite WorkflowEditor.vue with VueFlow; wire `decodeWorkflow` in watcher; verify existing workflows render |
| 4. Interactions | FlowEdge.vue, onConnect + auto-label, sidebar open/close, Configurator adapter |
| 5. DnD + Keys | useWorkflowDnD, keyboard shortcuts (delete, nudge, Ctrl+Z/S/C/V/A/Space) |
| 6. Clipboard | useWorkflowClipboard (copy/cut/paste) |
| 7. Advanced | useWorkflowHighlight (ctrl+click), renderTrace (dry-run trace) |
| 8. Polish | Visual styling, remove mxgraph from package.json |

---

## Verification Checklist
1. Load existing saved workflow → all nodes and paths render at correct positions
2. Drag step from toolbar → node appears on canvas, sidebar opens
3. Draw edge between nodes → correct auto-label for gateway/iterator/error-handler
4. Save → reload → graph unchanged (codec round-trip)
5. Undo/redo → state restores correctly
6. Copy/paste → new nodes with offset positions, no ID conflicts
7. Ctrl+click node → connected edges highlight blue
8. Test workflow (dry-run) → spinner on trigger, trace renders green/red
9. All 14 configurator types open and save correctly
10. Ctrl+Z/S, Delete, arrow nudge, Ctrl+Space all work
