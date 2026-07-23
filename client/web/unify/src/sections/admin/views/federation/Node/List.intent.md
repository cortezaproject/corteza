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

# Node List view

## Intention

Browse federation nodes (remote instances this one shares data with), see
where each pairing handshake stands, and initiate/confirm pairing.

## UX capabilities

- Paginated, searchable, sortable node list with pairing-status tags (`paired` / `pair_requested` / unknown); rows open the node editor.
- "Pair" dialog: paste a federation URI generated on another instance — creates the node and starts the handshake (`nodeCreate` + `nodePair`).
- Per-row actions: permissions (gated by grant), and confirm-pending-pair when status is `pair_requested` (`nodeHandshakeConfirm`).
- Wildcard permissions button gated by `federation/` grant; resource `corteza::federation:node/<id|*>`.

## Routes

`federation.nodes` at `/federation/nodes`; navigates to `.create` (new button) and `.edit` with `nodeID` (row click).

## When changing this

- Pairing is a two-sided flow (this dialog consumes the URI the editor's generate-URI dialog produces on the other instance) — keep both consistent with the server's pairing protocol.
- The federation admin group is feature-flagged in the sidebar; the routes exist regardless.
