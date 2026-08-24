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
- `DryRunField.vue` — one input row in the dry-run/test dialog, typed from the workflow's declared input.
- `Help.vue` — static controls/shortcuts reference dialog.
- `C3.js`, `WorkflowEditor.c3.js` — component-catalog scenarios.

## Data touched

- WorkflowEditor reads `$AutomationAPI` (functionList, eventTypesList, workflowExec, sessionRead/Cancel), `$SystemAPI`/`$ComposeAPI` for runAs and dry-run lookups, the injected `$eventBus` for realtime session messages, and the labels store. Persistence stays with the parent view via `save`/`delete`/`undelete` emits; graph state serializes through `lib/codec`.

## When changing this

- Configurators mutate `item.config` / `edges[id].config` in place; deep watchers in WorkflowEditor bridge those maps back into VueFlow nodes/edges. Breaking the bridge silently loses step config at save time.
- Connection legality lives in `lib/connectionRules` (shared with nodes and tests); parallel-gateway `ref` (fork/join) is auto-derived from edge balance; out-edge labels of excl gateway / iterator / error-handler are positional and renumbered on delete.
- `eventBus` (mitt) carries `change-detected` / `trigger-updated` from nested configurators — dirty tracking depends on it.
- A dry run watches its session while it is running and stops the moment it is not. `prompted` (waiting on a person) and `suspended` (waiting on a delay) are live, resumable states the server never completes on its own, so the poll renders the partial trace and stops; the realtime `workflowSessionResumed` message re-attaches it and the trace carries on filling in. Because that message names only a state, the `workflowSessionPrompt` messages are read for the session each state belongs to. A session that neither finishes nor suspends is dropped at a wall-clock ceiling, the only bound a runaway workflow has.
- The poll owns the spinner and the stop control, and must not outlive either: leaving the editor and stopping the test both abort it, the second without waiting on a cancel the server can refuse.
- Running a dry run and reading its trace are separate permissions: `execute` on the workflow starts it, while the session poll additionally needs `sessions.search` on the automation component and `sessions.manage` on the workflow (`lib/dry-run`'s `canReadTrace`). Without both the run still happens and the poll is skipped, so the toast says the trace is not visible rather than that the test failed; only the workflow's own reported error is a failed test.
