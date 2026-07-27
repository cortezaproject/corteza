---
kind: folder
covers: '.'
owner: fe
depends-on:
  - lib/vue
  - lib/js
touched-by:
  - client/web/unify/src/sections/index.js
tests:
  - client/web/unify/e2e/sections/project/wizard.spec.ts
  - client/web/unify/e2e/sections/project/project-list.spec.ts
  - client/web/unify/e2e/sections/project/wizard-steps-graph.spec.ts
  - client/web/unify/e2e/sections/project/lifecycle-dashboard.spec.ts
  - client/web/unify/e2e/sections/project/members.spec.ts
---

# Project section

## Intention

Project-centric building: a project bundles resources (data model, pages,
automations, agents, chatbots, connections, members) into one governed,
publishable unit, built through a wizard. This doc records what is **locked**
(interview-ruled 2026-07-24); WIP notes mark what must not be relied on.

## Locked contracts

- **Wizard tabs**: Build / Govern / Manage & Monitor — the wizard's shape.
- **Create flow**: creation collects name and description only. Governance
  determination moved into the Govern tab's FRIA flow (ruled 2026-07-28); the
  deployer-category questions are removed when that flow lands.
- **Approval UX**: wizard-header request → approve → publish (right-aligned in
  the wizard's tab row; ruled 2026-07-24, moved out of the app topbar), with
  direct per-step review. First publish flips status; the live status is
  `active` — the backend never sets a `published` status.
- **Lifecycle & revisions**: `draft` → `active` (via publish), plus `archived`,
  `suspended`, soft-`deleted`. A project is a revision chain — each revision is
  its own row (root/parent/number); publish locks what that revision built, and
  a new revision branches from an active one (one draft per chain).
- **Build tab is canvas-centric**: the graph is the primary surface. Default
  visibility is progressive by design — seeded to the kinds introduced by steps
  up to the active one — the full graph is always reachable via layer chips.
- **Step model**: the step mechanism is locked, and the **Build** step set is
  locked: data model, connections, automations, agents, chatbots, pages, roles,
  permissions, users (source of truth: `config/pipeline.js` STEPS). Govern
  steps are governance content — WIP, may change.
- **Resource dialog standard**: every resource kind has a CreateDialog and a
  DetailDialog; the Wizard provides `inspectResource(kind, id)` and
  `createResource(kind)`; clicking anything in the graph always opens the
  editable detail dialog.
- **Per-resource permissions**: shared ResourcePermissionsSection embedded in
  each resource dialog; matrix takes a scope prop; automation appears as
  runtime "Run" kind (read + execute); connection deliberately omitted.
- **Members**: users hold per-project roles; the members dialog (with `?new=1`
  auto-open deep link) is the management surface. Groups are intended but
  **not built** — a prior broken experiment was removed (2026-07-24).
- **Dashboards**: two surfaces. The live dashboard on its own routes (view set
  locked: Overview, Category, Backlog, All Events) covers the whole project
  across revisions; the wizard's Manage & Monitor tab covers one revision —
  board, metrics and activity over the same items. Work items carry a revision
  the way an issue carries a milestone: filed against the project, assigned to
  a revision, reassignable, and left behind when a revision publishes.

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
