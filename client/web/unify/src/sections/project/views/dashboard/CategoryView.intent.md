---
kind: file
covers: CategoryView.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/project/components/dashboard/CategoryPanel.vue
touched-by:
  - client/web/unify/src/sections/project/index.js
tests: []
---

# CategoryView view

## Intention

Route target for one governance category's page, chain-wide. A thin wrapper:
the screen itself is `components/dashboard/CategoryPanel.vue`, the same
component the wizard's Manage & Monitor tab mounts scoped to one revision.
The view only turns the `:category` route param into a prop.

## UX capabilities

All of them live in `CategoryPanel`: the report-driven metrics band (KPI trio,
breakdowns, windowed trend) above the category's item list, the read-only
detail drawer registered with the shared right-sidebar store, create / edit /
delete of records, and the create dialog's queued backlog titles.

## Routes

- `project.overview.category` at `category/:category`, a child of the
  dashboard layout; an unknown category key renders a muted empty note.
  Linked from Overview's cards and DashboardNav.

## When changing this

- Keep it thin. Category behaviour belongs in `CategoryPanel`, so the
  dashboard and the wizard cannot drift — that is the point of the panel.
- Config-driven: columns, charts, KPIs and form schema all come from
  `config/categories.js` — extend there, not with view-local branches.
