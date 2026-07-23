---
kind: file
covers: RecordView.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - client/web/unify/src/sections/compose/composables/usePageVisibility.ts
  - lib/vue/src/stores/useRecordStore.js
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
  - client/web/unify/src/sections/compose/components/Record/RecordModal.vue
tests: []
---

# Record View (public record page)

## Intention

The record page in all three modes — view, edit (`?edit=1`), create
(`recordID === '0'`) — rendered from the page's block layout. The same component
also runs inside `RecordModal` (`inModal` + `modalPageID`/`modalRecordID` props,
navigation via query params instead of route params).

## UX capabilities

- Record toolbar: back/cancel, delete, save-as-copy, new, edit, save — each individually toggleable via the active layout's `config.buttons`, further gated by record `can*` flags.
- Prev/next navigation across the record set last listed by a RecordList block (`recordStore.paginationRecordIDs`).
- Create supports cloning (`?cloneFromID`) and reference prefill (`?refField`/`?refValue`, multi-aware); `ownedBy` defaults to the current user.
- Save flushes pending file uploads (collected from field editors via the provided `$fileUploadContext`) before create/update; server-side field errors map back onto the form; required fields validated client-side.
- View↔edit switches in place on a pristine clone — cancel restores without refetch; unsaved-changes guard on route leave; block visibility re-evaluated on mode/record change.
- Provides `recordViewContext` (mode/record/isNew/isSaving) that RecordBlock consumes.

## Routes

`page.record` at `pages/:pageID/records/:recordID` under `namespace.view`. Non-record pages redirect to `page`. Editors get topbar links to module edit, builder, and page edit.

## When changing this

- Record loads are abortable; a superseded load must not clobber newer state.
- Fast-swap (loadRecord only) applies solely between two existing records on the same page — any transition involving create mode must go through full `loadPage()`.
- After save/clone the query params `cloneFromID`/`refField`/`refValue` must be stripped so later actions don't inherit them.
- Listens to the `refetch-records` event-bus signal (automation buttons) to reload.
