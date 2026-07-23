---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - client/web/unify/src/sections/admin/components/Permissions/CPermissionGrid.vue
touched-by: []
tests: []
---

# Permissions Index view

## Intention

Give admins the system-component permission matrix: grant/deny operations on
system resources per role.

## UX capabilities

- Renders the shared `CPermissionGrid` with `component="system"` and the
  injected `$SystemAPI` — the grid owns all matrix behavior (role columns,
  rule editing, saving).
- This view contributes only the topbar title and the wiring.

## Routes

- `system.permissions` → `/system/permissions`.
- Sibling views under `views/compose|automation|federation/Permissions/`
  wrap the same grid for their components.

## When changing this

- Keep this a thin wrapper: matrix fixes belong in `CPermissionGrid`
  (`sections/admin/components/Permissions/`), where every component's
  permission view inherits them.
- Anything system-specific must go through props, not a forked grid.
