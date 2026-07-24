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

The shell drawer for the project section: a flat tree of the user's projects
plus an "All projects" entry, so any project is one click away from anywhere
in the section.

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

- Entries: "All projects" (routes to `project.list`), then a divider group of
  projects sorted by name.
- Routing rule (locked lifecycle): a live project (`status: active`; the BE
  never sets `published`, handled only defensively) opens its dashboard
  (`project.overview`); anything else opens the wizard (`project.wizard`).
- Archived projects are hidden from the tree (soft-deleted ones never reach
  the store); the All Projects list still shows everything.
- Each project carries a one-letter status badge; severity mapping mirrors
  ProjectList's and must stay in sync with it.

## When changing this

- Keep the live-status routing rule identical to ProjectList's — both derive
  "live" from `active`.
- The tree relies on the section preload; do not add its own fetch.
