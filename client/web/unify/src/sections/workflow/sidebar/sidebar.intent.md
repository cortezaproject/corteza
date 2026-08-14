---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/workflow/index.js
tests: []
---

# Workflow sidebar

## Intention

Left navigation for the section: a deep link to every workflow, so any
workflow's canvas is one click away from anywhere in the section. The
"Workflows" header is itself the link to the list — clicking it navigates, its
chevron expands — so the list needs no second entry of its own.

## Map

- `WorkflowSidebar.vue` — CSidebarNav tree: a "Workflows" header (`workflow.list`) over an alphabetized (by display label) child list of workflows → `workflow.edit`.

## Data touched

- Shared `workflowStore` from `lib/vue`: fetches the list on first mount only
  if empty — the shell does global setup, each section loads its own nav data.
- No section-local state; renders purely from the store list.

## When changing this

- Labels resolve `meta.name` → `handle` → untitled fallback; keep this in sync with how the list and editor title workflows.
- Freshness relies on Home/Editor calling `workflowStore.updateInList`/`removeFromList` on every lifecycle change — the sidebar never refetches on its own.
