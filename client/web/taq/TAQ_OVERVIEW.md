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
│   │   └── Builder.vue            # Visual flow editor canvas (~438 lines)
│   ├── components/
│   │   ├── builder/
│   │   │   ├── NodePicker.vue     # Node selection dialog (~197 lines)
│   │   │   └── ConfigSidebar.vue  # Node configuration panel (~145 lines)
│   │   └── flow/
│   │       ├── TriggerNode.vue    # Trigger node component
│   │       ├── StepNode.vue       # Step node component
│   │       ├── BranchNode.vue     # Conditional branching node
│   │       ├── EndNode.vue        # Visual terminal node
│   │       └── AddableEdge.vue    # Custom edge with + button (~169 lines)
│   ├── composables/
│   │   └── useFlowEditor.ts       # Builder-specific state & logic (~716 lines)
│   ├── stores/
│   │   └── automation.ts          # Pinia store for shared state (~181 lines)
│   ├── utils/
│   │   ├── taq-parser.ts          # API <-> VueFlow conversion (~474 lines)
│   │   └── flow-constants.ts      # Layout & dimension constants
│   ├── router/
│   │   └── index.js               # Route definitions
│   └── plugins/                   # Auth, API, i18n setup
```

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

The application maintains **dual representations** of automation data:

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
- Triggers - loaded from API catalog (`store.triggers`)
- Logic - frontend-defined (Branch/Gateway)
- Actions - loaded from API catalog (`store.functions`, excluding gateway kind)

**Filtering:**
- When adding first node: shows only triggers
- When inserting on edge: shows logic + actions (no triggers)

### ConfigSidebar.vue

**Features:**
- Shows node label and description
- For branch nodes:
  - Lists branch outputs in reorderable DataTable
  - Displays "If", "Else If", "Else" tags
  - Drag-to-reorder support
  - "Add Else If Branch" button
- Delete button at bottom

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
functions: AutomationFunction[]  // Function catalog
eventTypes: AutomationEventType[] // Event type catalog
loading: boolean
error: string | null

// Actions
fetchList(api, filter)   // Load automation list
create(api, data)        // Create new automation
remove(api, automationID) // Delete automation
loadCatalog(api)         // Load functions + event types
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

// History
history: string[]               // JSON snapshots for undo/redo
canUndo: boolean
canRedo: boolean

// Computed
name: string                    // Automation name (get/set)
automationId: string            // Current automation ID
isEmpty: boolean                // True if no nodes

// Actions
load(id)                        // Load automation from API
save()                          // Save to API (create or update)
reset()                         // Clear to empty state
addNode(nodeType, insertionPoint) // Add node to flow
deleteNode(node)                // Remove node from flow
addBranchOutput(branchNode)     // Add Else If branch
reorderBranchEdges(branchNodeId, newEdgeOrder) // Reorder branches
undo() / redo()                 // History navigation
cleanupOrphanedNodes()          // Remove disconnected nodes
```

### Why This Approach?

| Concern | Location | Reason |
|---------|----------|--------|
| Automation list | Store | Shared between Dashboard views |
| Function catalog | Store | Loaded once, used by NodePicker |
| Current automation | Composable | Only needed in Builder |
| VueFlow nodes/edges | Composable | Transient visual state |
| Undo/redo history | Composable | Builder-specific feature |

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

## Remaining Inconsistencies

### 1. **Missing Validation**

**Issue:** The `save()` function has minimal validation:

```javascript
async function save() {
  // No validation that automation has at least one trigger
  // No validation of step configurations
  // No validation of path integrity
}
```

**Impact:** Users can save invalid automations that may fail at runtime.

### 3. **Missing i18n for Node Labels**

**Issue:** Node labels in `NodePicker.vue` and `ConfigSidebar.vue` are hardcoded English strings:

```javascript
{ label: 'Manual', description: 'Trigger manually' }
{ label: 'Record Create', description: 'After record is created' }
"Add Else If Branch"
"Branches"
```

**Impact:** Breaks i18n compliance, not localizable.

---

## Improvement Recommendations

### Completed

- [x] **Consolidate State Management** - Implemented hybrid approach with Pinia store for shared state and `useFlowEditor` composable for builder-specific state
- [x] **Fix API Injection Consistency** - Store methods consistently accept API as parameter, composable injects API internally
- [x] **Extract Flow Logic to Composable** - Created `useFlowEditor.ts` with all builder logic
- [x] **End Nodes Persisted** - End nodes saved as termination steps (`kind: 'termination'`)
- [x] **Single ID Counter** - Unified ID counter prevents collisions between triggers and steps
- [x] **Orphan Node Cleanup** - Automatic cleanup of disconnected nodes on deletion
- [x] **Edge Highlighting** - Visual feedback showing path to selected node
- [x] **Multi-Output Branches** - Support for If/Else If/Else with dynamic edge ordering
- [x] **Branch Reordering** - Drag-and-drop reordering in ConfigSidebar
- [x] **Centralized Layout Constants** - Moved to `flow-constants.ts`
- [x] **Dynamic Catalog Usage** - NodePicker now uses triggers and functions from API instead of hardcoded values

### Medium Priority

#### 2. **Add Validation Layer**

Create validation utilities:

```typescript
// utils/validation.ts
export function validateAutomation(automation: NgAutomation): ValidationResult {
  const errors: string[] = []

  if (!automation.triggers?.length) {
    errors.push('Automation must have at least one trigger')
  }

  // Validate paths form connected graph
  // Validate required step arguments
  // etc.

  return { valid: errors.length === 0, errors }
}
```

#### 3. **Internationalize Node Labels**

Move labels to locale files:

```yaml
# locale/en/corteza-webapp-taq/builder.yaml
nodePicker:
  triggers:
    manual:
      label: Manual
      description: Trigger manually
configSidebar:
  branches: Branches
  addElseIf: Add Else If Branch
```

### Low Priority

#### 4. **Add Unit Tests**

Priority test areas:
- `taq-parser.ts` conversion functions
- `useFlowEditor.ts` composable methods
- Branch edge ordering logic
- Orphan cleanup logic
- Undo/redo functionality

```typescript
// __tests__/taq-parser.test.ts
describe('automationToVueFlow', () => {
  it('converts triggers to trigger nodes with prefixed IDs', () => {...})
  it('generates end nodes for leaf nodes', () => {...})
  it('handles branch nodes with multiple outputs', () => {...})
})
```

#### 5. **Improve Error Handling**

Add proper error boundaries and user feedback:

```javascript
async function save() {
  const validation = validateAutomation(automation.value)
  if (!validation.valid) {
    $toast.toastWarning(validation.errors.join('\n'), t('builder.validation.title'))
    return
  }

  try {
    // ... save logic
  } catch (e) {
    if (e.response?.status === 409) {
      $toast.toastWarning(t('builder.errors.conflict'))
    } else {
      $toast.toastDanger(t('builder.errors.saveFailed'))
    }
  }
}
```

---

## Summary

The TAQ application has a solid foundation with Vue Flow integration and clean visual design.

### Recent Improvements

- **State management refactored** - Clear separation between shared state (Pinia store) and builder-specific state (`useFlowEditor` composable)
- **Builder.vue simplified** - Reduced from ~760 lines to ~438 lines by extracting logic to composable
- **End nodes persisted** - End nodes are now saved as termination steps (`kind: 'termination'`) in the backend
- **Multi-output branches** - Branch nodes now support If/Else If/Else with multiple outputs
- **Branch management** - Drag-to-reorder branches in ConfigSidebar, "Add Else If" button
- **Edge highlighting** - Visual path highlighting when selecting a node
- **Single ID counter** - Prevents path collisions between triggers and steps
- **Orphan cleanup** - Automatic removal of disconnected nodes
- **Centralized constants** - Layout dimensions moved to `flow-constants.ts`

### Remaining Work

1. **Validation** - Add proper pre-save validation
2. **i18n compliance** - Move all strings to locale files
3. **Unit tests** - Add tests for parser and composable

### Code Metrics

| File | Lines | Purpose |
|------|-------|---------|
| useFlowEditor.ts | ~716 | Builder state & operations |
| taq-parser.ts | ~474 | API ↔ VueFlow conversion |
| Builder.vue | ~438 | Main editor UI |
| NodePicker.vue | ~197 | Node selection dialog |
| AddableEdge.vue | ~169 | Custom edge component |
| ConfigSidebar.vue | ~145 | Node configuration panel |
| automation.ts | ~181 | Pinia store |

Total estimated LOC: ~2,500 lines
