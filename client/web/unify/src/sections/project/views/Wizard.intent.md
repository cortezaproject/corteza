---
kind: file
covers: Wizard.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/pipeline.js
  - client/web/unify/src/sections/project/stores/projects.js
touched-by:
  - client/web/unify/src/sections/project/index.js
tests:
  - client/web/unify/e2e/sections/project/wizard.spec.ts
  - client/web/unify/e2e/sections/project/wizard-steps-graph.spec.ts
  - client/web/unify/e2e/sections/project/lifecycle-dashboard.spec.ts
  - client/web/unify/e2e/sections/project/members.spec.ts
---

# Wizard view

## Intention

The project build surface. Locked shapes: Build / Govern / Manage & Monitor /
Publish tabs, the step pipeline, the resource-graph canvas, the publish flow.

## UX capabilities

- Steps (set + tab membership locked, owned by `config/pipeline.js`) in a left
  nav beside the ever-present resource graph (resizable split) — the primary
  canvas, seeded per step to the kinds built so far (layer chips peek past).
- Publish tab (ruled 2026-07-29) — the whole request → approve → publish
  cycle, which was a cluster in the tab row until then. Deliberately the ONE
  tab with no left rail: publishing is a rare one-way action, so it reads as a
  sequence, not a place. One primary action at the foot, labelled for the next
  thing to do; members without the capability get no control. Direct per-step
  review — Approve / Request changes (note required) by grant-capable members
  on any step, any time; a flagged step blocks project approval.
- TWO screens behind it, chosen on whether a parent revision exists (ruled
  2026-07-29). A REVISION publishes through numbered stages: what changes
  (diff vs the parent) → bring the data across (record mappings) → get it
  approved → go live. A FIRST publish gets its own screen rather than that flow
  emptied out — nothing is replaced, nothing can be lost, so there is no diff
  and no migration. It answers what the project contains and who is about to
  get it (inventory + members, real data), says what going live does, and
  reports rather than judging readiness. Approval is still required.
- Publish blocks on exactly one thing: a destructive change (removed/retyped
  field with records behind it) whose mapping is undecided. Risk is stated as
  counts, never adjectives; unfinished work items warn but never block. A
  reviewer's note shows in the approval stage, not a tooltip (superseding
  2026-07-24). No confirm dialog — the go-live stage IS the confirmation, and
  typing the project handle is required when records will be lost. Success
  shows an in-tab receipt, not a redirect.
- Provides `inspectResource(kind, id)` / `createResource(kind)`; each kind's
  Create + Detail dialog mounts here once; graph and step clicks always open
  the editable detail dialog (`editField`/`createField` for module fields).
- Manage & Monitor tab — the build-time dashboard for the revision being
  worked on: a grouped left rail (monitor + the five categories) rendered from
  its own nav config, never from `pipeline.js` STEPS, so the "Manage & Monitor
  has no steps" contract holds. The content pane is a key→component map over
  `components/wizard/manage/`, one file per section, so sections stay
  independently buildable. Its header carries the revision, that revision's
  work-item completeness, and a new-item action; the same completeness
  reappears in the Publish tab's go-live stage, warning but never blocking.
- Topbar carries the shared `RevisionSwitcher` + `ProjectTopbarTools`
  (Members, View project) — defined once in `components/project/` and mounted
  by the dashboard too. `?new=1` auto-opens Members once for a just-created
  project, then is stripped via `router.replace`. No "View dashboard" button:
  the switcher's Dashboard entry is that navigation.
- Capabilities gate review actions and editing only, never tab/step visibility.

## Routes

- `project.wizard` at `/project/projects/:projectId/wizard`, query `step` /
  `tab` / `new`; links to `project.list`, `project.overview`, `namespace.view`.

> **WIP:** approval **persistence** is session-local (`stores/projects.js`) —
> the M&M completeness bar must derive from work items, never from it;
> **Govern/FRIA content** will change (flow shapes only are locked; Govern
> steps summary/resource-management/data-sensitivity may change, and the FRIA
> flow is due to absorb the governance questions leaving the create flow);
> **permission create-ops** are missing from the matrix pending a ruling.

## When changing this

- Tab shape, the Build step set, and the dialog/provide contract are locked —
  ruling first. Graph default visibility is step-scoped by design (full graph
  via layer chips).
