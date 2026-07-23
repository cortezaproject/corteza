---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Session List view

## Intention

One place to see what automation executed and with what status, across both
engines: TAQ (Trigger Action Query) executions and classic workflow sessions.

## UX capabilities

- Two tabs that lazy-fetch on activation: TAQ executions (unpaginated, filtered by string status and automationID via `ngAutomationAllExecutions`) and workflow sessions (paginated `sessionList`, filtered by numeric status code 0–5, sessionID, workflowID).
- Status tags with per-engine severity mapping; both tabs hide the text search.
- Workflow-session rows open the session detail view; TAQ rows are not clickable (no detail yet).

## Routes

`automation.sessions` at `/automation/sessions`; rows navigate to `automation.sessions.view` with `sessionID`.

## When changing this

- The two tabs have distinct status vocabularies — TAQ statuses are strings (created/running/…), workflow filter statuses numeric codes. Don't unify without a server-side contract change.
- TAQ executions have no detail view or pagination yet — deliberate; extend the TAQ tab, not the workflow path, when adding them.
