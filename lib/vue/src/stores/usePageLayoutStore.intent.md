---
kind: file
covers: usePageLayoutStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/compose
tests: []
---

# usePageLayoutStore

## Intention

Compose page-layout cache for the active namespace — the layout variants a
page can render (per role/condition), used by page views and the builder.

## State owned

`set` — frozen `compose.PageLayout` list (weight-ordered) for `namespaceID`.

## API surface consumed

`$ComposeAPI.pageLayoutListNamespace/List/Read/Create/Update/Delete`.

## Consumers

Compose page views (layout selection), page builder/admin, and the module
admin editor (creates a default layout for a new module's record page).

## Invariants

- `load` fetches ALL layouts of a namespace in one call
  (`pageLayoutListNamespace`); `findByPageID` is cache-first and falls back to
  the per-page list endpoint.
- `getByPageID` returns every layout of a page in weight order — selection
  logic depends on that ordering.
- Cached layouts are frozen; reads/mutations return fresh `compose.PageLayout`
  copies.
