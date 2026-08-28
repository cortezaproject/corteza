---
kind: folder
covers: '.'
owner: fe
depends-on:
  - client/web/unify/src/sections/agentic/sidebar/AgenticSidebar.vue
  - client/web/unify/src/sections/agentic/views
touched-by:
  - client/web/unify/src/sections/index.js
tests: []
---

# Agentic section

## Intention

Administration area for AI agents: define an agent (model, behavior, tool and
data access, invocation surfaces) and test it in a built-in chat with full
execution traces. This doc governs the section entry (`index.js`); `toolAccess.js`
carries its own sidecar, and `views/`, `components/`, `composables/` and
`sidebar/` carry their own docs.

## Contract (section values declared by index.js)

- `id: 'agentic'`; every route is tagged `meta.section: 'agentic'` so the shell resolves the agentic sidebar.
- `app: 'agentic/'` — the registry application gating entry to the section; the user must hold `access` on it (see `sections.intent.md`).
- Routes are declared inline with full paths, mirroring the legacy agentic app 1:1 under the `/agentic` prefix; names are namespaced `agentic.*` to avoid cross-section collisions.
- Route table: `agentic` at `/agentic` (Home, the agent list), `agentic.create` at `/agentic/create` and `agentic.edit` at `/agentic/:agentID/edit` (both Editor).
- The root list route sets `meta.hideSidebar: true` — the list keeps the sidebar collapsed, like the legacy app.
- `sidebar: AgenticSidebar`; no topbar overrides — shell defaults apply.

## Map

- `index.js` — section contract: id, application gate, inline route table, sidebar wiring.
- `toolAccess.js` — what an agent's tool grants add up to, shared by the Editor and the tool dialog (own sidecar).
- `views/` — Home (list) and Editor, each with its own sidecar.
- `components/` — AiChat / AiTrace test-bench pair used by the Editor (own doc).
- `composables/` — editor split-pane behavior (own doc).
- `sidebar/` — section navigation tree (own doc).

## When changing this

- The chatbot section's scenario editor deep-links here by literal path (`/agentic/:agentID/edit`, `/agentic/create` in new tabs) — path changes must be coordinated there.
- Route names are referenced by the sidebar and by cross-view navigation; uniqueness across sections is convention-only.
