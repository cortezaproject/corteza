---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - client/web/unify/src/sections/compose/composables/usePageVisibility.ts
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
  - lib/vue/src/stores/useRecordStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Public page views

## Intention

The end-user face of compose: render pages built in the page builder, and view,
create, and edit records on record pages. These screens consume what `Admin/`
produces — they never define pages, only display them.

## Data touched

- `usePageStore` / `usePageLayoutStore` — pages and layouts are read from the
  stores preloaded by `Namespace/View.vue`; no direct page fetches here.
- `useRecordStore` — record CRUD plus `paginationRecordIDs` for prev/next.
- `usePageVisibility` — resolves the active layout and evaluates block/layout
  visibility expressions and roles (uses `$SystemAPI` sink for evaluation).
- `$ComposeAPI` — record attachment uploads on save.

## Map

- `Index.vue` — namespace landing: onboarding or redirect to the home page.
- `View.vue` — renders a non-record page's block grid.
- `RecordView.vue` — record page in view/create/edit mode; also embedded by
  `RecordModal` (see its sidecar).

## When changing this

- Block positioning contract: layout blocks intersected with page block
  definitions; a block present in the layout but missing from the page is
  dropped. `meta.hidden` blocks stay in the list (Grid/tabs need them) while
  expression/role-invisible blocks are removed entirely. Keep `View.vue`,
  `RecordView.vue`, and the builder's mirror logic consistent.
