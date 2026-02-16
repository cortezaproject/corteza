# Taq (Automation Builder)

Visual automation/workflow builder using VueFlow. Create triggers, steps, branches, and end nodes in a flow editor.

## Routes

| Path | Name | View | Description |
|------|------|------|-------------|
| `/` | `dashboard` | `Dashboard.vue` | List automations |
| `/builder` | `builder` | `Builder.vue` | Create new automation |
| `/builder/:id` | `builder-edit` | `Builder.vue` | Edit existing automation |
| `/:pathMatch(.*)*` | — | redirect to `/` | Catch-all |

## Stores

### useAutomationStore (`stores/automation.ts`)

- **State:** `list[]` (NgAutomationInstance), `loading`, `error`, `functions[]`, `triggers[]`
- **Actions:**
  - `fetchList(api, filter)` — list automations
  - `create(api, data)` — create and add to list
  - `remove(api, automationID)` — delete
  - `loadFunctions(api)` / `loadTriggers(api)` — load catalogs
  - `loadCatalog(api)` — load both in parallel

## Composables

### useSegmentForm (`composables/useSegmentForm.ts`)

Shared composable for segment-driven forms. Accepts a value strategy (read/write callbacks) and returns `processedSegments` + `updateValue()`. Handles dependency resolution, cascade clearing, and disabled state computation. Used by FunctionForm and TriggerForm.

**Options:** `{ segments, parameters, getValue, getAggregateValue, onUpdate, getReferenceInfo?, upstreamResults? }`
**Returns:** `{ processedSegments, updateValue }`

When `getReferenceInfo` and `upstreamResults` are provided, each processed input includes `isReference`, `referenceLabel` fields for displaying step result references.

### useFlowEditor (`composables/useFlowEditor.ts`)

Core composable managing the VueFlow automation editor. ~800 lines.

**State:** `automation`, `nodes[]`, `edges[]`, `loading`, `saving`, `running`

**Key operations:**
- `load(id)` — fetch automation, convert to VueFlow via `automationToVueFlow()`
- `save()` — convert back via `vueFlowToAutomation()`, push to API
- `exec()` — run automation with tracing
- `addNode(type, insertionPoint)` — insert node with auto-layout
- `deleteNode(node)` — remove with edge reconnection
- `undo()` / `redo()` — JSON-serialized history
- `addBranchOutput(branchNode)` — add "Else If" output
- `reorderBranchEdges(id, order)` — reorder branch outputs
- `relayout()` — re-apply dagre layout
- `getUpstreamResults(nodeId)` — walk edges backward to find ancestor step nodes, return their function results for step-result referencing

## Views

### Dashboard.vue

Card grid of automations. Create dialog with name/description. Delete with confirmation. Shows trigger/step counts and last update date.

### Builder.vue

Full VueFlow canvas with:
- Custom node types: `trigger`, `step`, `branch`, `end`
- Custom edge type: `addable` (+ button to insert nodes)
- Bottom toolbar: zoom controls, undo/redo, run, save, back
- Right-side config sidebar (resizable drawer)
- Reference panel (slides in left of sidebar) for selecting step result references
- Node picker dialog for adding triggers/steps
- Edge highlighting (ancestor walk-up on selection)
- Keyboard deletion (Delete/Backspace)

## Components

### Common
- `TaqIcon.vue` — renders an `IconDef` object. Supports `type: 'name'` (PrimeIcons), `type: 'url'` (image URL), `type: 'attachment'` (attachment ID, placeholder). Props: `icon?: IconDef`, `fallback?: IconDef`.

### Flow Nodes
- `TriggerNode.vue` — entry point node with delete action
- `StepNode.vue` — action node with delete action
- `BranchNode.vue` — conditional gateway with multiple outputs
- `EndNode.vue` — termination node (non-selectable)
- `AddableEdge.vue` — edge with + button to insert nodes between

### Builder UI
- `NodePicker.vue` — dialog to select functions/triggers from catalog
- `ConfigSidebar.vue` — right drawer to configure selected node (shows FunctionForm for steps, TriggerForm for triggers). Exposes `applyReference()` for applying step result references from the ReferencePanel.
- `ReferencePanel.vue` — panel showing upstream step results for binding to input fields. Displays an accordion of ancestor steps with their available result outputs.

### Form System
Architecture: `useSegmentForm` composable (shared logic) → type-specific controller → `DynamicForm` renderer → `DynamicInput` components.

- `DynamicForm.vue` — pure renderer: takes `processedSegments`, renders `DynamicInput` list
- `FunctionForm.vue` — controller for step/function nodes. Reads/writes `Expr[]` arguments. Supports step result references via `scope`/`source` fields on Expr. Exposes `onReferenceSelect()` for applying references. Uses `useSegmentForm`.
- `TriggerForm.vue` — controller for trigger nodes. Reads/writes `TriggerConstraint[]` (constraint values with `{ name, op: "=", values: [{ @value, @type }] }`). Uses `useSegmentForm`.
- `DynamicInput.vue` — individual input field component. Supports dual-mode: static value or step result reference (shows reference chip with clear button)
- `CInputFieldValueMap.vue` — field mapping input type

## Utilities

### flow-constants.ts

- `NODE_DIMENSIONS` — width: 260px, height: 100px, h-sep: 200px, v-sep: 100px
- `DEFAULT_TRIGGER_ICON` / `DEFAULT_ACTION_ICON` — `IconDef` fallbacks (re-exported from `DEFAULT_ICONS`)
- `TRIGGER_META` — `IconDef`/label map for 6 event types (onManual, onInterval, onTimestamp, afterCreate, afterUpdate, afterDelete)

### taq-parser.ts

Converts between API format (NgAutomation) and VueFlow format (nodes/edges).

- `automationToVueFlow(automation, catalog?)` — API → VueFlow (creates nodes, edges, end nodes, applies layout). Pass `{ functions, triggers }` catalog to resolve icons inline from the construct library.
- `vueFlowToAutomation(state)` — VueFlow → API (maps nodes to triggers/steps/paths)
- `applyDagreLayout(state)` — auto-position nodes top-to-bottom using dagre

## Dependencies

Extra dependencies beyond the shared stack:
- `@vue-flow/core` + `@vue-flow/background` + `@vue-flow/controls` — flow editor
- `dagre` — graph layout algorithm
