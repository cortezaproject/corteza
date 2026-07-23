---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/ApiGateway/CFilterParamsEditor.vue
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# ApiGateway Editor view

## Intention

Create or edit one gateway route (endpoint, method, enabled, async, meta)
and assemble its filter chain across prefilter/processer/postfilter steps.

## UX capabilities

- Endpoint is required; create redirects to edit after first save.
- Filters panel exists only in edit mode: per-step tabs, add from
  server-defined kinds (used ones disabled), drag reorder (weight),
  per-filter config modal (`CFilterParamsEditor`), remove with confirm.
- Delete route, per-route permissions, unsaved guard covering route fields
  and pending filter edits; back lands on `system.apiGateway`.

## Routes

- `system.apiGateway.create` → `/system/api-gateway/new`;
  `system.apiGateway.edit` → `/system/api-gateway/:routeID`.

## When changing this

- Filters save as a diff: only entries flagged created/updated/deleted hit
  `apigwFilter*`, only after the route update succeeds — keep that
  reconciliation or orphaned filters accumulate server-side.
- The `response` filter's params are shape-converted FE↔BE; filter identity
  in the editor is `ref` (definition name) — no `filterID` until saved.
