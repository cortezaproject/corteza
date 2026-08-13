---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/usePageStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Page List (tree) view

## Intention

Show the namespace's pages as the hierarchy end users navigate, and let admins
restructure it: reorder siblings and re-parent pages by drag and drop.

## UX capabilities

- Tree from `pageStore.loadTree`, sorted by weight, all parents expanded initially; root-level pages (`selfID === '0'`) are visually distinguishable from nested ones.
- Client-side lenient filtering by title; clicking a node opens the page editor.
- Drag and drop persists the full new structure: re-parented pages get their `selfID` updated first, then each level is reordered, recursively; afterwards both the tree and the flat page list are re-fetched so the sidebar reflects the change.
- Create button gated by `canCreatePage`.

## Routes

`admin.pages` at `admin/pages` under `namespace.view`; navigates to `admin.pages.create` and `admin.pages.edit` (`:pageID`).

## When changing this

- Reorder must run update-selfID before `pageStore.reorder` per level — reordering against stale parents scrambles weights.
- The refetch after a successful drop is what keeps the sidebar in step; keep it.

> **DRIFT:** a FAILED drop is never corrected. The refetch sits inside the same
> `try` as the reorder, so a failed persist jumps to `catch` (log + toast) and
> the optimistic tree stays on screen showing an order the server rejected,
> until the user reloads.
