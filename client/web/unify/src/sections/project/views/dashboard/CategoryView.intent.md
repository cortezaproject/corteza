---
kind: file
covers: CategoryView.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/categories.js
  - client/web/unify/src/sections/project/stores/events.js
  - client/web/unify/src/sections/project/stores/report.js
touched-by: []
tests: []
---

# CategoryView view

## Intention

One governance category's page: a report-driven metrics band (KPIs,
breakdown charts, windowed trend) above the category's item list, with full
create / inspect / edit / delete of records and their linked backlog items.

## UX capabilities

- KPI trio (total/open/overdue) and breakdowns come from the report endpoint,
  so they stay accurate above the events store's 200-row list cap; the list
  itself is client-side filtered/sorted over the store's rows.
- Row click opens a read-only detail drawer (registered with the shared
  right-sidebar store — exclusive with other right panels); its Edit button,
  or the row kebab, opens the edit dialog. Deletes confirm first.
- The create dialog can queue backlog titles that become backlog items linked
  to the new record (event creation succeeding is not undone by their failure).
- Ranked enums (severity/risk/status) sort and chart in canonical order.

## Routes

- `project.overview.category` at `category/:category`; an unknown category key
  renders a muted empty note. Linked from Overview's cards and DashboardNav.

## When changing this

- Mutations refresh metrics silently (no skeleton flash); keep the split
  staleness guards (cards vs trend) when touching the load flow.
- Config-driven: columns, charts, KPIs and form schema all come from
  `config/categories.js` — extend there, not with view-local branches.
