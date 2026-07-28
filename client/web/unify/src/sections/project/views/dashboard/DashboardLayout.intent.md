---
kind: file
covers: DashboardLayout.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/components/dashboard/DashboardNav.vue
touched-by:
  - client/web/unify/src/sections/project/index.js
tests:
  - client/web/unify/e2e/sections/project/lifecycle-dashboard.spec.ts
---

# DashboardLayout view

## Intention

Shell for a project's dashboards — the locked Manage & Monitor content
family. Owns the topbar (project name, revision switcher, open the compose
namespace) and an in-view left rail (DashboardNav) beside the routed child
view.

## UX capabilities

- Left-rail navigation between the locked dashboard view set; the child view
  renders in a bordered panel and owns its own scroll.
- Scope is the whole revision CHAIN, not one revision: the dashboard is the
  project-wide view, so it loads items across every revision and shows which
  revision each belongs to. The wizard's Manage & Monitor tab is the
  single-revision counterpart (ruled 2026-07-28).
- Topbar carries the shared revision switcher and `ProjectTopbarTools` — the
  same components the wizard header mounts, defined once in
  `components/project/` so neither surface redefines them. The switcher's
  entries are this chain-wide Dashboard plus every revision (each opening
  that revision's wizard), and the action to branch a new revision; the tools
  hold Members and View project. Together they replaced the old version tag,
  the jump-back-to-wizard button and this layout's own namespace link.

## Routes

- Parent route at `/project/projects/:projectId` with children `project.overview`
  (index), `.events`, `.category`, `.reports` (stub), `.backlog`. Child route
  meta carries `titleKey`/`icon` consumed by the nav and the stub view.

## When changing this

- The view set is locked — adding/removing a child view needs a ruling.
- The layout must fetch the project itself (the sidebar drawer is lazily
  mounted; relying on its store load leaves the topbar empty).
- Children must not Teleport into the topbar — this layout owns it.
