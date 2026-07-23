---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - client/web/unify/src/sections/admin/components/UserGroup/UserGroupMembers.vue
  - client/web/unify/src/sections/admin/components/UserGroup/UserGroupRoles.vue
  - lib/js
  - lib/vue
touched-by: []
tests: []
---

# UserGroup Editor view

## Intention

Create or edit one user group: identity, hierarchy placement, and (in edit
mode) its members and attached roles.

## UX capabilities

- Name (`meta.short`, required), handle, description; parent hierarchy as a
  list of `{selfID, name}` in `config.path` — hidden for the root group.
- Edit-only panels: members and roles, both self-managed sub-components
  keyed by `userGroupID` (they persist independently of the Save button).
- Delete and undelete (restore); per-group permissions button;
  unsaved-changes guard on leave.

## Routes

- `system.userGroups.create` → `/system/user-groups/new`;
  `system.userGroups.edit` → `/system/user-groups/:userGroupID`.
- Create redirects to `.edit`; back action → `system.userGroups`.

## When changing this

- Save payload is `{ handle, meta, config }` (+ `userGroupID` on update) —
  hierarchy changes travel inside `config.path`.
- Default-parent lookup fetches groups flat with `limit: 100`; revisit if
  group counts grow. Groups are not roles — no permission rules here.
