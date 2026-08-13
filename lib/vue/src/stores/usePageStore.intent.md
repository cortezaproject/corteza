---
kind: file
covers: usePageStore.js
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/compose
tests: []
---

# usePageStore

## Intention

Compose page cache for the active namespace: drives compose navigation and
the page admin/builder, including tree structure operations.

## State owned

`set` — frozen `compose.Page` list (weight-ordered) for `namespaceID`.

## API surface consumed

`$ComposeAPI.pageList/Read/Create/Update/Delete/Tree/Reorder`.

## Consumers

Compose sidebar navigation, page views, page admin/builder, and
`CFieldRecordViewer` (reads `set` directly to build record links).

## Invariants

- Delete with `cascade`/`rebase` strategy refetches the whole namespace —
  server-side child deletion/reparenting makes the local tree stale.
- `reparent` is implemented as `pageUpdate` with a new `selfID` plus a sibling
  `reorder` — NOT via pageReorder alone (pageReorder breaks for nested
  parents server-side).
- Cached pages are frozen; `findByID` / mutations return fresh `compose.Page`
  copies.
