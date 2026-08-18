---
kind: file
covers: View.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue
  - client/web/unify/src/sections/compose/composables/usePageVisibility.ts
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests:
  - client/web/unify/src/sections/compose/views/Pages/View.title.test.js
---

# Page View (public, non-record)

## Intention

Render a compose page for end users: pick the layout whose condition/roles match
the current user and screen, title the screen the way that layout asks, position
the page's blocks accordingly, and hide blocks whose visibility rules fail.

## UX capabilities

- Blocks render through the shared `Grid`; layout selection honours an explicit `?layoutID=` (e.g. from a navigation block) when that layout's own condition/roles pass, else falls back to normal order; with no matching layout the page's own block positions are used as fallback.
- Editors (`canUpdatePage`) get topbar shortcuts to the page builder and page editor; a translator button appears per resource-translation settings.
- Distinct empty states for "page has no blocks" and "page not found".
- The topbar shows the page title, unless the chosen layout sets `config.useTitle` and its own `meta.title` — that title is interpolated against the signed-in user (`${user.name}`, `${userID}`). There is no record here, so a title reading one, an empty one, or one that will not evaluate leaves the page title standing.

## Routes

`page` at `pages/:pageID` under `namespace.view`. Reloads on `pageID` change. Links to `admin.pages.builder` and `admin.pages.edit`.

## When changing this

- The page comes from `pageStore.getByID` only — no API fallback; it relies on `Namespace/View.vue` preloading. A pageID missing from the store shows "not found".
- Visibility is evaluated before content shows (no flash); expression/role-hidden blocks are removed, while `meta.hidden` (tab children) must stay in the block list for `Grid`/TabsBlock.
- Layout/positioning logic mirrors `RecordView.vue` and the builder — change all together. So does the title override, minus the record: keep the two fallback chains the same shape.
