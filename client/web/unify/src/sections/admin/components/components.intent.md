---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/views
tests: []
---

# Admin components

## Intention

Reusable editor sub-panels and dialogs that the admin views compose, one folder
per resource family (mirroring `views/`). Nothing here is routed directly; each
component is embedded by an Editor/Index view and receives its resource or IDs
via props.

## Map

- `User/` — user editor panels: UserAvatar (avatar preview/upload/remove), UserExternalAuth (the sign-in methods list: password and linked external providers), UserPasswordDialog (sets a password from the editor's action row, both entries matching, stored at once), UserRoles (role-membership editor), UserSecurity (MFA toggles; hosts the `user/editor` / `passwordFooter` Corredor script slot and re-emits a button click as `script` for the editor to dispatch).
- `Role/` — RoleMembers (role membership editor), RolePermissionClone (dialog that clones another role's permission rules; emits `cloned`).
- `UserGroup/` — UserGroupMembers, UserGroupRoles: group membership and group-role assignment panels. A user is in exactly one group, so UserGroupMembers never removes anyone: a row moves its user to another group, and adding a user from another group moves them here. Every move is confirmed with a sentence naming both groups, and adding and moving both need `members.manage` on the group. UserGroupDeleteBlocked is the dialog the editor and list show instead of a delete confirm while a group still has members or child groups; `deleteBlockers.js` counts both the way the server does.
- `Template/` — template editor tooling: CTemplatePreview (render preview with variables/options JSON), CTemplateToolbox (snippet/partial helper). The code editor itself is `CCodeEditor` in lib/vue.
- `ApiGateway/` — CFilterParamsEditor (per-filter parameter form for gateway route filters: workflow picker, HTTP status, response type, …).
- `Workflow/` — WorkflowTriggers (trigger listing panel for the workflow editor).
- `Permissions/` — CPermissionGrid (role × resource-operation permission matrix; the shared grid behind the system, compose, automation and federation permission pages). Its role and user pickers are the shared `CInputRole`/`CInputUser`.

## Data touched

- Panels inject APIs directly (mostly `$SystemAPI`) to resolve names and manage memberships; they do not own Pinia stores.

## When changing this

- These are admin-local; promote a component to lib/vue only when a non-admin section needs it.
- CPermissionGrid changes affect four permission pages at once — verify all of them.
- The commit pattern is per-panel and mixed: UserGroup panels persist member/role changes immediately via `$SystemAPI`, while others (e.g. UserRoles via `initialMembershipIDs`) stage changes for the parent editor's save. Check which pattern a panel uses before reusing or refactoring it.
