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
tests: []
---

# Wizard view

## Intention

The project build surface. Locked shapes: Build / Govern / Manage & Monitor tabs, the step pipeline, the resource-graph canvas, the wizard-header approval flow.

## UX capabilities

- Steps (set + tab membership locked, owned by `config/pipeline.js`) in a left
  nav beside the ever-present resource graph (resizable split) — the primary
  canvas, seeded per step to the kinds built so far (layer chips peek past).
- Wizard-header tool cluster, right-aligned in the tab row — wrapping to its
  own right-aligned row when the tabs leave no room, content-driven, no fixed
  breakpoint (only the project
  title + version Tag remain Teleported to the app topbar): approval controls
  on the well-known `publish` governance step: request
  → approve → publish (confirmed, then dashboard handoff). No status Tag in
  the row (ruled 2026-07-24) — the single state-driven button carries the
  status, and its tooltip is the only surface for the reviewer's note on a
  changes-requested publish; members without a capability for the current
  state get no publish control at all. Direct per-step
  review — Approve / Request changes (note required) by grant-capable members
  on any step, any time; a flagged step blocks project approval.
- Provides `inspectResource(kind, id)` / `createResource(kind)`; each kind's
  Create + Detail dialog mounts here exactly once; graph and step clicks always
  open the editable detail dialog (`editField`/`createField` for module fields).
- Members dialog opens from the wizard-header cluster; `?new=1` auto-opens it once for a
  just-created project, then is stripped via `router.replace`.
- Capabilities gate review actions and editing only, never tab/step visibility.

## Routes

- `project.wizard` at `/project/projects/:projectId/wizard`, query `step` /
  `tab` / `new`; links to `project.list`, `project.overview`, `namespace.view`.

> **WIP:** approval **persistence** is session-local (`stores/projects.js`);
> **Govern/FRIA content** will change (flow shapes only are locked; Govern
> steps summary/resource-management/data-sensitivity may change); the
> **Manage & Monitor tab body** is a placeholder — intended to hold the
> build-time dashboard (current-version categories + backlog, partially shared
> with the live dashboard via version assignment);
> **permission create-ops** are missing from the matrix pending a ruling.

## When changing this

- Tab shape, the Build step set, and the dialog/provide contract are locked —
  ruling first. Graph default visibility is step-scoped by design (full graph
  via layer chips).
