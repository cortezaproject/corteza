---
kind: folder
covers: '.'
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/index.js
tests: []
---

# Project section

## Intention

Project-centric building: a project bundles resources (data model, pages,
automations, agents, chatbots, connections, members) into one governed,
publishable unit, built through a wizard. This doc records what is **locked**
(interview-ruled 2026-07-24); WIP notes mark what must not be relied on.

## Locked contracts

- **Wizard tabs**: Build / Govern / Manage & Monitor — the wizard's shape.
- **Create flow**: projects are created through the create flow, which contains
  the governance (FRIA) questions.
- **Approval UX**: topbar request → approve → publish, with direct per-step
  review. First publish flips status; the live status is `active` — the backend
  never sets a `published` status.
- **Lifecycle**: `draft` → `active` (via publish), plus `archived`,
  `suspended`, and soft-`deleted`.
- **Build tab is canvas-centric**: the graph is the primary surface. Default
  visibility is progressive by design — seeded to the kinds introduced by steps
  up to the active one — while the full graph (every resource and relation) is
  always reachable via the layer chips.
- **Step model**: the step mechanism is locked, and the **Build** step set is
  locked: data model, connections, automations, agents, chatbots, pages, roles,
  permissions, users (source of truth: `config/pipeline.js` STEPS). Govern-tab
  steps (summary, resource-management, data-sensitivity) are governance content
  — WIP, may change.
- **Resource dialog standard**: every resource kind has a CreateDialog and a
  DetailDialog; the Wizard provides `inspectResource(kind, id)` and
  `createResource(kind)`; clicking anything in the graph always opens the
  editable detail dialog.
- **Per-resource permissions**: shared ResourcePermissionsSection embedded in
  each resource dialog; matrix takes a scope prop; automation appears as
  runtime "Run" kind (read + execute); connection deliberately omitted.
- **Members**: users hold per-project roles; the members dialog (with `?new=1`
  auto-open deep link) is the management surface. Groups are part of the
  intended model but **not built yet** — a prior broken experiment was removed
  (2026-07-24); build fresh when the time comes.
- **Dashboards**: two surfaces are intended. The live dashboard on its own
  routes (view set locked: Overview, Category, Backlog, All Events), and the
  wizard's Manage & Monitor tab (currently a placeholder) holding the
  build-time view — current-version categories and backlog items, partially
  shared with the live dashboard since items are assignable to a version.

> **WIP:** approval **persistence** — the flow is session-local FE scaffolding
> (governance backend dropped); backend persistence is planned. Do not rely on
> approval state surviving a reload.

> **WIP:** **Govern tab content and FRIA question content** — anything
> governance-content is due to change; only the flow shapes above are locked.

> **WIP:** **permission create-ops** — parent-scoped create/list/app-access
> operations are missing from the matrix; the scoping split awaits a ruling.

## Map

- `views/` — ProjectList, Wizard, dashboard views (own docs).
- `components/` — wizard machinery, graph, per-resource-kind dialog families,
  permissions, members (roles/users/group), dashboards (own docs).
- `sidebar/`, `composables/`, `stores/`, `utils/`, `config/` — own docs.

## When changing this

- New resource kinds must follow the dialog standard and the permissions
  section contract — no bespoke editors.
- Governance-content changes are expected; flow-shape changes are intent
  changes and need a ruling first.
