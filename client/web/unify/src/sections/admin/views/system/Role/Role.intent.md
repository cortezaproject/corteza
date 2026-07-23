---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - client/web/unify/src/sections/admin/components/Role/RolePermissionClone.vue
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# Role

## Intention

Manage RBAC roles: their identity, membership, lifecycle (archive/delete),
and context-role rules. Roles are the unit permissions attach to.

## Data touched

- `$SystemAPI.role*` incl. `roleMemberList/Add/Remove`,
  `roleArchive/Unarchive`, `roleDelete/Undelete`. Class-based `system.Role`.

## Map

- `List.vue` — role list with deleted/archived filters (see sidecar).
- `Editor.vue` — role form, members, context rules, clone (see sidecar).

## When changing this

- Role changes do not reach already-active sessions (server caches the
  user's roles until re-login) — don't "fix" that in the UI.
- Shared member/clone UI lives in `sections/admin/components/Role/`.
