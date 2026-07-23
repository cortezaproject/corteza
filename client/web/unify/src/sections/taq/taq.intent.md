---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAutomationStore.js
  - lib/js/src/automation/types/taq.ts
touched-by:
  - client/web/unify/src/sections/index.js
  - client/web/unify/src/sections/admin/views/automation/Taq/Editor.vue
tests: []
---

# TAQ section

## Intention

The TAQ (Trigger Action Query — never "Task Queue") section is the visual
automation builder of the unified app. A TAQ is an ng-automation resource:
triggers fire it, function steps act, gateways/iterators shape the flow. This
section owns listing TAQs and graph-editing them on a canvas; the admin
section's automation/Taq area owns metadata CRUD and deep-links into the
builder here (`/taq/builder/<automationID>`, opened in a new tab).

## Map

- `index.js` — section registration: id `taq`, sidebar component, three routes (below).
- `views/` — route targets `List.vue` and `Builder.vue`; each carries its own sidecar.
- `components/` — `flow/` (canvas node/edge renderers), `builder/` (picker, config sidebar, forms, trace), `common/` (config modal, icon).
- `composables/` — `useFlowEditor` (the builder's entire model) and `useSegmentForm`.
- `sidebar/TaqSidebar.vue` — drawer nav listing all automations (see `sidebar/sidebar.intent.md`).
- `utils/` — API↔VueFlow conversion, dagre layout, icon/dimension constants.

## Data touched

All persistence goes through `$AutomationAPI` `ngAutomation*` endpoints
(list/read/create/update/exec/execution-trace). Shared cache/catalog state
lives in `useAutomationStore` (automation list + functions/triggers catalog).
RBAC resource: `corteza::automation:ng-automation/<id>`.

## When changing this

- Routes mirror the legacy standalone app 1:1: `taq` at `/taq` (list,
  `meta.hideSidebar`), `taq.builder` at `/taq/builder` (new), and
  `taq.builder-edit` at `/taq/builder/:id`. Names are namespaced `taq.*` to
  avoid collisions with other sections — renaming breaks the admin deep-link
  and the sidebar's named-route navigation.
- The shell does global setup only: this section lazy-loads its own catalog
  (Builder) and automation list (sidebar); never assume either is preloaded.
