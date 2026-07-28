---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/pipeline.js
  - client/web/unify/src/sections/project/config/kinds.js
  - client/web/unify/src/sections/project/stores/projects.js
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
  - client/web/unify/src/sections/project/components/dashboard
tests:
  - client/web/unify/e2e/sections/project/wizard-steps-graph.spec.ts
  - client/web/unify/e2e/sections/project/lifecycle-dashboard.spec.ts
---

# Wizard machinery

## Intention

The moving parts of the project wizard: step-nav chrome, review banner, the
generic governance form renderer, and one panel per pipeline step (`steps/`). The Wizard view (`views/Wizard.vue`) owns state and dispatch;
everything here is presentational or a thin list over store caches.

**Locked:** the step mechanism AND the step set. `config/pipeline.js` `STEPS`
is the single source of truth — key, type (`form|resource|sensitivity|
permissions`), resource `kind`, and `tab` (Build/Govern; Manage & Monitor has
no steps). Adding a step = pipeline entry + panel here + dispatch branch in
Wizard.vue; `kindsThroughStep` then scopes the graph automatically.

**Locked approval UX:** direct per-step review — a granter may Approve or
Request changes (required note) on any Build/Govern step, any status, any
time; there is no per-step submit stage — those actions live in the step
header (ruled 2026-07-28). Only the well-known `publish` governance step keeps
a request → approve → publish cycle, surfaced in the wizard header row beside
the tabs. Capabilities gate ACTIONS only, never tab/step visibility;
governance status never locks editing.

## Map

- `StepNav.vue` — step list; badges flip to amber (changes-requested) /
  green (approved) from per-step governance status.
- `ManageNav.vue` — the Manage & Monitor rail. Deliberately NOT a step list:
  prop-driven (`activeKey` in, `select` out) over `config/manageNav.js`, so
  that tab navigates without steps. Visual twin of the dashboard's own rail.
- `manage/` — one component per M&M section, mounted by Wizard.vue's
  key→component map; most just mount a shared panel from
  `components/dashboard/` with the open revision. One file per section is
  deliberate: they are built independently and must not collide. The board's
  columns always render — an empty revision must stay usable, with per-column
  quick-add. Item viewing/editing REUSE the dashboard's dialogs and drawers;
  never fork a board-local editor, or the two surfaces drift.
- No bottom toolbar (ruled 2026-07-28): per-step review moved to the step
  header, Save to a footer inside the step panel (form steps only), prev/next
  dropped — the step nav already lists every step.
- `StepStatusBanner.vue` — surfaces the changes-requested review note.
- `GovernanceForm.vue` — schema-driven form renderer; also reused by the
  dashboard event/backlog dialogs (cross-folder consumer).
- `LlmProviderDialog.vue` (ResourceManagementStep); `YesNo.vue` — boolean
  glyph, also used by project/MembersDialog.
- `fria/` + Fria* step panels — the Art. 27 assessment: a determination step
  plus a scenarios step. Taxonomy keys are persisted data; scenario state is
  session-local scaffolding like other Govern steps, resetting on reload.
- `steps/` — one panel per STEPS entry. Resource steps are lists that open
  the Wizard-mounted dialogs via injected `createResource(kind)` /
  `inspectResource(kind, id)` (module fields: `editField`/`createField`) —
  a step NEVER mounts its own create/detail dialog.

## Data touched

`useProjectsStore`: per-kind caches (`loadX`/`xFor`/`removeX`), `touch()`
after mutations (drives graph refresh), and the governance surface
(`governanceStatus/Note/Values`, `saveStepForm`, `transitionStep`).

> **WIP:** approval **persistence** is session-local scaffolding in
> stores/projects.js (governance backend dropped) — the UX above is locked,
> surviving reload is not.

> **WIP:** Govern-step content — GovernanceForm schemas (summaryForm,
> resourceManagementForm) and any FRIA content are due to change.

> **WIP:** permission create-ops — PermissionsStep embeds
> ProjectPermissionMatrix, which still omits parent-scoped create/list ops.

## When changing this

- Never bypass the injected openers or the pipeline registry; flagging one
  step also sends a pending `publish` approval back for review.
