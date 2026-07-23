---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - lib/vue/src/stores/useModuleStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Admin Record List view

## Intention

Show every record of a module in a full-page table, regardless of whether the
namespace has a record-list page built — the admin's raw data view.

## UX capabilities

- Renders a single synthetic `RecordList` block (all fields, 20 per page, selectable, export allowed) through the shared `Grid`, so the table behaves exactly like a RecordList block on a public page.
- Row and create navigation is redirected to admin record routes by providing a `$recordRoutes` resolver (view / create with optional `cloneFromID`).
- Topbar links back to the module editor; unknown moduleID shows "not found".

## Routes

`admin.modules.record.list` at `admin/modules/:moduleID/records` under `namespace.view`; reached from the module editor and record screens; rows go to `admin.modules.record.view`, create to `admin.modules.record.create`.

## When changing this

- `$recordRoutes` is the override point RecordListBlock checks before using public page routes — removing it silently reroutes rows to public record pages.
- The module must already be in the module store (preloaded by `Namespace/View.vue`).
