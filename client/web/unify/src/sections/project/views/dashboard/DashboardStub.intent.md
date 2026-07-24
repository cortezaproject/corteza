---
kind: file
covers: DashboardStub.vue
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/project/index.js
tests: []
---

# DashboardStub view

## Intention

Placeholder body for dashboard views whose real content is not built yet —
currently only Reports (pending the custom-reports design). Keeps the nav and
routing demonstrable.

## UX capabilities

- Renders the route's `meta.titleKey` and `meta.icon` with a "coming soon"
  note; no data, no interaction.

## Routes

- `project.overview.reports` at `reports` under the dashboard layout; any
  future stubbed child can point here too.

## When changing this

- Replace by swapping the route's component in `../index.js`, not by growing
  this stub; keep the `titleKey`/`icon` route-meta contract (DashboardNav
  reads it too).
