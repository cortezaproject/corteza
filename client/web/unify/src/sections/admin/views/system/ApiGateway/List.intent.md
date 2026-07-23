---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# ApiGateway List view

## Intention

Overview of all integration gateway routes plus the global gateway settings
(profiler mode, proxy redirect-following) — the entry point for creating,
editing, and profiling routes.

## UX capabilities

- Settings panel: profiler tri-state (disabled / filter / global) and proxy
  follow-redirects toggle, persisted only on explicit Save.
- Paged/sortable/filterable route list (`useResourceList`) with
  deleted-state filter popover; row click opens the editor; New creates;
  wildcard and per-route permissions; delete via row actions.
- Profiler button appears only when the profiler is not disabled.

## Routes

- `system.apiGateway` → `/system/api-gateway`; navigates to
  `system.apiGateway.create`, `.edit` (`:routeID`), and `.profiler`.

## When changing this

- The profiler tri-state maps to two settings keys, `apigw.profiler.enabled`
  and `apigw.profiler.global`; keep the derivation consistent both ways.
- Keep the profiler entry button gated on the setting — the profiler pages
  are useless (empty) while disabled.
