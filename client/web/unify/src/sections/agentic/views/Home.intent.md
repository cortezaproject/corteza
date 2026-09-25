---
kind: file
covers: Home.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAgentStore.js
touched-by:
  - client/web/unify/src/sections/agentic/index.js
tests: []
---

# Agentic Home view

## Intention

The agent catalog: browse and manage all agents the user can see; entry point
into the agent editor.

## UX capabilities

- Server-driven list (`useResourceList` over `$SystemAPI.agentListCancellable`): search, sort, paginate; name column shows `meta.short` with the description as subtitle, plus handle, active/inactive status tag, last-change date. It sorts by that same COALESCE expression (`changedAtField`), so the order matches the value shown.
- Tri-state deleted filter (Without / Including / Only) in a popover; it defaults to Without, and it is the only way a deleted row, and so its Restore action, reaches the list.
- Row click opens the editor only when the row grants update or delete.
- Create button gated by RBAC `agent.create`; wildcard permissions button (`corteza::system:agent/*`) gated by `grant`.
- Per-row actions: permissions (per-agent resource), edit, duplicate (inactive copy with `_copy` handle and "(Copy)" name, then navigates to it), delete with confirm, or undelete for soft-deleted rows.

## Routes

`agentic` at `/agentic` (`meta.hideSidebar: true` — sidebar stays collapsed on the list); links to `agentic.create` and `agentic.edit`.

## When changing this

- Duplicate copies the agent config block-by-block (`meta`, `behavior`, `execution`, `access`, `invocation`) — new top-level Agent fields must be added there or copies silently lose them.
- Delete/duplicate sync `useAgentStore` so the sidebar tree stays consistent.
