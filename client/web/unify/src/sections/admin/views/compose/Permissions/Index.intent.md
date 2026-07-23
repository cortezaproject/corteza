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

# Compose Permissions Index view

## Intention

Component-wide RBAC matrix for the compose component (namespaces, modules,
records and the compose service itself, as wildcard resource rules).

## UX capabilities

- Mounts the shared `CPermissionGrid` with `component="compose"` and the injected `$ComposeAPI` — all grid behavior (role columns, rule editing, saving) lives in that component.
- Sets the topbar title; nothing else is view-specific.

## Routes

`compose.permissions` at `/compose/permissions`; no params.

## When changing this

- Behavior changes belong in the shared `CPermissionGrid`, not here — this wrapper must stay identical in shape to the automation/federation ones.
- Per-resource compose permissions are managed inside the compose section, not from this grid.
