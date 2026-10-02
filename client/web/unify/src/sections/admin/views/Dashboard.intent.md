---
kind: file
covers: Dashboard.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/components/chart/CChart.vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests:
  - client/web/unify/e2e/sections/admin/corredor-admin-slots.spec.ts
---

# Dashboard view

## Intention

Admin landing page: at-a-glance instance health over a chosen range — how many
of each resource exist and in what state, how activity, sign-ins and automation
runs moved over time, and what needs attention.

## UX capabilities

- One aggregate call feeds the page: `GET /system/stats/?from&to&bucket`. The range presets (7d, 30d, 90d, 1y) pick the bucket (day, week, month) and the choice persists in localStorage `admin.dashboard.range`. Weeks start on Monday and days are taken in the server's zone (`server/system/service/statistics.go`).
- Sixteen resource tiles in a four-column grid, each with a status bar and sparkline; a tile opens `ResourceDialog` on `GET /system/stats/{resource}` (status counts, created/updated/deleted per day, newest rows with labels).
- Activity, sign-ins and one automation-runs card (workflow and TAQ charts) are `TrendChart`s on CChart. A click anywhere in a chart column opens `BucketDialog` on `GET /system/stats/events/{kind}` (activity, signins, workflows, taqs); any listed event opens `EventDialog`, with actor names via `useActors`.
- An attention list of recent failures, labelled through `eventLabel.js`, which reuses the action-log vocabulary so a resource/action pair reads as a sentence.
- A section the caller may not read (`*.search` on its component, or `action-log.read` for activity and sign-ins) is omitted from the payload and nothing is drawn for it; errors never block rendering. Sets the topbar title via the `#topbar-title` teleport.
- Corredor manual scripts for `system` bound to ui page `dashboard`, slot `toolbar`, render as `CManualScriptButtons` above the tiles; a click dispatches a `system` event on `$ScriptBus`.

## Routes

`dashboard` at `/dashboard` — target of the `admin` section-index redirect; no params.

## When changing this

- The pieces live in `sections/admin/components/Dashboard/`: `useSystemStats` owns the range, the request and the bucket labels, `useDashboardData` the derived shapes, `resources.js` the tile list. A new tile needs an entry there and a resource in `statistics.go`.
- TAQ run counts come from one action-log row per terminal run (`automation:ng-automation` / `run`, written by `recordNgAutomationRun`); the exec ledger itself is in-memory, so runs before those rows existed are not counted.
- Chart colours derive from live PrimeVue theme CSS vars (`chartTheme.js`) so theming keeps working. `TrendChart` maps a click to a column with `convertFromPixel({ gridIndex: 0 })`; the axis-index finder and vue-echarts' `@zr:click` return nothing here.
