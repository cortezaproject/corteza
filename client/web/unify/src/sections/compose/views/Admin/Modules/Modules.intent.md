---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useModuleStore.js
  - client/web/unify/src/sections/compose/components/Admin/Module
  - client/web/unify/src/sections/compose/components/ModuleFields/Configurator
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Modules admin

## Intention

Define and maintain the namespace's data model: module metadata, fields and
their configuration, storage (DAL) settings, and the bridge into related pages
and records. `Records/` (own folder doc) handles admin record CRUD.

## Data touched

- `useModuleStore` — CRUD keeps the shared set in sync for every other view.
- `usePageStore` / `usePageLayoutStore` — record page and record-list page
  creation from the module editor, including seeded default layouts.
- `$ComposeAPI` — module field translations; `$SystemAPI` —
  `dalSchemaAlterationList` for schema-issue resolution.
- `$Settings` — feature gates for federation and discovery settings modals.
- RBAC: namespace `canCreateModule`/`canExportModules`/`canManageNamespace`,
  per-module `can*`; permission resources `corteza::compose:module/...`,
  `module-field/...`, `record/...`.

## Map

- `List.vue` — paginated module list with import/export.
- `Edit.vue` — the module editor (fields, DAL, unique values, revisions, issues).
- `Records/` — admin record list/create/view+edit (see `Records/Records.intent.md`).

## When changing this

- Module export/import JSON shape `{ type: 'module', list: [...] }` must stay
  compatible with `ModuleImporter`.
- Record/record-list page creation seeds a `primary` layout from the page's
  blocks using the `compose.PageLayout` type — see `Edit.intent.md` traps.
