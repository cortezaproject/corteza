---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests:
  - client/web/unify/e2e/sections/admin/users.spec.ts
---

# User List view

## Intention

Find user accounts, act on their lifecycle (suspend/unsuspend/delete)
without opening the editor, and reach per-user permissions.

## UX capabilities

- Search + status filter (Active / Suspended / Deleted, popover); sort/paginate
  via `useResourceList`. A suspended or deleted user is marked by the shared
  row-state tag in the last-change cell; the list has no state column of its own.
- The last column is the shared "Last change" column (`changedAtField`): the most recent of deletedAt/updatedAt/createdAt, sorted by the same COALESCE expression.
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
