---
kind: file
covers: Overview.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/project/components/dashboard/OverviewPanel.vue
touched-by:
  - client/web/unify/src/sections/project/index.js
tests: []
---

# Overview view

## Intention

Route target for the dashboard's landing screen, chain-wide. A thin wrapper:
the screen itself is `components/dashboard/OverviewPanel.vue`, the same
component the wizard's Manage & Monitor tab mounts (via `ManageOverview.vue`)
scoped to one revision.

## UX capabilities

All of them live in `OverviewPanel`: per-category health cards (open count +
status doughnut) linking to each category's page, the admin-only audit-events
activity panel, and the cross-category created-over-time trend with its
time-range control. Counts come from the server-side report endpoint, so they
stay accurate beyond any list row cap.

## Routes

- `project.overview`, the index child of the dashboard layout; links to
  `project.overview.category` per card.

## When changing this

- Keep it thin. Overview behaviour belongs in `OverviewPanel`, so the
  dashboard and the wizard cannot drift — that is the point of the panel.
- "Open" means not-Completed — the single definition lives in
  `stores/events.js#isOpenStatus`; don't restate the rule here.
