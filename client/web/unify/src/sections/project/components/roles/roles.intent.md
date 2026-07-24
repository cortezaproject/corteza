---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/stores/users.js
  - client/web/unify/src/sections/project/components/permissions
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
  - client/web/unify/src/sections/project/components/wizard/steps/RolesStep.vue
tests: []
---

# Project role dialogs (members model)

## Intention

The role slice of the locked members model: project access roles are what
users hold and what all permissions attach to. Implements the resource dialog
standard for the `role` kind — create a role, then manage its meta, its
members, and its full permission set in one detail dialog. Opened via the
Wizard's `createResource('role')` / `inspectResource('role', id)`.

## Data touched

- `useProjectsStore`: `rolesFor`, `addRole`, `updateRole`,
  `setProjectUserRole` (membership writes), `projectUsersFor`, plus the
  resource loaders the embedded matrix needs (resources, pages, automations,
  agents, chatbots, connections, project users).
- `useProjectUsersStore`: platform directory for member names/emails.
- Deep link: `system.roles.edit` (full admin role editor, new tab).

## Map

- `RoleCreateDialog.vue` — minimal name + description create; emits
  `created(id)`.
- `RoleDetailDialog.vue` — mixed persistence model: role meta is staged
  (committed on Save), while membership changes persist immediately through
  `RoleMemberList`. Embeds the single-role permissions view:
  `ProjectPermissionMatrix` with `roles=[role]` + `hideRoleHeader` (all
  capabilities per row, no columns). Reload-after-open never clobbers
  in-progress meta edits.
- `RoleMemberList.vue` — role membership editor shared with RolesStep: pick
  any platform user via `CInputUser` (adds instantly), remove instantly;
  members derived from the project-user list filtered to this role; a local
  pick-cache renders users the 500-user directory hasn't loaded.

## When changing this

- Membership writes go through `setProjectUserRole` only — the same primitive
  the user dialog uses; keep both sides consistent.
- Preserve the staged-meta vs immediate-members split; a single Save that also
  commits membership would change the members-model UX contract.
