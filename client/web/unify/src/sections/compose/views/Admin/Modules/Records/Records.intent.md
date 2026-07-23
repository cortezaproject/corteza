---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/useRecordStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Admin record views

## Intention

Record CRUD for a module without needing any built page: every module gets an
"all records" table, a create form, and a view/edit screen for free. All three
render synthetic page blocks through the same `Grid` used by public pages, so
admin and public record UIs stay behaviorally identical.

## Data touched

- `useModuleStore` — module resolved from the preloaded store by `:moduleID`
  route param; a missing module renders "not found" (no fetch fallback).
- `useRecordStore` — findByID/create/update/delete.
- `$ComposeAPI` — record attachment upload endpoint on save.

## Map

- `List.vue` — synthetic RecordList block, admin-routed row navigation.
- `Create.vue` — synthetic Record block in create mode.
- `Edit.vue` — view + edit modes for one record (two routes, one component).

## When changing this

- The synthetic-block pattern is the contract: block `options.fields: []` means
  "all fields", and a synthetic `page` object with `pageID: '0'` is passed so
  blocks never resolve a real page.
- Create/Edit provide `recordViewContext` and `$fileUploadContext` exactly like
  `Pages/RecordView.vue` — RecordBlock depends on both; keep the shapes in sync.
- `compose.Record` throws on a module without fields — create is blocked (with
  a message) until the module has at least one field.
