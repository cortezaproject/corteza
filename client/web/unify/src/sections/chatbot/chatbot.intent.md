---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/chatbot/sidebar/ChatbotSidebar.vue
  - client/web/unify/src/sections/chatbot/views
touched-by:
  - client/web/unify/src/sections/index.js
tests: []
---

# Chatbot section

## Intention

Administration area for embeddable website chatbots: define a chatbot (identity,
scenario flow, styling, embed key), preview it live against the real widget, and
operate its visitor sessions. This doc governs only the section entry
(`index.js`); `views/`, `sidebar/` and `composables/` carry their own docs.

## Contract (section values declared by index.js)

- `id: 'chatbot'`; every route is tagged `meta.section: 'chatbot'` so the shell resolves the chatbot sidebar.
- Routes are declared inline with full paths, mirroring the legacy chatbot app 1:1 under the `/chatbot` prefix. Legacy generic names are namespaced: `root` → `chatbot`, `sessions` → `chatbot.sessions`.
- Route table: `chatbot` at `/chatbot` (List), `chatbot.create` at `/chatbot/create` and `chatbot.edit` at `/chatbot/:chatbotID/edit` (both Editor), `chatbot.sessions` at `/chatbot/sessions` (Sessions).
- `sidebar: ChatbotSidebar`; no topbar overrides — shell defaults apply.

## Map

- `index.js` — section contract: id, inline route table, sidebar wiring.
- `views/` — route-target views (own sidecars); `views/Editor/` sub-panels carry their own folder doc.
- `sidebar/` — section navigation tree (own doc).
- `composables/` — preview API client (own doc).
- `components/` — empty; no section-local shared components yet.

## Data touched

- `useChatbotStore` (lib/vue) — shared chatbot list feeding sidebar, Sessions view and list CRUD sync.
- `$SystemAPI` `chatbot*` endpoints; the Editor preview additionally talks to `/api/system/chatbot/preview/*` and renders the real widget from `client/web/chatbot-widget`.

## When changing this

- Route names are the API: sidebar and cross-view navigation reference them; route-name uniqueness across sections is convention-only.
- The Editor's scenario panel deep-links into the agentic section by literal path (`/agentic/:agentID/edit`, `/agentic/create`) — coordinate path changes with that section.
