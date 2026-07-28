---
kind: file
covers: BoardView.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/project/components/dashboard/BoardPanel.vue
touched-by:
  - client/web/unify/src/sections/project/index.js
tests:
  - client/web/unify/e2e/sections/project/lifecycle-dashboard.spec.ts
---

# BoardView view

## Intention

Route target for the dashboard's **Board** — the project's kanban across
every revision in the chain. A thin wrapper: the screen itself is
`components/dashboard/BoardPanel.vue`, the same component the wizard's
Manage & Monitor Board section mounts with a revision. Chain-wide here, so
each card carries the revision it belongs to; the wizard drops that chip
because everything there shares one revision.

## UX capabilities

- Work items of all six types as cards in the four shared status columns
  (`config/eventForm.js` EVENT_STATUS); drag between columns writes status
  optimistically and rolls back on failure.
- Dragging changes status ONLY — an item's revision is never altered by a
  move, which matters here precisely because this board spans revisions.
- Columns always render, each with its own quick-add, so an empty project
  is still workable rather than a dead end.

## Routes

- `project.overview.board` at `/project/projects/:projectId/board`, a child
  of the dashboard layout; route meta carries `titleKey`/`icon` for the rail.

## When changing this

- Keep it thin. Board behaviour belongs in `BoardPanel`, so the dashboard and
  the wizard cannot drift — that is the whole point of the shared panel.
- The layout loads the events/backlog stores chain-wide for its children;
  this view must not load them itself.
