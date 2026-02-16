# TAQ Application Overview

TAQ (Task Automation Query) is a visual flow builder for creating and managing automations in the Corteza low-code platform. Users design automation workflows with triggers, steps, branching logic, and conditional flows using a node-and-edge interface powered by Vue Flow.

---

## Architecture

### Directory Structure

```
client/web/taq/
├── src/
│   ├── views/
│   │   ├── Dashboard.vue          # Automation list & management
│   │   └── Builder.vue            # Visual flow editor canvas
│   ├── components/
│   │   ├── builder/
│   │   │   ├── NodePicker.vue     # Node selection dialog
│   │   │   ├── ConfigSidebar.vue  # Node configuration panel
│   │   │   └── form/
│   │   │       ├── DynamicForm.vue   # Renders inputs from function segments
│   │   │       ├── DynamicInput.vue  # Resolves input type to component
│   │   │       └── inputs/
│   │   │           ├── registry.ts          # Maps input type strings to Vue components
│   │   │           └── CInputFieldValueMap.vue  # Field-value table for aggregate record values
│   │   └── flow/
│   │       ├── TriggerNode.vue    # Trigger node component
│   │       ├── StepNode.vue       # Step node component
│   │       ├── BranchNode.vue     # Conditional branching node
│   │       ├── EndNode.vue        # Visual terminal node
│   │       └── AddableEdge.vue    # Custom edge with + button
│   ├── composables/
│   │   └── useFlowEditor.ts       # Builder-specific state & logic
│   ├── stores/
│   │   └── automation.ts          # Pinia store for shared state & catalog
│   ├── utils/
│   │   ├── taq-parser.ts          # API -> VueFlow conversion & dagre layout
│   │   └── flow-constants.ts      # Layout dimensions, trigger metadata, icon defaults
│   ├── router/
│   │   └── index.js               # Route definitions
│   └── plugins/                   # Auth, API, i18n setup
```

### Routes

| Path | Name | View | Purpose |
|------|------|------|---------|
| `/` | `dashboard` | `Dashboard.vue` | Automation list & management |
| `/builder` | `builder` | `Builder.vue` | Create new automation |
| `/builder/:id` | `builder-edit` | `Builder.vue` | Edit existing automation |
| `/*` | — | — | Catch-all redirect to `/` |

Also depends on shared input components from `lib/vue/src/components/input/`:
- `CInputNamespace.vue` — Namespace selector with search
- `CInputModule.vue` — Module selector (depends on namespaceID)
- `CInputUser.vue` — User selector with search
- `CInputRecord.vue` — Record selector with search
- `CInputSelect.vue` — Base searchable select

The local input registry (`inputs/registry.ts`) maps backend input type strings to these components (plus PrimeVue's `InputText` as the default fallback).

### Technology Stack

- **Vue 3** with Composition API
- **Vue Flow** for graph visualization
- **Dagre** for automatic graph layout
- **Pinia** for state management
- **PrimeVue** for UI components
- **Tailwind CSS** for styling

---

## Core Logic & Data Flow

### Data Model Conversion

The application maintains **dual representations** of automation data. Conversion from API to VueFlow uses `automationToVueFlow()` in `taq-parser.ts`. The reverse conversion (VueFlow to API) is done inline in `useFlowEditor.save()` — note that `taq-parser.ts` also exports a `vueFlowToAutomation()` function, but it is currently **unused dead code**.

1. **Backend (API) Format** - `NgAutomation`:
   ```typescript
   {
     automationID, handle, meta, enabled,
     triggers: NgAutomationTrigger[],  // Array of triggers
     steps: NgAutomationStep[],        // Array of steps
     paths: NgAutomationPath[]         // Connections between nodes
   }
   ```

2. **Frontend (VueFlow) Format**:
   ```typescript
   {
     nodes: Node<FlowNodeData>[],  // Visual nodes with position
     edges: Edge[]                 // Connections (order matters for branches)
   }
   ```

3. **FlowNodeData** (per-node data):
   ```typescript
   interface FlowNodeData {
     label: string
     description?: string
     icon?: string               // PrimeIcon class (e.g. "pi pi-database")
     nodeType: string            // Function ref or trigger eventType
     config: Record<string, unknown>  // Triggers only (maps to trigger.input)
     arguments: Expr[]           // Steps only (maps to step.arguments)
     ref: string                 // Handle reference
     stepID?: string
     triggerID?: string
   }
   ```

4. **Expr** (argument format, source of truth for step configuration):
   ```typescript
   interface Expr {
     argumentName?: string  // Which parameter (e.g. "namespace", "module")
     target?: string        // For aggregate params (e.g. field name)
     type: string           // From function definition types[0] (e.g. "ID")
     value?: any            // Literal value
     expr?: string          // Expression (future use)
     scope?: string         // Scope (future use)
   }
   ```

### ID Mapping Strategy

- Backend uses simple numeric IDs with a **single shared counter** for both triggers and steps
- Frontend prefixes with type to ensure uniqueness: `trigger_1`, `step_2`, `end_3`
- IDs are assigned in order: triggers first, then terminations, then other steps
- Path handles use format: `path_{sourceId}_{targetId}`
- All IDs and handles use underscores (`_`) as separators for consistency

### End Nodes (Termination Steps)

- End nodes are saved as **termination steps** (`kind: 'termination'`) in the backend
- Display as "End" to users, stored as `termination` internally
- Auto-generated for:
  - Leaf nodes (nodes with no outgoing edges)
  - New branch outputs
- On save: all end nodes become termination steps with paths
- On load: termination steps are displayed as end nodes

### Layout System

Uses Dagre library for automatic hierarchical layout with constants from `flow-constants.ts`:

```typescript
NODE_DIMENSIONS = {
  WIDTH: 260,
  HEIGHT: 100,
  HORIZONTAL_SEP: 200,  // Between sibling branches
  VERTICAL_SEP: 100,    // Between parent-child ranks
}
```

Post-processing:
- Aligns branch children to same Y level
- Positions children left-to-right matching edge array order
- Ensures minimum 120px distance for end nodes connected to branches

---

## Branch Node System

Branch nodes support **multiple outputs** (not just yes/no):

### Edge-Based Branching

- First output: "If" (leftmost)
- Middle outputs: "Else If" (can have multiple)
- Last output: "Else" (rightmost)

Branch output order is determined by **edge array order**, not by any handle identifier.

### Adding/Reordering Branches

```typescript
// Add new "Else If" branch (inserted before last "Else")
addBranchOutput(branchNode)

// Reorder branches via drag-and-drop in ConfigSidebar
reorderBranchEdges(branchNodeId, newEdgeOrder)
```

### Custom Edge Rendering (AddableEdge.vue)

For branch edges, renders a forked path:
1. Drops 30px from source center
2. Horizontal segment to target X position
3. Vertical drop to target

Labels positioned near top of vertical segment, "+" button centered on vertical segment.

---

## Component Behaviors

### Builder.vue (Main Flow Editor)

**Responsibilities:**
- Renders VueFlow canvas with custom node/edge types
- Manages node selection and configuration sidebar
- Handles undo/redo history (JSON snapshots)
- Coordinates node addition/deletion
- Highlights incoming edges when node is selected
- Saves automation to API

**Key Interactions:**
- Middle-click double-tap: Fit to view
- Delete/Backspace: Remove selected elements
- Click on edge "+": Insert node
- Node click: Open config sidebar, highlight ancestor edges

**Viewport Settings:**
- Min zoom: 0.5x, Max zoom: 2x, Default: 1.5x
- Pan on drag (all mouse buttons)
- Zoom on scroll disabled

### NodePicker.vue

**Categories:**
- Triggers - loaded from API catalog (`store.triggers`), labels/descriptions/icons from backend `meta`
- Logic - frontend-defined (Branch/Gateway), uses i18n translations
- Actions - loaded from API catalog (`store.functions`, excluding gateway kind), labels/descriptions/icons from backend `meta`

**Filtering:**
- When adding first node: shows only triggers
- When inserting on edge: shows logic + actions (no triggers)

**Icons:** Stored on backend `meta.icon` without prefix (e.g. `"database"`), frontend prepends `"pi pi-"`. Falls back to `"bolt"` for triggers, `"cog"` for actions.

### ConfigSidebar.vue

**Features:**
- Shows node label and description
- **DynamicForm** for step configuration (rendered from backend function segments)
- For branch nodes:
  - Lists branch outputs in reorderable DataTable
  - Displays "If", "Else If", "Else" tags
  - Drag-to-reorder support
  - "Add Else If Branch" button
- Delete button at bottom

### DynamicForm System

Step configuration is driven by **backend-defined segments**:

1. Each function defines `segments > sections > elements` in Go
2. Each element has `input.type` (e.g. `"NamespaceSelector"`), `input.argument`, and optional `input.context.dependsOn` (context comes from backend JSON; not fully typed in `SegmentInput` interface)
3. **DynamicForm** reads `Expr[]` from node data to populate inputs via `getArgValue()` / `getAggregateValue()`
4. **DynamicInput** resolves `input.type` to a Vue component via the input registry
5. `dependsOn` creates data dependencies (e.g. module selector needs namespace value as `namespaceID` prop)
6. Inputs are auto-disabled when dependencies are unresolved, with contextual placeholder text
7. On change, `updateArgument()` creates proper `Expr` objects with types from function definition and cascade-clears dependent inputs
8. Aggregate parameters (e.g. `FieldValueMap`) produce multiple `Expr` entries sharing the same `argumentName` with different `target` values

### AddableEdge.vue

**Features:**
- Custom forked path rendering for branch edges
- Dynamic labels based on edge order (If/Else If/Else)
- Centered "+" button for node insertion
- Supports edge highlighting (when ancestor is selected)

---

## State Management Architecture

The application uses a **hybrid approach** separating shared state from builder-specific state:

### Pinia Store (`stores/automation.ts`)

Manages **shared state** used across the app:

```typescript
// State
list: NgAutomationInstance[]     // Automation list (Dashboard)
functions: AutomationFunction[]  // Function catalog (with segments, parameters, meta)
triggers: AutomationTrigger[]    // Trigger catalog (with meta for labels/icons)
loading: boolean
error: string | null

// Actions
fetchList(api, filter)   // Load automation list
create(api, data)        // Create new automation
remove(api, automationID) // Delete automation
loadCatalog(api)         // Load functions + triggers (parallel)
loadFunctions(api)       // Load function catalog only
loadTriggers(api)        // Load trigger catalog only
reset()                  // Clear store state
```

### Flow Editor Composable (`composables/useFlowEditor.ts`)

Manages **builder-specific state**:

```typescript
// State
automation: NgAutomation        // Current automation being edited
nodes: Node<FlowNodeData>[]     // VueFlow nodes
edges: Edge[]                   // VueFlow edges
loading: boolean
saving: boolean
running: boolean                // Execution in progress

// History
history: string[]               // JSON snapshots for undo/redo
canUndo: boolean
canRedo: boolean

// Computed
name: string                    // Automation name (get/set)
automationId: string            // Current automation ID
isEmpty: boolean                // True if no nodes

// Actions
load(id)                        // Load automation from API (resolves icons from catalog)
save()                          // Save to API (create or update)
exec()                          // Execute automation via API
reset()                         // Clear to empty state
addNode(nodeType, insertionPoint) // Add node to flow
deleteNode(node)                // Remove node from flow
addBranchOutput(branchNode)     // Add Else If branch
reorderBranchEdges(branchNodeId, newEdgeOrder) // Reorder branches
updateNodeData(nodeId, data)    // Update node data (arguments, config)
undo() / redo()                 // History navigation
saveToHistory()                 // Save current state to undo history

// Internal (not exported)
// cleanupOrphanedNodes()       // Remove disconnected nodes (called by deleteNode)
// relayout()                   // Re-apply dagre layout (called by deleteNode)
```

### Why This Approach?

| Concern | Location | Reason |
|---------|----------|--------|
| Automation list | Store | Shared between Dashboard views |
| Function/trigger catalog | Store | Loaded once, used by NodePicker + icon resolution |
| Current automation | Composable | Only needed in Builder |
| VueFlow nodes/edges | Composable | Transient visual state |
| Undo/redo history | Composable | Builder-specific feature |
| Step arguments (Expr[]) | Node data | Source of truth, read directly on save |

---

## Edge Highlighting

When a node is selected, the Builder highlights all **incoming edges** by walking backwards through the graph:

```javascript
// Find all edges on the path from any trigger to selected node
function getAncestorEdges(nodeId) {
  const ancestorEdges = new Set()
  const visited = new Set()
  const queue = [nodeId]

  while (queue.length > 0) {
    const currentId = queue.shift()
    if (visited.has(currentId)) continue
    visited.add(currentId)

    // Find all incoming edges
    const incoming = edges.filter(e => e.target === currentId)
    incoming.forEach(edge => {
      ancestorEdges.add(edge.id)
      queue.push(edge.source)
    })
  }

  return ancestorEdges
}
```

Highlighted edges receive `data.highlighted = true` and are styled with primary color.

---

## Improvement Recommendations

### Completed

- [x] **Consolidate State Management** - Hybrid approach: Pinia store for shared state, `useFlowEditor` composable for builder state
- [x] **Fix API Injection Consistency** - Store methods accept API as parameter, composable injects API internally
- [x] **Extract Flow Logic to Composable** - Created `useFlowEditor.ts` with all builder logic
- [x] **End Nodes Persisted** - Saved as termination steps (`kind: 'termination'`)
- [x] **Single ID Counter** - Prevents path collisions between triggers and steps
- [x] **Orphan Node Cleanup** - Automatic removal of disconnected nodes on deletion
- [x] **Edge Highlighting** - Visual path highlighting when selecting a node
- [x] **Multi-Output Branches** - If/Else If/Else with dynamic edge ordering
- [x] **Branch Reordering** - Drag-to-reorder in ConfigSidebar
- [x] **Centralized Layout Constants** - Moved to `flow-constants.ts`
- [x] **Dynamic Catalog Usage** - NodePicker uses triggers and functions from API
- [x] **Expr[] as source of truth** - Step arguments stored directly as `Expr[]`, no lossy config conversion
- [x] **Backend-driven labels/icons** - Function and trigger labels, descriptions, icons from backend `meta`; only logic nodes use frontend i18n
- [x] **DynamicForm system** - Backend-defined segments drive step configuration UI
- [x] **Input component registry** - Extensible mapping of input types to Vue components
- [x] **Catalog-before-load** - `catalogReady` ref gates route watcher to prevent race conditions
- [x] **Execution support** - Run button calls `ngAutomationExec` API
- [x] **Placeholder system** - DynamicForm computes disabled/enabled placeholders from `dependsOn` context

### Medium Priority

#### 1. **Add Validation Layer**

Create validation utilities for pre-save checks (at least one trigger, required arguments filled, path integrity).

#### 2. **Add Unit Tests**

Priority test areas:
- `taq-parser.ts` conversion functions (round-trip fidelity)
- `useFlowEditor.ts` composable methods
- DynamicForm argument handling (aggregate, cascade clearing)
- Branch edge ordering logic

### Low Priority

#### 3. **Improve Error Handling**

Add proper error boundaries and user feedback for save conflicts, network errors, etc.

---

## Summary

The TAQ application provides a visual flow builder for creating automations with a backend-driven configuration system.

### Key Design Decisions

- **Expr[] is the source of truth** for step arguments — no intermediate config object, no lossy conversion
- **Backend drives UI** — function labels, descriptions, icons, and form segments come from the construct library API
- **Frontend only handles logic nodes** — Branch/Gateway labels use i18n; everything else comes from backend
- **Catalog loads before automation** — prevents race conditions with icon/metadata resolution
- **Input registry is extensible** — adding a new input type = one component + one registry entry. Current mappings:
  - `UserSelector` / `User` → `CInputUser`
  - `NamespaceSelector` / `Namespace` → `CInputNamespace`
  - `ModuleSelector` / `Module` → `CInputModule`
  - `RecordSelector` / `Record` → `CInputRecord`
  - `Text` / `String` / `Number` → PrimeVue `InputText`
  - `Select` / `Dropdown` → `CInputSelect`
  - `FieldValueMap` → `CInputFieldValueMap` (local to taq)

### Remaining Work

1. **Validation** - Add proper pre-save validation
2. **Unit tests** - Add tests for parser, composable, and DynamicForm
3. **More input types** - Expression editor, code editor, etc. (User/Record selectors already added)
4. **Trigger configuration** - Triggers currently use `config` (flat object), segments not yet defined for triggers
5. **Dead code cleanup** - Remove unused `vueFlowToAutomation()` from `taq-parser.ts` (save logic is inline in `useFlowEditor`)
