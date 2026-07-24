---
kind: file
covers: events.js
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/users.js
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/stores/dateUtils.js
touched-by:
  - client/web/unify/src/sections/project/views/dashboard
  - client/web/unify/src/sections/project/stores/backlogItems.js
tests: []
---

# events store

## Intention

Dashboard (Manage & Monitor) event data: one flat list over the five
per-project category resources (incident, feature, privacy, task, review),
so the category views, KPIs and nav badges all read a single source.

## State owned

- `events` — the active project's events across all categories, wholesale
  replaced per `load(projectId)`; each row carries a stable `id`, its
  `category`, and owner refs resolved to display names (raw IDs kept as
  `<key>Id` for round-tripping).
- Derived getters: `byCategory`, `countByCategory`, `kpis`
  (total/open/overdue), `breakdown` (bar-chart buckets), `ownerOptions`
  (project members + access users, falling back to the full directory).
- `isOpenStatus` (exported) — THE open rule: open ⇔ status !== 'Completed'.
  Shared with backlogItems and the Overview report cards; never re-derive.

## API surface consumed

`$SystemAPI.project{Incident,Feature,Privacy,Task,Review}{List,Create,
Update,Delete}` — one resource per category, bound via the CATS table.
Loads users + project user sets first so owner refs resolve.

## Consumers

Dashboard views (Overview, CategoryView, AllEventsView), NewEventDialog,
DashboardNav badges, backlogItems (reuses `isOpenStatus`).

## Invariants

- Loads cap at 200 rows per category (accepted for v1); aggregate-correct
  numbers come from the report store instead.
- Mutations patch the local list in place — lists/KPIs/badges must react
  without a refetch.
- Payload dates go through `normalizeDates` (YYYY-MM-DD); owner fields are
  sent as user IDs, shown as names.
