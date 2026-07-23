---
kind: file
covers: Home.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/workflow/stores/labels.js
  - client/web/unify/src/sections/workflow/components/Import.vue
  - client/web/unify/src/sections/workflow/components/Export.vue
  - client/web/unify/src/sections/workflow/components/NamespaceModuleSelector.vue
  - lib/vue
touched-by: []
tests: []
---

# Workflow list view

## Intention

Find workflows and act on their lifecycle (enable/disable, delete/undelete,
export/import) without opening the editor; entry point to the canvas.

## UX capabilities

- Search/sort/paginate via `useResourceList`; filter popover: tri-state
  subWorkflow / disabled / deleted radios + namespace/module label filter
  (`ref_namespace=[…]` / `ref_module=[…]` label queries); rows resolve
  namespace/module tags through the labels store.
- Create + Import gated by automation `workflow.create`; wildcard permissions
  gated by `grant`. Row actions: edit, enable/disable, single export
  (workflow + triggers JSON), delete with confirm, undelete.
- Import creates workflows then triggers, owned by current user, `runAs`
  reset to `'0'`; per-handle failures reported as skipped, not fatal.

## Data touched

- `$AutomationAPI` workflow + trigger CRUD; labels store; shared `workflowStore` (`updateInList`/`removeFromList`) keeps the sidebar fresh.

## Routes

- `workflow.list` → `/workflow/list` (sidebar hidden). Row click / edit →
  `workflow.edit` with `workflowID`; create button → `workflow.create`.

## When changing this

- Export JSON shape (handle/enabled/meta/keepSessions/steps/paths/triggers)
  must stay importable here and in the editor's Configurator import.
