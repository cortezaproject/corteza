---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
  - lib/js
touched-by: []
tests: []
---

# User

## Intention

Administer user accounts: identity, credentials, role membership, and
lifecycle (suspend/delete). The primary identity surface of the platform.

## Data touched

- `$SystemAPI.user*` incl. `userSetPassword`, `userSuspend/Unsuspend`,
  `userSessionsRemove`, `userMembershipList`; `roleMember*`,
  `userGroupList`. Class-based resource `system.User`.
- Shared user store (`useUserStore`) is kept in sync on mutations.

## Map

- `List.vue` — user list with lifecycle actions (see sidecar).
- `Editor.vue` — full account editor with sub-panels (see sidecar).

## When changing this

- Password fields must remain write-only and optional on update.
- Role/flag changes do not affect the user's active sessions (server-side
  session cache) — session revocation is the lever for immediate effect.
- Sub-panels (security/roles/avatar/external auth) live in
  `sections/admin/components/User/`.
