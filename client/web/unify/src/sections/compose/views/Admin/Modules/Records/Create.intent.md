---
kind: file
covers: Create.vue
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

# Admin Record Create view

## Intention

Create a record for a module directly from the admin area, using a synthetic
all-fields Record block — same form behavior as a public record page in create
mode, minus the page.

## UX capabilities

- New records default `ownedBy` to the current user; `?cloneFromID` prefills values from a source record; `?refField`/`?refValue` prefill a reference field (multi-aware).
- Required-field validation client-side; server field errors mapped back onto the form; pending file uploads flushed before create.
- Save replaces to the record's admin view; cancel replaces to the admin record list.
- Modules without fields show an explanatory message instead of a form (a record cannot be constructed).
- Topbar links to the module editor and the admin record list.

## Routes

`admin.modules.record.create` at `admin/modules/:moduleID/records/create` under `namespace.view`. Re-initializes when moduleID or the clone/prefill query params change.

## When changing this

- Keep `recordViewContext` / `$fileUploadContext` provides in lockstep with `Pages/RecordView.vue` — RecordBlock and field editors consume them.
- A failed clone-source load degrades to a blank record rather than blocking creation.
