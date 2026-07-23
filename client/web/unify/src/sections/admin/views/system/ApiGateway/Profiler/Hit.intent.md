---
kind: file
covers: Hit.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# ApiGateway Profiler Hit view

## Intention

Inspect a single profiled gateway request: what came in, what went out, and
whether it succeeded — the terminal drill-down of the profiler.

## UX capabilities

- Shows hit ID, status (tag colored success for 2xx, warn otherwise), and
  timestamp.
- Request and response rendered as pretty-printed JSON in scrollable blocks.
- Back navigates via browser history (-1); a not-found message renders when
  the hit no longer exists (e.g. purged).

## Routes

- `system.apiGateway.profiler.hit` →
  `/system/api-gateway/profiler/:routeID/hit/:hitID`.
- Reached from `system.apiGateway.profiler.route`; links nowhere further.

## When changing this

- Only `hitID` is used for fetching (`apigwProfilerHit`); `routeID` exists
  in the URL purely for hierarchy/back-navigation context.
- Profiler data is ephemeral — always keep the missing-hit fallback state.
