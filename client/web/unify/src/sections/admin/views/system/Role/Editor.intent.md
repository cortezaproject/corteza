---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - client/web/unify/src/sections/admin/components/Role/RoleMembers.vue
  - client/web/unify/src/sections/admin/components/Role/RolePermissionClone.vue
  - lib/js
  - lib/vue
touched-by: []
tests: []
---

# Role Editor view

## Intention

Define one role fully: identity, contextual (rule-driven) vs
membership-based behavior, members, and lifecycle.

## UX capabilities

- Edit name/handle/description; contextual toggle swaps the members panel
  for a context panel (expression + applicable resource types).
- Members managed via `RoleMembers`; saved as a diff (add/remove) on submit;
  saving a contextual role strips any leftover explicit members.
- Archive/unarchive, delete/undelete, clone permission rules to other roles
  (`RolePermissionClone` dialog), per-role permissions button.
- `isSystem` roles cannot be archived/deleted; `isClosed` roles are read-only.

## Routes

- `system.roles.create` → `/system/roles/new`; `system.roles.edit` →
  `/system/roles/:roleID`. Create redirects to `.edit`; back → `system.roles`.

## When changing this

- Contextual vs membership is exclusive — toggling clears the other side's
  state; preserve that invariant or roles end up with both.
- Membership sync is a batched diff after the role update; a partial failure
  leaves membership half-applied — keep the operations grouped.
- `initialMemberIDs` is that diff's other side, not the unsaved-changes
  baseline — the guard tracks membership through its own `extra`. Removing it
  as guard bookkeeping would silently stop the save sending member changes.
