---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useChatbotStore.js
touched-by:
  - client/web/unify/src/sections/chatbot/index.js
tests: []
---

# Chatbot List view

## Intention

Browse and manage all chatbots the user can see; the entry point into the
chatbot editor.

## UX capabilities

- Server-driven list (`useResourceList` over `$SystemAPI.chatbotListCancellable`): search, sort, paginate; columns name (with handle subtitle), handle, enabled tag, last-change date.
- Row click opens the editor only when the row grants update or delete.
- Create button gated by RBAC `chatbot.create`; wildcard permissions button (`corteza::system:chatbot/*`) gated by `grant`.
- Per-row actions: permissions (per-chatbot resource), edit, duplicate (copies full config as a disabled chatbot with `_copy` handle, then navigates to it), delete with confirm, or undelete for soft-deleted rows.

## Routes

`chatbot` at `/chatbot`; links to `chatbot.create` and `chatbot.edit`.

## When changing this

- Duplicate builds its payload field-by-field — new Chatbot fields must be added there or they are silently dropped from copies.
- Delete/duplicate/undelete also sync `useChatbotStore` so the sidebar tree and Sessions view stay consistent.
