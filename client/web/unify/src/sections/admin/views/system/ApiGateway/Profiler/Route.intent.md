---
kind: file
covers: Route.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# ApiGateway Profiler Route view

## Intention

List the recorded hits for one profiled gateway path so an admin can find
the specific request worth inspecting.

## UX capabilities

- Hit table (hitID, status, timestamp), newest first; the decoded path is
  shown in the topbar title.
- "Load older" cursor pagination in pages of 50; the footer button hides
  once a short page signals the end.
- Back button returns to the aggregate profiler; row click (or row action)
  opens the hit detail.

## Routes

- `system.apiGateway.profiler.route` →
  `/system/api-gateway/profiler/:routeID` (`routeID` = base64 path).
- Navigates to `system.apiGateway.profiler.hit` with `routeID` + `hitID`.

## When changing this

- `routeID` is passed to `apigwProfilerRoute` in its base64 form — decode is
  for display only; keep pairing with `Index.vue`'s `btoa(path)`.
- Cursor = `before` set to the last row's `hitID`; `hasMore` is inferred
  from a full page, so a page-size change must keep that heuristic.
