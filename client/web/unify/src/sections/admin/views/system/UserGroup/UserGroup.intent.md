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

Manage user groups — a hierarchy that places every user in exactly one group,
used for organizing and scoping (e.g. auth-client visibility). Roles attached
to a group are live permissions: a member is evaluated against the roles of
their own group and of every group beneath it, and a change to a group's
members or roles applies without a server restart.

## Data touched

- `$SystemAPI.userGroup*` incl. `userGroupUndelete`; class-based resource
  `system.UserGroup`.

## Map

- `List.vue` — group list with deleted/archived filters (see sidecar).
- `Editor.vue` — group form, hierarchy, members/roles panels (see sidecar).

## When changing this

- Groups are not roles: nothing here should touch permission rules. Moving a
  user into a group higher in the tree widens what they can reach.
- Members/roles panel components live in
  `sections/admin/components/UserGroup/`.
