---
kind: file
covers: projects.js
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config/pipeline.js
  - client/web/unify/src/sections/project/config/kinds.js
  - client/web/unify/src/sections/project/config/sensitivity.js
  - client/web/unify/src/sections/project/utils/fields.js
  - lib/js
touched-by:
  - client/web/unify/src/sections/project/views
  - client/web/unify/src/sections/project/components
  - client/web/unify/src/sections/project/sidebar/ProjectSidebar.vue
  - client/web/unify/src/sections/project/stores/events.js
tests: []
---

# projects store

## Intention

The canonical store for everything project-shaped: project CRUD and the
locked lifecycle, members, every per-project resource collection the wizard
builds, permission/effective-access evaluation, the resource graph's view
state, and (temporarily) the governance workflow.

## State owned

- `projects` — list of lib `system.Project` instances; `absorb(raw)` applies
  payloads onto the cached instance so mutations stay reactive.
- Per-project caches keyed by projectID, never on the Project object:
  members, resources (compose modules incl. field mapping), connections (+
  shared connection library), automations (TAQs), agents, chatbots, access
  roles, project end-users (derived from role membership), compose pages (+
  layouts).
- DAL sensitivity levels (standard scheme seeded on demand).
- Graph view state: `graphVersion` (bumped by `touch()` after every
  persisting mutation — the graph and effective-access caches invalidate on
  it) and the session-only layer/access-overlay visibility sets.
- Effective-access caches per project (role- and user-axis), versioned to
  `graphVersion`, with optimistic cell patching on writes.
- `governanceByProject` — per-step status/note/form values (see WIP below).

## API surface consumed

`$SystemAPI` (project*, projectPublish, projectGetDeploymentPlan, projectGraph, members, connection*,
role*, user*, dalSensitivityLevel*, permissionsTrace/Update), `$ComposeAPI`
(module*, page*, pageLayout*, permissions), `$AutomationAPI` (automations,
permissions) — routed per resource string by `apiForResource`.

## Invariants

- Lifecycle is the locked set: publish promotes `draft` → `active` (BE never
  sets `published`); `archived`, `suspended`, soft-`deleted` complete it.
- `publishProject` must carry real `mappings`: the BE skips record migration
  entirely on an empty set and then soft-deletes the old namespace, so an empty
  publish silently drops every record. `deploymentPlan` is where they come
  from — always read fresh, never cached, since it describes edits made
  seconds ago.
- Revisions are project rows in a chain (root / parent / number), not a field
  to increment: `createRevision` branches from an active project and the BE
  allows one draft per chain; `listRevisions` loads the chain. Work items are
  filed against the root project and point at the revision they're assigned
  to, so a revision reads as a milestone over one shared item pool.
- Resources are fetched-by-projectID only; mutations resync via load, and
  every persisting mutation calls `touch()`.
- Update flows are optimistic with rollback on push failure.

> **WIP:** the whole per-step governance section is SESSION-LOCAL
> scaffolding — the governance backend was removed pending redesign; the
> submit → approve / request-changes rules (publish gate, auto-clear on
> resubmit, immediate send-back) mirror it 1:1 in memory only. A reload
> resets all approval state. Deliberate; do not persist it ad hoc.

> **WIP:** FRIA content — `create()` still sends the deployer-category answers
> that drive the BE `FriaRequired` derivation, but that is ruled for removal
> (2026-07-28): creation becomes name + description only and the Govern FRIA
> flow owns the determination. Removal lands with that flow, never before —
> deleting it early leaves nothing deciding whether a FRIA is required.

> **WIP:** permission create-ops — the effective-access matrix omits
> parent-scoped create/list/app-access operations; scoping awaits a ruling.
