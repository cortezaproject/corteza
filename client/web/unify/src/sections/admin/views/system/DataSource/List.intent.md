---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# DataSource List view

## Intention

Show the built-in primary database alongside manageable external DAL
connections — see where records live, add or remove external sources.

## UX capabilities

- Primary panel (fetched separately by `primary-dal-connection` type):
  name, handle, location, ownership; edit pencil goes to the shared editor;
  never deletable, never in the paged list.
- External list fixed to type `dal-connection`: search, sort, paginate via
  `useResourceList`; deleted-state filter popover; New button; delete via
  row actions; row click opens the editor.

## Routes

- `system.dataSources` → `/system/data-sources`; navigates to
  `system.dataSources.create` and `.edit` (`:connectionID`) — the same
  editor serves both primary and external connections.

## When changing this

- The primary/external type split is the screen's core contract; keep the
  type filter hard-wired in the list fetch and the primary fetch separate.
- Not the integration-connections screen — same `:connectionID` param name,
  different API (`dalConnection*`).
