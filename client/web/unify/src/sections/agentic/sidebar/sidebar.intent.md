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

The section's left navigation: a Home entry plus an "Agents" group with one
child link per agent, so any agent's editor is one click away.

## Map

- `AgenticSidebar.vue` — `CSidebarNav` tree (exact route matching, all groups expanded, divider before the agents group).

## Data touched

- `useAgentStore` — the sidebar triggers `fetchList()` on mount when the list is empty (the shell does global setup only; sections load their own nav data) and renders reactively from it.

## When changing this

- Per-agent labels fall back `meta.short` → handle → i18n "untitled"; entries are sorted by that label.
- Nav freshness relies on Home/Editor CRUD keeping the store in sync (`updateInList` / `removeFromList`) — don't add a second fetch path here.
- Entries route by name (`agentic`, `agentic.edit`); renames in `index.js` must be swept here. Note the list route hides the sidebar entirely (`meta.hideSidebar`), so this tree is mostly seen from the editor.
