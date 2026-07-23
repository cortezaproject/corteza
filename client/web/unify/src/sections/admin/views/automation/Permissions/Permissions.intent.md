---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/Permissions/CPermissionGrid.vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Automation — component permissions

## Intention

Component-wide RBAC matrix for the automation component. Thin wrapper: all
grid behavior lives in the shared `CPermissionGrid` admin component; this
area only binds it to the automation API client.

## Map

- `Index.vue` — mounts `CPermissionGrid` (see `Index.intent.md`).

## Data touched

- `$AutomationAPI` permission endpoints, via `CPermissionGrid`.

## When changing this

- Behavior changes belong in the shared `CPermissionGrid`, not here — this
  wrapper must stay identical in shape to the compose/federation ones.
- Per-resource (single workflow/TAQ) permissions are handled by
  `CPermissionsButton` in the resource views, not by this grid.
