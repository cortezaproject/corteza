---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useChatbotStore.js
touched-by:
  - client/web/unify/src/sections/chatbot/index.js
tests: []
---

# Chatbot sidebar

## Intention

The section's left navigation: fixed entries for the chatbot list and the
cross-chatbot Sessions inbox, plus one child link per chatbot so any bot is one
click away.

## Map

- `ChatbotSidebar.vue` — `CSidebarNav` tree (exact route matching, all groups expanded).

## Data touched

- `useChatbotStore` — the sidebar triggers `fetchList()` on mount when the list is empty (the shell does global setup only; sections load their own nav data) and renders reactively from it.

## When changing this

- Per-chatbot labels fall back name → handle → i18n "unnamed" → ID; entries are sorted by that label.
- Nav freshness relies on List/Editor CRUD keeping the store in sync (`updateInList` / `removeFromList`) — don't add a second fetch path here.
- Entries route by name (`chatbot`, `chatbot.sessions`, `chatbot.edit`); renames in `index.js` must be swept here.
