---
kind: file
covers: useChatbotStore.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/chatbot
tests: []
---

# useChatbotStore

## Intention

Shared chatbot list cache for the chatbot section's list/editor views.

## State owned

`list` — raw chatbot objects (name-sorted), `loading`.

## API surface consumed

`$SystemAPI.chatbotList`.

## Consumers

Chatbot section list/editor/sidebar views.

## Invariants

- `fetchList` loads ALL chatbots (`limit: 0`) and swallows errors into an
  empty list.
- Editors call `updateInList` (upsert) / `removeFromList` after their own CRUD
  instead of refetching.
