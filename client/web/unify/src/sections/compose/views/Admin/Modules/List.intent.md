---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
  - lib/vue/src/stores/useModuleStore.js
  - client/web/unify/src/sections/compose/components/Modules/ModuleImporter.vue
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Module List view

## Intention

Browse the namespace's modules and reach their editors; create, import, export,
and delete modules, and open their permission dialogs.

## UX capabilities

- Paginated, searchable, sortable table (name/handle/changedAt) scoped to the namespace; rows open the editor only when the module is updatable or deletable.
- Federated modules (label `federation`) carry a badge.
- Header: create (`canCreateModule`), export-all JSON (`canExportModules`), importer (re-fetches on success), wildcard permissions (`compose/` `grant`).
- Row menu: per-module permissions (`canGrant`), export, delete (confirm, store delete, refetch).

## Routes

`admin.modules` at `admin/modules` under `namespace.view`; navigates to `admin.modules.create` and `admin.modules.edit` (`:moduleID`).

## When changing this

- Deleting through `moduleStore` (not raw API) keeps the sidebar and page views consistent.
- Export JSON shape `{ type: 'module', list: [...] }` must stay importable.
