---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/components/permissions
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
tests: []
---

# Agent dialogs

## Intention

The agent resource-kind's implementation of the project resource dialog
standard: CreateDialog + DetailDialog for a project's AI agents, rendered once
by the Wizard and opened via `createResource('agent')` /
`inspectResource('agent', id)`. Agent behavior is authored in the agentic
section (`agentic.edit`), reached by deep link; these dialogs own the
project-facing metadata.

## Data touched

- `useProjectsStore`: `agentsFor`, `addAgent`, `updateAgent`.
- Router: resolves `agentic.edit` (`agentID`) for the "open builder" links.

## Map

- `AgentCreateDialog.vue` — name + description create with plain Create and
  Create-and-open-builder buttons (tab opened synchronously before the await
  to survive popup blockers; closed on failure). Never auto-opens the detail
  dialog.
- `AgentDetailDialog.vue` — staged name/description/active (committed on
  Save). The Active toggle maps to the agent's `status` field
  (`active`/`inactive`) — the store patch sends `status`, not a boolean.
  Embeds `ResourcePermissionsSection` (`kind="agent"`); footer deep-links to
  the agent editor (new tab).

## When changing this

- Keep the boolean-toggle ↔ `status` string mapping intact when touching the
  save patch; the backend has no `active` boolean.
- In the permissions matrix an agent is a single runtime "Use" (read) toggle;
  these dialogs must not grow build-time permission UI.
- Dialogs open only through the Wizard's provides (dialog standard).
