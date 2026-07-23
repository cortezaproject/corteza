---
kind: file
covers: useAgentChatStore.ts
backfilled: true
owner: fe
depends-on: []
touched-by:
  - lib/vue/src/components/agent/CAgentSidebar.vue
tests: []
---

# useAgentChatStore

## Intention

Session-local chat state for the agent sidebar: per-agent conversation tabs,
plus lazy-loaded server-persisted history — chats survive panel close and
section navigation.

## State owned

`availableAgents` (sorted, with auto-selected `activeAgentID`), per-agent
`conversations` tab lists + active tab index, per-agent `history` /
`historyLoading` / `historyLoadedAt`.

## API surface consumed

`$SystemAPI.aiConversationList/Delete` (message send is done by the sidebar
component, not the store).

## Consumers

`CAgentSidebar` and its tab/history subcomponents.

## Invariants

- Every agent always has at least one conversation tab; closing the last tab
  re-inits an empty one.
- Tab `id`s are monotonically increasing per agent (max+1) — used as keys.
- `openConversationFromHistory` focuses an already-open tab for the same
  conversationID, and reuses a single empty fresh tab instead of stacking.
- `deleteConversation` soft-deletes server-side, then prunes both history and
  any open tabs for that conversationID.
