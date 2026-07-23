---
kind: file
covers: Sessions.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useChatbotStore.js
touched-by:
  - client/web/unify/src/sections/chatbot/index.js
tests: []
---

# Chatbot Sessions view

## Intention

Cross-chatbot operator inbox: monitor visitor sessions and handle handoffs
across every chatbot the user can see, in one place.

## UX capabilities

- Single `CChatbotInbox` (lib/vue) fed with the IDs of all chatbots from `useChatbotStore`.
- Status filter spans the full lifecycle: `handoff_requested`, `handoff_active`, `active`, `closed`; filter UI shown; 5s auto-refresh.
- All strings resolved from the `chatbot.inbox.*` locale namespace via `makeChatbotInboxTranslations`.

## Routes

`chatbot.sessions` at `/chatbot/sessions`; linked from the section sidebar.

## When changing this

- "All chatbots" is expressed by passing the full ID set — the inbox issues one list call and filters client-side, so this stays cheap; a chatbot missing from the store is invisible here.
- The view fetches the store list on mount only when empty; it trusts the list/editor CRUD flows to keep the store current.
