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

Shell for a live project's dashboards — the locked Manage & Monitor content
family. Owns the topbar (project name + 1-based version tag, jump back to the
wizard, open the compose namespace) and an in-view left rail (DashboardNav)
beside the routed child view.

## UX capabilities

- Left-rail navigation between the locked dashboard view set; the child view
  renders in a bordered panel and owns its own scroll.
- Loads the project itself plus the events and backlog stores on entry and on
  every project switch — child views assume these stores are loading/loaded.

## Routes

- Parent route at `/project/projects/:projectId` with children `project.overview`
  (index), `.events`, `.category`, `.reports` (stub), `.backlog`. Child route
  meta carries `titleKey`/`icon` consumed by the nav and the stub view.

## When changing this

- The view set is locked — adding/removing a child view needs a ruling.
- The layout must fetch the project itself (the sidebar drawer is lazily
  mounted; relying on its store load leaves the topbar empty).
- Children must not Teleport into the topbar — this layout owns it.
