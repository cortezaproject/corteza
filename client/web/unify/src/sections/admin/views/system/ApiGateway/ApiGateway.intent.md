---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# ApiGateway

## Intention

Manage integration gateway routes — custom API endpoints composed of filter
chains (prefilter/processer/postfilter) — plus the request profiler
(`Profiler/`) for inspecting gateway traffic.

## Data touched

- `$SystemAPI.apigwRoute*` (plain raw objects — no resource class),
  `apigwFilterDefFilter` / `apigwFilter*`, `apigwProfiler*`; profiler and
  proxy enablement live in system settings (`settingsList`/`settingsUpdate`).

## Map

- `List.vue` — route list + global gateway settings (sidecar).
- `Editor.vue` — route create/edit + filter chain assembly (sidecar).
- `Profiler/Index.vue` — aggregated per-path stats (sidecar).
- `Profiler/Route.vue` — hit list for one profiled path (sidecar).
- `Profiler/Hit.vue` — single request/response detail (sidecar).

## When changing this

- Route order in `routes.js` is load-bearing: `/profiler` paths must stay
  declared before the `/:routeID` catch, or "profiler" is read as an ID.
- Profiler route params carry the base64-encoded path, not a server ID —
  the encode/decode contract spans all three profiler views.
