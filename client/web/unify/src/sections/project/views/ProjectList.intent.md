---
kind: file
covers: ProjectList.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/components/project/NewProjectDialog.vue
touched-by:
  - client/web/unify/src/sections/project/index.js
tests:
  - client/web/unify/e2e/sections/project/project-list.spec.ts
  - client/web/unify/e2e/sections/project/lifecycle-dashboard.spec.ts
---

# ProjectList view

## Intention

The project section's entry point: every project in one backend-driven list,
from which a project is created, opened (wizard or dashboard), renamed,
archived or deleted.

## UX capabilities

- Server-side searched, sorted and paginated project list (name + description,
  lifecycle status tag, updated date); mutations re-list from the backend.
- Filter by lifecycle status from the toolbar. It is a user-facing filter and
  rides the reactive filter state; the chain-head constraint below is not —
  that one is structural and must never be clearable from the UI.
- One row per project, not per revision: the list requests chain HEADS
  (`headsOnly`) — the row no other revision names as parent. Server-side by
  necessity, since the list is paginated and searched on the backend;
  grouping after fetch would split chains across pages.
- Create a project — on success the user lands in the wizard with `?new=1`,
  the just-created flag the Wizard uses to auto-open the members dialog once.
- Row click favours the dashboard once a chain has ever published (ruled
  2026-07-28): it opens `project.overview` even when a draft is in progress,
  because the dashboard is the project's home and the wizard is entered
  deliberately from its revision switcher. Only a chain that has never
  published — draft, no parent revision — opens the wizard directly.
- Per-row actions: rename (edits `meta.short`), archive/unarchive
  (status `archived` ↔ `draft`), delete with confirmation.

## Routes

- `project.list` at `/project/projects`; the section index `/project`
  redirects here. Links out to `project.wizard` (with `?new=1` on create)
  and `project.overview`.

## When changing this

- Live means status `active` — the locked contract says the backend never
  sets `published`; the extra `published` check here is defensive only.

- Search matches handle OR the meta JSON name (`meta.short`) via the rdbms
  filter override (`server/store/adapters/rdbms/filter.go` `f.Project`) —
  ruled and verified 2026-07-24; the generated filter alone is handle-only.

- The name column is not sortable (meta JSON column has no sort ident).
- Keep the `?new=1` handoff contract in sync with Wizard.vue.
