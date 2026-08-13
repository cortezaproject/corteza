---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# ApiGateway Profiler Index view

## Intention

Monitor gateway traffic at a glance: aggregated per-path statistics (hit
count, payload size, timing) with drill-down into individual routes.

## UX capabilities

- Aggregate table; sizes rendered in kB, times in ms.
- Row click (or row action) drills into that path's hit list.
- Purge-all button clears collected profiler data (shown only when there is
  data).
- Auto-refresh every 10 s with a visible countdown on the refresh button;
  manual refresh restarts the cycle.

## Routes

- `system.apiGateway.profiler` → `/system/api-gateway/profiler`.
- Drill-down navigates to `system.apiGateway.profiler.route` with
  `routeID = btoa(path)`.

## When changing this

- The `routeID` param is the base64-encoded request path, not a server ID —
  keep the encode here paired with the decode in `Route.vue`.
- The refresh timer is cleared during load and on unmount; keep it that way
  to avoid overlapping fetch loops after navigation.

> **DRIFT:** the columns declare `sortable: true` but the table binds no
> `@sort` handler and the aggregation call sends no sort params, so clicking
> a header does nothing. Either wire sorting or drop the flags.
