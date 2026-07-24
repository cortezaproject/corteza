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

The moving parts of the project wizard: step-nav chrome, toolbar, review
banner, the generic governance form renderer, and one panel per pipeline step
(`steps/`). The Wizard view (`views/Wizard.vue`) owns state and dispatch;
everything here is presentational or a thin list over store caches.

**Locked:** the step mechanism AND the step set. `config/pipeline.js` `STEPS`
is the single source of truth — key, type (`form|resource|sensitivity|
permissions`), resource `kind`, and `tab` (Build/Govern; Manage & Monitor has
no steps). Adding a step = pipeline entry + panel here + dispatch branch in
Wizard.vue; `kindsThroughStep` then scopes the graph automatically.

**Locked approval UX:** direct per-step review — a granter may Approve or
Request changes (required note) on any Build/Govern step, any status, any
time; there is no per-step submit stage. Only the well-known `publish`
governance step keeps a request → approve → publish cycle, surfaced in the
wizard header row beside the tabs (Wizard.vue), not in this toolbar. Capabilities gate ACTIONS only,
never tab/step visibility; governance status never locks editing.

## Map

- `StepNav.vue` — step list; badges flip to amber (changes-requested) /
  green (approved) from per-step governance status.
- `WizardToolbar.vue` — back, centered prev/next stepper, Save (form steps
  only), capability-gated Approve / Request changes.
- `StepStatusBanner.vue` — surfaces the changes-requested review note.
- `GovernanceForm.vue` — schema-driven form renderer; also reused by the
  dashboard event/backlog dialogs (cross-folder consumer).
- `LlmProviderDialog.vue` — LLM provider create (ResourceManagementStep).
- `YesNo.vue` — boolean glyph (also used by project/MembersDialog).
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
