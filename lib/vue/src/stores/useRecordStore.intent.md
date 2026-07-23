---
kind: file
covers: useRecordStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
  - lib/vue/src/stores/useModuleStore.js
touched-by:
  - lib/vue/src/components/field/viewers/CFieldRecordViewer.vue
  - lib/vue/src/components/input/CInputRecord.vue
tests:
  - client/web/unify/src/sections/compose/stores/record.test.ts
  - lib/vue/src/components/field/viewers/CFieldRecordViewer.test.ts
---

# useRecordStore

## Intention

Shared compose record CRUD plus a cheap label-resolution path so record
references render as titles anywhere without each viewer fetching.

## State owned

Two caches: `records` Map of full `compose.Record` objects (module must be in
the module store) and `labelCache` Map of viewer-grade records (raw API shape
when the module isn't loaded). Also `paginationRecordIDs` for prev/next
navigation on record pages.

## API surface consumed

`$ComposeAPI.recordList/Read/Create/Update/Delete`.

## Consumers

Compose record views/blocks, record field viewers/editors, `CInputRecord`.

## Invariants

- `findByID` THROWS when the module isn't in the module store — viewers must
  use `resolveRecordLabels({namespaceID, moduleID, recordIDs})` instead.
- `getByID` checks both `records` and `labelCache`.
- `resolveRecordLabels` coalesces same-tick requests into one query per
  (namespace, module), dedups in-flight IDs, falls back to per-ID reads for
  unresolved IDs, and awaits until every requested ID settles.
- `findByID({ force: true })` bypasses the cache — required when opening a
  record for edit; supports AbortSignal for cancel-on-navigate.
