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
- Create a project — on success the user lands in the wizard with `?new=1`,
  the just-created flag the Wizard uses to auto-open the members dialog once.
- Row click routes by lifecycle: a live project opens its dashboard
  (`project.overview`), anything still draft opens the wizard.
- Per-row actions: rename (edits `meta.short`), archive/unarchive
  (status `archived` ↔ `draft`), delete with confirmation.

## Routes

- `project.list` at `/project/projects`; the section index `/project`
  redirects here. Links out to `project.wizard` (with `?new=1` on create)
  and `project.overview`.

## When changing this

- Live means status `active` — the locked contract says the backend never
  sets `published`; the extra `published` check here is defensive only.

> **DRIFT:** the search box is dead in practice (found 2026-07-24 via e2e):
> the backend `ProjectFilter` matches `query` against `handle` only
> (`server/store/adapters/rdbms/filters.gen.go`), and the create flow never
> sets a handle — so searching by the visible name finds nothing. Which side
> to fix (BE matches name/meta vs FE derives a handle on create) awaits a
> ruling; the e2e specs deliberately avoid the search box until then.

- The name column is not sortable (meta JSON column has no sort ident).
- Keep the `?new=1` handoff contract in sync with Wizard.vue.
