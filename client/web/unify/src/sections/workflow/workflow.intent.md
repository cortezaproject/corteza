---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/index.js
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/index.js
tests: []
---

# Workflow section

## Intention

The legacy automation/workflow webapp folded into the unify shell as one section:
list workflows, then build them on a visual canvas (VueFlow) that stays
byte-compatible with the legacy mxGraph model stored server-side. This doc governs
the section entry (`index.js`); every subfolder carries its own doc.

## Contract (section values declared by index.js)

- `id: 'workflow'`; every route tagged `meta.section: 'workflow'`.
- Unlike sections with a separate `routes.js`, this section declares fully-prefixed paths inline in `index.js`, mirroring the legacy app 1:1 under `/workflow`. Legacy route names were already namespaced (`workflow.*`) and are kept verbatim; the legacy `root` redirect became the section index route `workflow` → `workflow.list`.
- Routes: `workflow` (`/workflow`, redirect), `workflow.list` (`/workflow/list`, `meta.hideSidebar: true` — list keeps the sidebar collapsed like the legacy app), `workflow.create` (`/workflow/new`), `workflow.edit` (`/workflow/:workflowID/edit`); the last two share `views/Editor.vue`.
- `sidebar: WorkflowSidebar`; no topbar/profileItems/agentContext overrides.

## Map

- `index.js` — section contract: routes + sidebar wiring.
- `views/` — route targets: `Home.vue` (list), `Editor.vue` (canvas editor); own sidecars.
- `components/` — editor building blocks, `WorkflowEditor.vue` foremost (own doc; `Configurator/`, `FlowNodes/` nested docs).
- `composables/` — stateful canvas behaviors: history, DnD, clipboard, highlight (own doc).
- `lib/` — pure logic: graph↔API codec, connection rules, toolbar/style tables (own doc).
- `sidebar/` — left nav with per-workflow deep links (own doc).
- `stores/` — labels store: namespace/module display-name cache (own doc).

## When changing this

- Route names are the API: sidebar, list rows, and save-redirects push by name.
- Server model compatibility (steps/paths/triggers with mxGraph visual meta) is owned by `lib/codec.js` — any new step kind must land there, in the toolbar/style tables, and in a Configurator panel together.
