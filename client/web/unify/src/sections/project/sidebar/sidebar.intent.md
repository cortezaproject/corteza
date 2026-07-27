---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - lib/vue
touched-by:
  - client/web/unify/src/sections/project/index.js
tests: []
---

# Project section sidebar

## Intention

The shell drawer for the project section: one "Projects" root entry that
routes to the list, with the user's projects as its collapsible children — so
any project is one click away from anywhere in the section.

## Map

- `ProjectSidebar.vue` — the whole sidebar; renders a `CSidebarNav` tree.

## When it shows

Registered as the section's `sidebar` in `sections/project/index.js`, with no
`sidebarDisabledRoutes` — the shell offers it on every `/project/**` route.
The drawer is lazily mounted; the projects list it renders is loaded by the
section `preload` (shell runs it on section entry), so the tree is ready
before the drawer first opens, and store `absorb` keeps it current after
mutations.

## Structure & contracts

- Single root entry "Projects" (ruled 2026-07-24, replacing the earlier
  "All projects" link + divider group): clicking the label routes to
  `project.list` (auto-expanding the children, never collapsing them); the
  right-side chevron toggles the children independently — behavior inherited
  from `CSidebarNavItem`, not implemented here. Children: the projects,
  sorted by name.
- Routing rule (locked lifecycle): a live project (`status: active`; the BE
  never sets `published`, handled only defensively) opens its dashboard
  (`project.overview`); anything else opens the wizard (`project.wizard`).
- Archived projects are hidden from the tree (soft-deleted ones never reach
  the store); the All Projects list still shows everything.
- One entry per revision chain, not per revision. The filter is client-side
  here — unlike ProjectList, which asks the backend for heads — because the
  shared store cache stops being heads-only as soon as a wizard visit absorbs
  a whole chain. Detect heads by which rows are named as a `parentRevisionID`.
- Each project carries a one-letter status badge; severity mapping mirrors
  ProjectList's and must stay in sync with it.

## When changing this

- Keep the live-status routing rule identical to ProjectList's — both derive
  "live" from `active`.
- The tree relies on the section preload; do not add its own fetch.
