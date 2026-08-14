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
tests:
  - client/web/unify/src/sections/compose/views/Pages/RecordView.layout.test.js
  - client/web/unify/src/sections/compose/views/Pages/RecordView.title.test.js
---

# Record View (public record page)

## Intention

The record page in all three modes — view, edit (`?edit=1`), create
(`recordID === '0'`) — rendered from the page's block layout. The same component
also runs inside `RecordModal` (`inModal` + `modalPageID`/`modalRecordID` props,
navigation via query params instead of route params).

## UX capabilities

- Record toolbar: back/cancel, delete, save-as-copy, new, edit, save — each individually toggleable via the active layout's `config.buttons`, further gated by record `can*` flags.
- Back leaves edit mode by dropping the `?edit` query; from view mode it returns to the previous screen, or the namespace's pages when the record was opened directly.
- Prev/next navigation across the record set last listed by a RecordList block (`recordStore.paginationRecordIDs`).
- Create supports cloning (`?cloneFromID`) and reference prefill (`?refField`/`?refValue`, multi-aware); `ownedBy` defaults to the current user.
- Save flushes pending file uploads (collected from field editors via the provided `$fileUploadContext`) before create/update; server-side field errors map back onto the form; required fields validated client-side.
- View↔edit switches in place on a pristine clone — cancel restores without refetch; unsaved-changes guard on route leave.
- `resolveLayout()` picks the layout once the record it needs is resolved, and re-runs on initial load, on a same-page record swap, on view/edit/create mode changes, and after save — but never on field-value edits, since a layout swap remounts the block set and would discard in-progress input. An explicit `?layoutID=` (e.g. from a navigation block) wins only if that layout's own condition/roles pass, else selection falls back to normal order.
- Block visibility, unlike layout, is reactive to record values: re-evaluated (debounced 300ms) on value edits as well as mode/record changes.
- Provides `recordViewContext` (mode/record/isNew/isSaving) that RecordBlock consumes.
- The displayed title is the page's, unless the active layout sets `config.useTitle` — then its `meta.title` is interpolated against the open record (`${record.values.x}`, `${recordID}`, `${ownerID}`, `${userID}`), so one page can title itself per layout. Any failure to evaluate falls back to the page title rather than surfacing a broken string.

## Routes

`page.record` at `pages/:pageID/records/:recordID` under `namespace.view`. Non-record pages redirect to `page`. Editors get topbar links to module edit, builder, and page edit.

## When changing this

- Record loads are abortable; a superseded load must not clobber newer state. Layout resolution and block-visibility evaluation carry the same guarantee via their own sequence counters.
- Fast-swap (loadRecord only) applies solely between two existing records on the same page — any transition involving create mode must go through full `loadPage()`; because it skips `loadPage()`, it must call `resolveLayout()` itself.
- The pageID/recordID/query-param watcher uses separate per-value getters, not one getter returning an array: an array is compared by identity and re-fires on every touched dependency, including the edit-query toggle it must ignore.
- After save/clone the query params `cloneFromID`/`refField`/`refValue` must be stripped so later actions don't inherit them.
- Listens to the `refetch-records` event-bus signal (automation buttons) to reload.
