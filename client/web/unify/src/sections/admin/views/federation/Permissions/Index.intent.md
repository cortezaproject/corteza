---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/components/Permissions/CPermissionGrid.vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Federation Permissions Index view

## Intention

Component-wide RBAC matrix for the federation component (node service and
wildcard resource rules).

## UX capabilities

- Mounts the shared `CPermissionGrid` with `component="federation"` and the injected `$FederationAPI` — all grid behavior (role columns, rule editing, saving) lives in that component.
- Sets the topbar title; nothing else is view-specific.

## Routes

`federation.permissions` at `/federation/permissions`; no params.

## When changing this

- Behavior changes belong in the shared `CPermissionGrid`, not here — this wrapper must stay identical in shape to the automation/compose ones.
- Per-node permissions are handled by `CPermissionsButton` in the node views, not by this grid.
