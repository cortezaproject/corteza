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
tests: []
---

# Project user dialogs (members model)

## Intention

The user slice of the locked members model: **users hold per-project roles**.
These dialogs implement the resource dialog standard for the `user` kind —
adding a member (with roles, which are mandatory) and inspecting/editing a
member's role set plus their resulting effective access. Opened via the
Wizard's `createResource('user')` / `inspectResource('user', id)`.

## Data touched

- `useProjectsStore`: `projectUsersFor` (userId + held roleIds),
  `rolesFor`, `addProjectUser` (invite new by email),
  `assignProjectUserRoles`, `setProjectUserRole`, plus the resource loaders
  the embedded matrix charts (`loadResources/Pages/Agents/Chatbots`,
  `loadProjectUsers`).
- `useProjectUsersStore` (`stores/users.js`): platform user directory for
  names/emails (`load`, `reload`, `findUser`).
- Deep link: `system.users.edit` (full admin user editor, new tab).

## Map

- `UserCreateDialog.vue` — two modes: pick an existing platform user (users
  already on the project are excluded) or invite a new one by email; at least
  one project role is required before Add is enabled. Persists user (if new)
  then role assignments; emits `created(userId)`.
- `UserDetailDialog.vue` — everything applies **live** (footer is Close-only,
  no Save): role chips toggle membership optimistically via
  `setProjectUserRole` (rolled back on failure), and an embedded
  `ProjectPermissionMatrix` with `evalUserId` shows a read-only "Evaluated"
  column (the user's resolved access) beside editable columns for exactly the
  roles they hold — columns track chip toggles immediately.

## When changing this

- Roles are the only unit of access: never grant a user permissions directly;
  the evaluated column must stay read-only.
- `resourceId` here is the **userId** (project-user rows are keyed by user,
  not a separate membership id).
- Keep the seed/reload dance in the detail dialog: reload never clobbers
  membership edits made while loading.
