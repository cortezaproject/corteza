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
tests: []
---

# User Editor view

## Intention

Administer one user account end to end: identity, group, credentials, role
memberships, avatar, external auth links, and lifecycle.

## UX capabilities

- Create: email (required) / name / handle / user group (default group
  preselected). Edit adds panels: security, roles, avatar, external auth.
- Password is set only when filled (admin override, write-only); role
  memberships save as a diff via `roleMemberAdd/Remove`.
- Suspend/unsuspend, revoke all sessions (disabled for yourself), delete;
  unsaved-changes guard on leave.

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
