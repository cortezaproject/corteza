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
tests:
  - client/web/unify/e2e/sections/project/project-list.spec.ts
  - client/web/unify/e2e/sections/project/members.spec.ts
---

# Project-level dialogs

## Intention

Dialogs that operate on the project as a whole — create, rename, members —
as opposed to the per-resource-kind dialog families in the sibling folders.
Scope note: this doc governs only `components/project/`; the section-wide
contract lives one level up in `sections/project/project.intent.md`.

## Map

- `NewProjectDialog.vue` — two-step create: (1) name/description, (2) the AI
  Act **Deployer category questions**, which drive the backend's FriaRequired
  derivation. Their placement here was locked, but that ruling was REVERSED
  on 2026-07-28: creation becomes name + description only and the Govern
  tab's FRIA flow takes over the determination. The questions come out when
  that flow lands, never before — deleting them early leaves nothing deciding
  whether a FRIA is required.
- `MembersDialog.vue` — THE members management surface (replaced the old
  members wizard step): per-member role preset select (config/roles),
  derived capability columns (read / write / request / grant approval), add
  and remove. Mounted by `ProjectTopbarTools.vue`, so it is reachable from
  every project surface rather than the wizard alone; auto-opened once for a
  just-created project via the `?new=1` deep link (ProjectList sets it).
- `RevisionSwitcher.vue` / `ProjectTopbarTools.vue` — the shared topbar
  cluster, mounted by BOTH the wizard and the dashboard so neither redefines
  it. The switcher navigates the revision chain and branches new revisions;
  the tools hold Members and View project. The publish/approval cluster is
  deliberately NOT here — its home in the wizard's tab row is locked.
- `RenameProjectDialog.vue` — pure prompt; emits `rename`, caller persists.

## Data touched

`useProjectsStore` — `create` (sends the `deployer` answers), `membersFor`,
`addMember`/`updateMember`/`removeMember` (backend-persisted, unlike
governance state); `useProjectUsersStore` for the user directory;
`ROLE_PRESETS`/`rolePreset` from config/roles define the capability model
the Wizard later reads for review gating.

> **WIP:** FRIA/deployer questions — both their content AND their placement
> are now unsettled: the 2026-07-28 ruling moves the determination out of
> creation entirely, so treat this dialog's step 2 as scheduled for removal.

## When changing this

- Member roles feed capability resolution for the whole wizard (write /
  approval gating) — changing presets or columns changes review behaviour.
- Keep create-flow failures non-destructive (dialog stays open with input).
