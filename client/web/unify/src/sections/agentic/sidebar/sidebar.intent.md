---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAgentStore.js
touched-by:
  - client/web/unify/src/sections/agentic/index.js
tests: []
---

# Agentic sidebar

## Intention

The section's left navigation: an "Agents" group holding one child link per
agent, so any agent's editor is one click away. The group header is itself the
link to the agent list — clicking it navigates, its chevron expands — so the
list needs no second entry of its own.

## Map

- `AgenticSidebar.vue` — `CSidebarNav` tree (exact route matching, all groups expanded).

## Data touched

- `useAgentStore` — the sidebar triggers `fetchList()` on mount when the list is empty (the shell does global setup only; sections load their own nav data) and renders reactively from it.

## When changing this

- Per-agent labels fall back `meta.short` → handle → i18n "untitled"; entries are sorted by that label.
- Nav freshness relies on Home/Editor CRUD keeping the store in sync (`updateInList` / `removeFromList`) — don't add a second fetch path here.
- Entries route by name (`agentic`, `agentic.edit`); renames in `index.js` must be swept here. Note the list route hides the sidebar entirely (`meta.hideSidebar`), so this tree is mostly seen from the editor.
