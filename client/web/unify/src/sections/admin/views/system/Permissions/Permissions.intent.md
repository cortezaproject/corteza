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

# Permissions

## Intention

System-component permission matrix: grant/deny operations on system
resources per role. The view is a thin route shell.

## Data touched

- `$SystemAPI` passed to the shared `CPermissionGrid`
  (`sections/admin/components/Permissions/`) with `component="system"`.

## Map

- `Index.vue` — wrapper mounting `CPermissionGrid` (see sidecar).

## When changing this

- Keep this a wrapper. Matrix behavior/fixes belong in `CPermissionGrid`,
  where the sibling compose/automation/federation permission views also get
  them.
