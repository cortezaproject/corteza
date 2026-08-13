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

- `CTopbarSearch.vue` — search trigger + dialog: cancellable query via
  `$DiscoveryAPI`, recent-searches history, grouped results. It also owns
  record navigation, resolving the hit's namespace/page via `$ComposeAPI`.
- `items/ItemGroup.vue` / `items/RecordItem.vue` — internal result rendering
  (not exported); purely presentational — RecordItem only emits
  `click`/`open-new-tab` and the parent resolves the target.

## Contracts consumers rely on

- Requires `$DiscoveryAPI` injection — apps without the discovery service
  configured must not render this component.
- All labels come via the `labels` prop (with English defaults), including
  function-valued entries (`noResults`) — no i18n inside.
- Search fires on Enter or a recent-search click only. There is deliberately
  no search-as-you-type: discovery queries are too expensive per keystroke.
- Record navigation depends on vue-router being present (compose-style routes
  with namespace slug params).

## When changing this

Search is fired against whatever discovery indexed — result-shape changes on
the discovery side (highlight structure, resource kinds) surface here first.
