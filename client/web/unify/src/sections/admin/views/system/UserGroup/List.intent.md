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

# UserGroup List view

## Intention

Browse user groups and reach the group editor, permissions, or delete.

## UX capabilities

- Search + status filter (Active / Archived / Deleted, popover); sort/paginate
  via `useResourceList`; name column shows `meta.short` with description.
- The last column is the shared "Last change" column (`changedAtField`): the most recent of deletedAt/updatedAt/createdAt, sorted by the same COALESCE expression.
- Create button gated by `user-group.create`; wildcard permissions
  (`user-group/*`) gated by system `grant`.
- Row action menu: per-group permissions (row `canGrant`), delete
  (`canDeleteUserGroup`). Delete counts the group's members and child groups
  first and opens `UserGroupDeleteBlocked` instead of the confirm while either
  remains.

## Routes

- `system.userGroups` → `/system/user-groups`.
- Create → `system.userGroups.create`; row click → `system.userGroups.edit`
  with `userGroupID`, only when the row grants `canUpdateUserGroup` or
  `canDeleteUserGroup`.

## When changing this

- Display name is `meta.short` (groups have no top-level name); keep the
  `meta.short || handle || userGroupID` fallback order consistent.
- The RBAC resource is `user-group` (hyphenated) while route names use
  `userGroups` — don't "normalize" either.
