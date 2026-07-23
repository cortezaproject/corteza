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

# User List view

## Intention

Find user accounts, act on their lifecycle (suspend/unsuspend/delete)
without opening the editor, and reach per-user permissions.

## UX capabilities

- Search + tri-state suspended AND deleted filters (popover); sort/paginate
  via `useResourceList`; state column badges: suspended / deleted / enabled.
- Create button gated by `user.create`; wildcard permissions (`user/*`)
  gated by system `grant`.
- Row action menu: permissions (row `canGrant`), suspend/unsuspend
  (`canUpdateUser`), delete (`canDeleteUser`).

## Routes

- `system.users` → `/system/users`.
- Create → `system.users.create`; row click → `system.users.edit` with
  `userID`, only when the row grants `canUpdateUser` or `canDeleteUser`.

## When changing this

- Lifecycle actions also update the shared user store
  (`storeUsers`/`removeUsers`) so cached labels elsewhere stay fresh.
- Suspension does not end the user's active sessions; revoke lives in the
  editor.
