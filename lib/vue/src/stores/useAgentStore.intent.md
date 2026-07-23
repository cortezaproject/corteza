---
kind: file
covers: useAgentStore.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/agentic
tests: []
---

# useAgentStore

## Intention

Shared AI agent list cache so agent list/editor views and pickers agree on one
dataset without refetching.

## State owned

`list` — raw agent objects (name-sorted), `loading`.

## API surface consumed

`$SystemAPI.agentList`.

## Consumers

Agentic section list/editor views; agent pickers.

## Invariants

- `fetchList` loads ALL agents (`limit: 0`) and swallows errors into an empty
  list (list views render empty rather than crash).
- Editors call `updateInList`/`removeFromList` after their own CRUD so the
  list stays fresh without a refetch; `updateInList` upserts.
