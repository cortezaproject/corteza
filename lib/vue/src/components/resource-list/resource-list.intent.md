---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
  - lib/vue/src/components/input/CInputSearch.vue
touched-by:
  - client/web/unify/src/sections/admin
  - client/web/unify/src/sections/chatbot/views/List.vue
tests:
  - lib/vue/src/composables/useResourceList.test.ts
---

# Resource list

## Intention

The one list-screen shell: every admin/chatbot List view is `useResourceList`
(state + fetching) wired into `CResourceList.vue` (card, search, lazy DataTable,
cursor pagination, row actions). New list screens follow this pairing instead of
composing a DataTable by hand.

## Contracts

- CResourceList is fully controlled and lazy: parent supplies `items`, `filter`,
  `sorting`, `pagination`, `loading`, `fields` ({ key, header/label, sortable,
  frozen, pt, ... }) and reacts to `update:filter`, `sort`, `row-click`,
  `page-change` ({ pageCursor, page, limit? }). The component fetches nothing.
- Pagination is cursor-based: `pagination.prevPage`/`nextPage` cursors gate the
  buttons; empty cursor means first page. Per-page select emits a reset page-change.
- Search box binds `filter[queryField]` (default `query`); hide via `hideSearch`.
- `actionItems(row) => MenuItem[]` auto-appends a frozen, hover-revealed actions
  column with one centralized TieredMenu (recreated per open so anchoring never
  goes stale); a `fields` entry keyed `actions` is then dropped. Items support
  `route` for router-links. `hideActionsMenu()` is exposed.
- Cell content via `body-<field.key>` slots; also `header`, `filter`, `footer`,
  `expansion` slots. `translations` prop carries all user-facing strings.
- useResourceList(apiFn, options): `apiFn` gets encoded params ({ limit, sort,
  pageCursor, incTotal, ...filter }) and must return { response, cancel }
  (abortable). Returns readonly `items`/`loading`/`error`, reactive
  `filter`/`sorting`/`pagination`, and handlers (`fetchItems`, `filterList`,
  `handleSort`, `handlePageChange`, `abortRequests`) that map 1:1 onto
  CResourceList's props/events. State round-trips through the route query so
  lists are deep-linkable; in-flight requests are cancelled on refetch.
- Row navigation is NOT part of the pairing: every list screen binds its own
  `@row-click`, because the target route is per-resource. (The composable also
  returns a `handleRowClick`, but it is hardcoded to `namespace.edit` and no
  screen uses it.)

## When changing this

Every admin list screen inherits changes here — verify one cursor-paginated and
one filtered screen after touching pagination or filter emit shapes.
