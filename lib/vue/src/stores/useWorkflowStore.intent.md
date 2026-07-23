---
kind: file
covers: useWorkflowStore.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/workflow
tests: []
---

# useWorkflowStore

## Intention

Shared workflow list cache for workflow section views, admin workflow views
and workflow pickers (e.g. the TAQ Run Workflow step selector).

## State owned

`list` — raw workflow objects, `loading`.

## API surface consumed

`$AutomationAPI.workflowList`.

## Consumers

Workflow section home/editor/sidebar, admin automation views, workflow
selector inputs.

## Invariants

- `fetchList(params)` passes filter params straight through and returns the
  full `{ set, filter }` response (callers may paginate); on error it resolves
  to an empty response instead of throwing.
- Editors call `updateInList` (upsert) / `removeFromList` after CRUD instead
  of refetching.
