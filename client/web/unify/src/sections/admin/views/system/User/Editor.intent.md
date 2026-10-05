---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - client/web/unify/src/sections/admin/components/User/UserSecurity.vue
  - client/web/unify/src/sections/admin/components/User/UserRoles.vue
  - client/web/unify/src/sections/admin/components/User/UserAvatar.vue
  - client/web/unify/src/sections/admin/components/User/UserExternalAuth.vue
  - lib/js
  - lib/vue
touched-by: []
tests:
  - client/web/unify/e2e/sections/admin/corredor-admin-slots.spec.ts
  - client/web/unify/e2e/sections/admin/user-password.spec.ts
---

# User Editor view

## Intention

Administer one user account end to end: identity, group, credentials, role
memberships, avatar, external auth links, and lifecycle.

## UX capabilities

- Create: email (required) / name / handle / user group (default group
  preselected). Edit adds panels: multi-factor authentication, roles, avatar,
  sign-in methods (the user's credentials, password included).
- The edit action row, ordered by consequence and coloured to match: a
  Suspended/Deleted tag in the user list's colours, Set password (info),
  revoke all sessions (warn; disabled with a tooltip on your own account),
  suspend (danger) or unsuspend (success), permissions.
- Set password opens `UserPasswordDialog`: write-only, both entries must
  match, and it is stored at once through `userSetPassword`, outside the
  form's Save and the unsaved-changes guard. Hidden for system users and for
  anyone who cannot update the user.
- Role memberships save as a diff via `roleMemberAdd/Remove`; delete;
  unsaved-changes guard on leave.
- Corredor manual scripts for `system:user` bound to ui page `user/editor`
  render as `CManualScriptButtons` in two slots: `infoFooter` under the basic
  information panel, `passwordFooter` under the multi-factor authentication
  panel's toggles (raised by `UserSecurity` as `script`); both edit only, and a click
  dispatches the script on `$ScriptBus` with the user as the event's subject.

## Routes

- `system.users.create` → `/system/users/new`; `system.users.edit` →
  `/system/users/:userID`. Create redirects to `.edit`; back → `system.users`.

## When changing this

- `meta` must stay in the update payload — MFA toggles live there.
- Saves also push into the shared user store; role/flag changes do not reach
  active sessions — "revoke all sessions" is the immediate lever.
- On create, assign the returned user before navigating so child panels
  never fetch with an empty userID.
- `initialMembershipIDs` is the role diff the save sends, not the
  unsaved-changes baseline — the guard tracks memberships through its own
  `extra`. Removing it as guard bookkeeping would silently stop role changes
  being applied.
