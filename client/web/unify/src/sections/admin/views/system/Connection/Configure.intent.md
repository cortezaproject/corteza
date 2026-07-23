---
kind: file
covers: Configure.vue
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/views/system/Connection/ConfiguredConnectionsPanel.vue
  - lib/vue
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Connection Configure view

## Intention

Manage the configured instances of a catalog (or already-active) connection
without exposing its definition for editing — the "use it" counterpart to
the Editor's "author it".

## UX capabilities

- Read-only header card with the connection's name and description, plus a
  permissions button (active connections only, grant-gated).
- The rest of the page is `ConfiguredConnectionsPanel` (list/delete
  configured instances); the footer offers only Back — no save actions
  exist by design.

## Routes

- `system.connections.configure` →
  `/system/connections/:connectionID/configure`; a local draft loaded here
  redirects to `.edit`; fetch failure falls back to `system.connections`.

## When changing this

- Keep the two-way redirect contract with `Editor.vue` (draft-local → edit,
  catalog → here); this view stays mutation-free for the connection itself
  — instance management belongs to the shared panel.
