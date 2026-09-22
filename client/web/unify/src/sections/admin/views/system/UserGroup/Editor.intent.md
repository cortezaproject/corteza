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
  (they persist independently of the Save button). The members panel takes
  the group itself — it names the group in its move confirmations and gates
  on `canManageMembersOnUserGroup`; the roles panel takes `userGroupID`.
- Delete and undelete (restore); per-group permissions button;
  unsaved-changes guard on leave. Delete counts members and child groups first
  and opens `UserGroupDeleteBlocked` instead of the confirm while either
  remains.
- Corredor manual scripts for `system:user-group` bound to ui page
  `user-group/editor`, slot `infoFooter`, render as `CManualScriptButtons`
  under the basic information panel (edit only); a click dispatches the
  script on `$ScriptBus` with the group as the event's subject.

## Routes

- `system.userGroups.create` → `/system/user-groups/new`;
  `system.userGroups.edit` → `/system/user-groups/:userGroupID`.
- Create redirects to `.edit`; back action → `system.userGroups`.

## When changing this

- Save payload is `{ handle, meta, config }` (+ `userGroupID` on update) —
  hierarchy changes travel inside `config.path`.
- Default-parent lookup fetches groups flat with `limit: 100`; revisit if
  group counts grow. Groups are not roles — no permission rules here.
