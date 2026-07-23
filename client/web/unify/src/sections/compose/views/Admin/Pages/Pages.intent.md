---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
  - client/web/unify/src/sections/compose/components/PageBlocks
  - client/web/unify/src/sections/compose/components/Admin/Page/PageTranslator.vue
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Pages admin

## Intention

Author the namespace's pages: organize the page tree, edit page metadata and
layouts, and compose each layout's blocks visually in the builder. What is built
here is what `../../Pages/` renders for end users.

## Data touched

- `usePageStore` — tree load/reorder, page CRUD; mutations flow through the
  store so the sidebar and public views update immediately.
- `usePageLayoutStore` — per-page layouts: CRUD, reorder
  (`$ComposeAPI.pageLayoutReorder`), builder saves.
- `$ComposeAPI` — page icons (iconList/iconUpload/iconDelete/pageUpdateIcon),
  page translations.
- RBAC: namespace `canCreatePage`, per-page `canUpdatePage`/`canDeletePage`/
  `canGrant`; permission resource `corteza::compose:page/<namespaceID>/<pageID>`.

## Map

- `List.vue` — drag-and-drop page tree (order + parenting).
- `Edit.vue` — page metadata, icon, and layout list/config.
- `Builder.vue` — per-layout visual block builder (the core contract — see its
  sidecar).

## When changing this

- Page/layout/block model: blocks are defined on the page; each layout
  references a subset by blockID with its own positions. Builder, Edit, and the
  public views all assume this — never let a save drop page blocks that other
  layouts reference.
- New pages (create, clone, record-page helpers in Modules) always get a seeded
  `primary` layout built from the `compose.PageLayout` type; a plain object
  persists all record-toolbar buttons disabled.
