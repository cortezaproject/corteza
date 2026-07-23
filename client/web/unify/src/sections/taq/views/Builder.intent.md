---
kind: file
covers: Builder.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/taq/composables/useFlowEditor.ts
  - client/web/unify/src/sections/taq/utils/taq-parser.ts
  - lib/vue/src/stores/useAutomationStore.js
touched-by:
  - client/web/unify/src/sections/taq/index.js
  - client/web/unify/src/sections/admin/views/automation/Taq/Editor.vue
tests: []
---

# TAQ Builder view

## Intention

The graph editor for one TAQ: compose triggers, action steps, branches and
iterators into an executable flow on a VueFlow canvas, then run it in place
and inspect the execution trace. This screen owns the whole
trigger/action/query building UX; all state and persistence logic lives in
`useFlowEditor` — the view wires canvas, panels and dialogs together.

## UX capabilities

- Auto-laid-out canvas (dagre): trigger, step, branch, iterator, end and loop nodes; every edge carries a "+" to insert a node mid-flow; empty state prompts for the first trigger. Layout is always automatic by design — nodes are never draggable or manually positioned; the graph's shape derives solely from flow structure.
- Node picker dialog (step palette): triggers and function steps grouped from the `useAutomationStore` catalog, searchable; also used to replace an existing node in place (trigger↔trigger, step↔any step kind).
- Selecting a node opens the resizable config sidebar (segment-driven forms, constraints, gateway type, branch conditions, metadata) and highlights the ancestor path; clicking an argument's link opens the reference panel to bind upstream outputs (`{scope, source}`) to arguments or condition sides — including condition rows addressed as `condition:<edgeId>:<side>:<rowIndex>`.
- Run: dirty flow forces save-first dialog; trigger properties (or agent-trigger input schema) collected via RunModal, else executes directly. Trace overlay shows status/duration banner, per-node executed tint, traversed edges, and a per-step I/O TracePanel for the selected node.
- Toolbar: back, zoom in/out/fit (also middle-mouse double-click), undo/redo, save (Ctrl/Cmd+S). Delete/Backspace removes selection (with confirm from node menus); unsaved-changes route guard.
- Header controls: enabled toggle, Run As user, per-TAQ permissions button, "show all configurations" preview toggle; topbar title click opens `TaqConfigModal` to edit name/description/labels/runAs, which saves immediately by design — metadata is independent of the graph's dirty/save cycle.

## Data touched

`$AutomationAPI` `ngAutomationRead/Create/Update/Exec/ExecutionTrace` (through
`useFlowEditor`); `useAutomationStore` catalog (`loadCatalog` triggered here —
the shell does not preload it); `useRightSidebarStore` (`taq-config` panel key)
so only one right panel is open shell-wide.

## Routes

`taq.builder` at `/taq/builder` (new TAQ) and `taq.builder-edit` at
`/taq/builder/:id`. Loading waits for `catalogReady`; `:id` is watched so
sidebar navigation between TAQs reloads in place. Saving a new TAQ
`router.replace`s to its id URL. The admin automation/Taq editor deep-links
here in a new tab; back returns to `/taq`.

## When changing this

- The catalog must be loaded before parsing a flow — icons, forms and
  upstream results all resolve against it.
- End nodes are real termination steps on save; they are non-selectable and
  must never open the config sidebar or be deletable directly.
- `provide('activeReferenceArgument')` is the contract deep inputs
  (DynamicInput, CInputFieldValueMap, condition rows) rely on for active-link
  state — keep it when restructuring.
- Selection is VueFlow-owned but node data must be re-looked-up from
  `editor.nodes` (VueFlow selection refs go stale after updates).
