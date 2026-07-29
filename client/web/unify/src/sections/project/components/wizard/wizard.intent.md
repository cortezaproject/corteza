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
governance form renderer, one panel per pipeline step. `views/Wizard.vue` owns
state and dispatch; everything here is presentational or a thin store list.

**Locked:** the step mechanism AND the step set. `config/pipeline.js` `STEPS`
is the single source of truth — key, type (`form|resource|sensitivity|
permissions`), resource `kind`, and `tab` (Build/Govern; Manage & Monitor has
no steps). Adding a step = pipeline entry + panel here + dispatch branch in
Wizard.vue; `kindsThroughStep` then scopes the graph automatically.

**Locked approval UX:** direct per-step review — a granter may Approve or
Request changes (required note) on any Build/Govern step, any status, any
time; no per-step submit stage — the actions live in the step header (ruled
2026-07-28). Only the well-known `publish` step keeps request → approve →
publish, in the Publish tab (ruled 2026-07-29). Capabilities gate ACTIONS
only, never tab/step visibility; governance status never locks editing.

## Map

- `StepNav.vue` — step list; badges flip to amber (changes-requested) /
  green (approved) from per-step governance status.
- `ManageNav.vue` — the M&M rail; NOT a step list: prop-driven (`activeKey`
  in, `select` out) over `config/manageNav.js`; twin of the dashboard rail.
- `manage/` — one component per M&M section via Wizard.vue's key→component
  map; most mount a shared `components/dashboard/` panel with the open
  revision. One file per section (built independently). Columns always
  render — empty revisions stay usable (per-column quick-add); viewing/
  editing REUSES the dashboard's dialogs/drawers, never a forked editor.
- No bottom toolbar (ruled 2026-07-28): per-step review sits in the step
  header, Save in a step-panel footer (form steps only), no prev/next.
- `RevisionCompletenessBar.vue` — work-item progress at the M&M rail's foot
  (ruled 2026-07-28): a stacked status bar (fixed EVENT_STATUS order,
  shared palette) + a lifecycle-WEIGHTED percent — items count their
  status's even 0→1 ramp position (ruled 2026-07-29). Totals via
  `composables/revisionCompleteness.js` (shared with the Publish tab's
  go-live stage); refresh = stores' `mutations` counters, never capped lists.
- `publish/` — Publish tab stage components, mounted like `manage/`; the plan,
  mapping decisions and publish actions are store state, never stage-local.
- `StepStatusBanner.vue` — surfaces the changes-requested review note.
- `GovernanceForm.vue` — schema-driven form renderer; also reused by the
  dashboard event/backlog dialogs (cross-folder consumer).
- `LlmProviderDialog.vue` (ResourceManagementStep); `YesNo.vue` — boolean
  glyph, also used by project/MembersDialog.
- `fria/` + Fria* step panels — the Art. 27 assessment (determination +
  scenarios); taxonomy keys persist, scenario state is session-local.
- `steps/` — one panel per STEPS entry. Resource steps are lists that open
  the Wizard-mounted dialogs via injected `createResource(kind)` /
  `inspectResource(kind, id)` (module fields: `editField`/`createField`) —
  a step NEVER mounts its own create/detail dialog.

## Data touched

`useProjectsStore`: per-kind caches (`loadX`/`xFor`/`removeX`), `touch()` after
mutations (drives graph refresh), and the session-local governance surface.

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
