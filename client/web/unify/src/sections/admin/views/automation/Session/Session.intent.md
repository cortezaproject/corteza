---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Automation — Sessions / executions monitoring

## Intention

Observability for automation runs: one place to see what executed, its status,
and why it failed. Covers both engines — TAQ (Trigger Action Query) executions
and classic workflow sessions — plus a workflow-session detail view.

## Map

- `List.vue` — two-tab list: TAQ executions + workflow sessions (see `List.intent.md`).
- `View.vue` — read-only workflow-session detail with cancel (see `View.intent.md`).

## Data touched

- `$AutomationAPI`: ngAutomationAllExecutions (TAQ tab), sessionList/Read/Cancel (workflow tab + detail).
- `$SystemAPI.userRead` — best-effort createdBy resolution.

## When changing this

- The two engines have distinct status vocabularies (strings vs numeric codes)
  — don't unify without a server-side contract change.
- Screen-level contracts (tab lazy-fetch, cancelability) live in the sidecars.
