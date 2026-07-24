---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/stores/users.js
  - client/web/unify/src/sections/project/config/roles.js
touched-by:
  - client/web/unify/src/sections/project/views/ProjectList.vue
  - client/web/unify/src/sections/project/views/Wizard.vue
tests: []
---

# Project-level dialogs

## Intention

Dialogs that operate on the project as a whole — create, rename, members —
as opposed to the per-resource-kind dialog families in the sibling folders.
Scope note: this doc governs only `components/project/`; the section-wide
contract lives one level up in `sections/project/project.intent.md`.

## Map

- `NewProjectDialog.vue` — two-step create: (1) name/description, (2) the AI
  Act **Deployer category questions**. Answering yes to any makes a FRIA
  required. **Locked:** these deployer/FRIA questions must stay in the
  create flow (asked for every project — no build-mode gate).
- `MembersDialog.vue` — THE members management surface (replaced the old
  members wizard step): per-member role preset select (config/roles),
  derived capability columns (read / write / request / grant approval), add
  and remove. Opened from the Wizard topbar; auto-opened once for a
  just-created project via the `?new=1` deep link (Wizard.vue handles the
  query flag, ProjectList sets it).
- `RenameProjectDialog.vue` — pure prompt; emits `rename`, caller persists.

## Data touched

`useProjectsStore` — `create` (sends the `deployer` answers), `membersFor`,
`addMember`/`updateMember`/`removeMember` (backend-persisted, unlike
governance state); `useProjectUsersStore` for the user directory;
`ROLE_PRESETS`/`rolePreset` from config/roles define the capability model
the Wizard later reads for review gating.

> **WIP:** FRIA/deployer **question content** is governance content and due
> to change; only their placement in the create flow is locked.

## When changing this

- Member roles feed capability resolution for the whole wizard (write /
  approval gating) — changing presets or columns changes review behaviour.
- Keep create-flow failures non-destructive (dialog stays open with input).
