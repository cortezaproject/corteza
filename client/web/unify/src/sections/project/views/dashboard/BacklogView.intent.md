---
kind: file
covers: BacklogView.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/backlogItems.js
  - client/web/unify/src/sections/project/stores/events.js
touched-by: []
tests: []
---

# BacklogView view

## Intention

Flat cross-category list of follow-up work items — backlog items are their
own linked resource (not tags on category records) — with KPIs and charts for
triage.

## UX capabilities

- KPI trio and status/category/priority charts are computed client-side from
  the loaded items (no report endpoint; the trend chart is deliberately
  omitted because the row cap would silently truncate it).
- Create and row-edit share one dialog; row click opens a read-only drawer
  (shared right-sidebar store, exclusive with other right panels) that links
  through to the parent event's own edit dialog.
- Empty state carries a create CTA; search/sort are client-side, with
  priority/status sorted by canonical rank.

## Routes

- `project.overview.backlog` at `backlog` under the dashboard layout.

## When changing this

- Open/overdue rules are shared: `isOpenStatus` from `stores/events.js`;
  `dateDue` is a plain YYYY-MM-DD string, compared as a string.
- A backlog row resolves its linked event via the events store (loaded by
  DashboardLayout); a missing event must degrade to the raw `#id` reference.
