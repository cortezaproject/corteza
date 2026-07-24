---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/config
  - lib/js
touched-by:
  - client/web/unify/src/sections/project/views
  - client/web/unify/src/sections/project/components
  - client/web/unify/src/sections/project/sidebar/ProjectSidebar.vue
tests: []
---

# Project stores

## Intention

The section's Pinia layer — the locked stores/composables architecture puts
all cross-view project state here. Views and components never call project
APIs directly for this state; they read stores and call store actions.

## Map

Each store file carries its own sidecar doc:

- `projects.js` — the canonical project store: CRUD + lifecycle, per-project
  resource caches, graph view state, effective access, and the session-local
  governance scaffolding (see its WIP note).
- `events.js` — dashboard events across the five category resources.
- `backlogItems.js` — dashboard backlog sub-issues linked to category events.
- `report.js` — thin wrapper over the server-side aggregation endpoint.
- `users.js` — user directory for team/owner pickers.
- `dateUtils.js` — not a store: shared date-normalization helper for the
  dashboard stores' payloads.

## Cross-cutting concerns

- Lifecycle statuses handled here are exactly the locked set: `draft` →
  `active` (publish flips; the BE never sets `published`), plus `archived`,
  `suspended` and soft-`deleted`.
- Dashboard stores (events, backlogItems) share idioms: flat reactive list
  for the active project, wholesale replace on load, in-place patch after
  mutations (no refetch), user refs resolved to names on read / IDs on
  write, dates normalized via `dateUtils.js`.
- IDs are always compared as strings (`String(id)`).

> **WIP:** approval persistence — the per-step governance/approval state
> (status, note, form values) is session-local scaffolding inside
> `projects.js`; nothing of it reaches the API and a reload resets it. Do
> not rely on it surviving a reload and do not "fix" the lack of
> persistence.

## When changing this

- New per-project state follows the keyed-by-projectID cache pattern in
  `projects.js`, never fields patched onto the Project instance.
