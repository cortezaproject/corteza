---
kind: file
covers: Overview.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/report.js
  - client/web/unify/src/sections/project/config/categories.js
touched-by: []
tests: []
---

# Overview view

## Intention

Landing dashboard of a live project: per-category health at a glance plus the
project's activity pulse, driven entirely by the server-side report endpoint
(grouped counts — accurate beyond any list row cap).

## UX capabilities

- One card per governance category (open count + status doughnut with total),
  each linking to that category's page.
- Audit-events activity panel (admin-only; hides itself when unavailable).
- Cross-category created-over-time trend, windowed by a time-range control
  (default 6 months); changing the range reloads only the trend.
- Skeletons on first load; a failed report shows a visible retry instead of
  silently rendering zeros.

## Routes

- `project.overview`, the index child of the dashboard layout; links to
  `project.overview.category` per card.

## When changing this

- Keep the separate staleness guards for cards vs trend (a range change must
  not orphan an in-flight project load, and vice versa).
- "Open" means not-Completed — the single definition lives in
  `stores/events.js#isOpenStatus`; don't restate the rule here.
