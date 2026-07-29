---
kind: file
covers: backlogItems.js
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/events.js
  - client/web/unify/src/sections/project/stores/users.js
  - client/web/unify/src/sections/project/stores/dateUtils.js
touched-by:
  - client/web/unify/src/sections/project/views/dashboard
  - client/web/unify/src/sections/project/components/dashboard
tests: []
---

# backlogItems store

## Intention

Dashboard backlog: "sub-issues" linked to one category event (category +
eventID), backed by the ProjectBacklogItem system resource. Deliberately
mirrors the events store's shape and idioms so the two read as one layer.

## State owned

- `items` — the active project's backlog items, wholesale replaced per
  `load(projectId, revisionId)`; stable `id`, assignee resolved to a display
  name (raw ID kept as `assigneeId`). Same revision scoping as events.js:
  omit for the whole project, pass one for a single revision's board.
- `byEvent(category, eventID)` — items linked to one category event.
- `openCount` — open (per events' `isOpenStatus`) items project-wide; feeds
  the Backlog nav item's live badge.
- `mutations` — write-invalidation counter mirroring events.js#mutations
  (see that doc); bumped on every successful add/update/updateStatus/remove.

## API surface consumed

`$SystemAPI.projectBacklogItem{List,Create,Update,Delete}`; users store for
assignee names.

## Consumers

BacklogView, CategoryView/event detail (sub-issue lists), BacklogItemDialog,
DashboardNav (badge).

## Invariants

- Same 200-row load cap and local in-place mutation idiom as events.js —
  lists and badges react without a refetch. The cap is a known ceiling for a
  board spanning all six item types; revisit when it bites.
- `updateStatus` is optimistic with rollback, mirroring events.js.
- `assignee` is a display name on read, a user ID (or blank) on write;
  `dateDue` normalized to YYYY-MM-DD via `dateUtils.js`.
- The open rule is imported from events.js, never redefined.
