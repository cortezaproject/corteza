---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by: []
tests: []
---

# search/

## Intention

Global topbar search over the discovery service: a dialog that queries
indexed resources (compose records and friends), highlights hits, and
navigates to the matched resource.

## Map

- `CTopbarSearch.vue` — search trigger + dialog: debounced/cancellable query
  via `$DiscoveryAPI`, recent-searches history, grouped results.
- `items/ItemGroup.vue` / `items/RecordItem.vue` — internal result rendering
  (not exported); RecordItem resolves the record's namespace/page via
  `$ComposeAPI` to build the navigation target.

## Contracts consumers rely on

- Requires `$DiscoveryAPI` injection — apps without the discovery service
  configured must not render this component.
- All labels come via the `labels` prop (with English defaults), including
  function-valued entries (`noResults`, `numberOfResults`) — no i18n inside.
- Record navigation depends on vue-router being present (compose-style routes
  with namespace slug params).

## When changing this

Search is fired against whatever discovery indexed — result-shape changes on
the discovery side (highlight structure, resource kinds) surface here first.
