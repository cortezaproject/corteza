---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# Role List view

## Intention

Browse and manage RBAC roles, and reach the role editor or role permissions.

## UX capabilities

- Search + tri-state deleted AND archived filters (popover); sort/paginate
  via `useResourceList`; description shown under the name.
- Create button gated by `role.create`; wildcard permissions button
  (`role/*`) gated by system `grant`.
- Row action menu: per-role permissions (row `canGrant`), delete
  (`canDeleteRole`).

## Routes

- `system.roles` → `/system/roles`.
- Create → `system.roles.create`; row click → `system.roles.edit` with
  `roleID` — but only when the row grants `canUpdateRole` or
  `canDeleteRole`; otherwise the click is a no-op.

## When changing this

- Keep the row-click permission gate: users who can only read a role must
  not land in an editor they cannot save.
- Deleting here only refreshes the list; it does not touch active sessions —
  role changes reach users on re-login (server session cache).
