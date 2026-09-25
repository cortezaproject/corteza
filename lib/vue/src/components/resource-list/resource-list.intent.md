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
- Row state belongs to the shell, not the screen: the `changedAt` cell ends
  with a Deleted / Suspended / Archived tag after the date (`resourceState`, the most final
  one wins), and a deleted row's text is muted. A list does not render state
  itself.
- A list shows one status at a time: Active, or exactly one of the states its
  resource has (`states`: deleted, suspended, archived, disabled).
  `CResourceStatusFilter` is the one picker, a dropdown inside the list's
  filter popover. A status maps onto the API's per-state '0'/'1'/'2' filters
  (`statusFilter`): Active excludes every state; a state is "only" that one and
  excludes the rest, except Deleted, which keeps its rows whatever else they
  carry. There are no Without / Including / Only choices for states.
- Active filters show as a bar of removable chips above the table. A status
  other than Active is one chip, "Status" plus the row-state tag. Any other
  filter a screen labels in `filterLabels` (compared against `filterDefaults`
  from `useResourceList`) is a chip too: `{ key: label }` for a '0'/'1'/'2'
  filter that is not a state, `{ key: { label, value } }` for anything else.
  Removing a chip, or Reset filter, emits `update:filter` with Active and the
  defaults. The search box is never a chip.
- Cell content via `body-<field.key>` slots; also `header`, `filter`, `footer`,
  `expansion` slots. `translations` prop carries all user-facing strings.
- useResourceList(apiFn, options): `apiFn` gets encoded params ({ limit, sort,
  pageCursor, incTotal, ...filter }) and must return { response, cancel }
  (abortable). Returns readonly `items`/`loading`/`error`, reactive
  `filter`/`sorting`/`pagination`, the frozen `filterDefaults`, and handlers (`fetchItems`, `filterList`,
  `handleSort`, `handlePageChange`, `abortRequests`) that map 1:1 onto
  CResourceList's props/events. State round-trips through the route query so
  lists are deep-linkable; in-flight requests are cancelled on refetch.
- Screens must not call `filterList` from a radio's `@change`: the route write
  it makes lands after the next change and restores the old values; the filter
  watcher already refetches.
- Row navigation is NOT part of the pairing: every list screen binds its own
  `@row-click`, because the target route is per-resource. (The composable also
  returns a `handleRowClick`, but it is hardcoded to `namespace.edit` and no
  screen uses it.)

## When changing this

Every admin list screen inherits changes here — verify one cursor-paginated and
one filtered screen after touching pagination or filter emit shapes.
