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

# UserGroup

## Intention

Manage user groups — hierarchical groupings of users used for organizing
and scoping (e.g. auth-client visibility), separate from RBAC roles.

## Data touched

- `$SystemAPI.userGroup*` incl. `userGroupUndelete`; class-based resource
  `system.UserGroup`.

## Map

- `List.vue` — group list with deleted/archived filters (see sidecar).
- `Editor.vue` — group form, hierarchy, members/roles panels (see sidecar).

## When changing this

- Groups are not roles: nothing here should touch permission rules.
- Members/roles panel components live in
  `sections/admin/components/UserGroup/`.
