---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/workflow/lib
  - client/web/unify/src/sections/workflow/composables
  - lib/vue
touched-by:
  - client/web/unify/src/sections/workflow/views/Editor.vue
  - client/web/unify/src/sections/workflow/views/Home.vue
tests: []
---

# Workflow editor components

## Intention

Building blocks of the visual workflow editor. `WorkflowEditor.vue` is the heart:
a VueFlow canvas replacing the legacy mxGraph editor with feature parity —
everything else here serves it or the list view. `Configurator/` and `FlowNodes/`
carry their own docs.

## Map

- `WorkflowEditor.vue` — the editor: drag palette (toolbar strip), VueFlow canvas, right sidebar hosting the per-step Configurator, workflow-config dialog, dry-run/test with trace replay onto the graph, undo/redo, clipboard, keyboard shortcuts (Ctrl+S/Z/C/X/V/A, delete, arrow-nudge, Shift+?), Ctrl+wheel zoom, swimlane reparenting on drop.
- `FlowEdge.vue` — custom edge: double-click inline label editing, highlight/trace styling.
- `Import.vue` / `Export.vue` — JSON workflow exchange (workflows + triggers), shared by list and configurator.
- `NamespaceModuleSelector.vue` — multi-select of `ref_namespace`/`ref_module` label values; used by list filter and workflow configurator.
- `ExpressionEditor.vue` / `ExpressionTable.vue` — expression input and reorderable target/type/expr rows for configurators.
- `Help.vue` — static controls/shortcuts reference dialog.
- `C3.js`, `WorkflowEditor.c3.js` — component-catalog scenarios.

## Data touched

- WorkflowEditor reads `$AutomationAPI` (functionList, eventTypesList, workflowExec, sessionRead/Cancel), `$SystemAPI`/`$ComposeAPI` for runAs and dry-run lookups, and the labels store. Persistence stays with the parent view via `save`/`delete`/`undelete` emits; graph state serializes through `lib/codec`.

## When changing this

- Configurators mutate `item.config` / `edges[id].config` in place; deep watchers in WorkflowEditor bridge those maps back into VueFlow nodes/edges. Breaking the bridge silently loses step config at save time.
- Connection legality lives in `lib/connectionRules` (shared with nodes and tests); parallel-gateway `ref` (fork/join) is auto-derived from edge balance; out-edge labels of excl gateway / iterator / error-handler are positional and renumbered on delete.
- `eventBus` (mitt) carries `change-detected` / `trigger-updated` from nested configurators — dirty tracking depends on it.
