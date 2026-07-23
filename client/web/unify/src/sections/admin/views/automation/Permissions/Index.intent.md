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

# Automation Permissions Index view

## Intention

Component-wide RBAC matrix for the automation component (workflow/TAQ service
and wildcard resource rules).

## UX capabilities

- Mounts the shared `CPermissionGrid` with `component="automation"` and the injected `$AutomationAPI` — all grid behavior (role columns, rule editing, saving) lives in that component.
- Sets the topbar title; nothing else is view-specific.

## Routes

`automation.permissions` at `/automation/permissions`; no params.

## When changing this

- Behavior changes belong in the shared `CPermissionGrid`, not here — this wrapper must stay identical in shape to the compose/federation ones.
- Per-resource (single workflow/TAQ) permissions are handled by `CPermissionsButton` in the resource views, not by this grid.
