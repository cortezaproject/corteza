---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/project/index.js
tests: []
---

# Project views

## What this area is

Route-target views of the project section. Every file here (and in
`dashboard/`) is pointed at by a route in `../index.js` and carries its own
sidecar intent doc — this folder doc only maps the family.

## Map of contents

- `ProjectList.vue` — section entry: list, create, open, archive, delete.
- `Wizard.vue` — the build surface (tabs, steps, graph, approval, dialogs).
- `dashboard/` — the locked Manage & Monitor view family for live projects:
  layout + Overview, CategoryView, AllEventsView, BacklogView, DashboardStub
  (each with a sidecar; no folder doc — there are no non-route files).

## Cross-cutting concerns

- Routes and route meta live in `../index.js`; views set the topbar via
  Teleport (`#topbar-title` / `#topbar-tools`) — one owner per screen
  (DashboardLayout owns it for all dashboard children).
- Lifecycle split is by CHAIN, not by row status: a project whose chain has
  never published opens the Wizard, anything else opens the dashboard — so a
  branched `draft` revision stays on the dashboard (`config/publishState.js`
  `chainHasPublished`).

## When changing this

- New route-target views need their own sidecar (SPEC §3 file tier).
- The dashboard view set is locked; adding/removing one needs a ruling. Set
  re-ruled and PINNED 2026-07-28: Overview, Board, Activity (the renamed All
  Events), Category, Backlog, Reports. Pinned because panels were twice built
  against a set that then moved — treat further changes as expensive. The dashboard and the wizard's Manage & Monitor
  tab are ONE surface differing only in scope — same sections, same shared
  panel components, chain-wide here and revision-filtered there. A section
  built for one must be built as a panel taking an optional revision.
