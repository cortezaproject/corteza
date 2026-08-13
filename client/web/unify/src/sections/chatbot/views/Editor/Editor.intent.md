---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/chatbot/composables/usePreviewClient.ts
  - client/web/chatbot-widget/src/engine.ts
  - client/web/chatbot-widget/src/ui.ts
touched-by:
  - client/web/unify/src/sections/chatbot/views/Editor.vue
tests: []
---

# Chatbot Editor sub-panels

## Intention

The config panels and live preview composed by `views/Editor.vue`. Data flows by
in-place mutation: each panel receives the shared chatbot object (Scenarios and
Styling take just their slice) and edits it directly, so the parent's dirty check
and preview remount react automatically. Emits carry actions, never data — the
only one is General's `regenerate-key`.

## Map

- `General.vue` — identity (name/handle/enabled), read-only embed snippet + widget key with copy and confirm-guarded regenerate (emits `regenerate-key`), allowed origins list, session TTL.
- `Scenarios.vue` — ordered scenario steps (reorderable, selectable) of types `static_message` / `conversation` / `form` / `consent`, each with type-specific config and per-step before/after TAQ hooks; handoff section with `onRequested`/`onAccepted` hooks; the conversation type's agent picker deep-links into the agentic section (open/create in new tab).
- `Styling.vue` — widget theming: launcher label/position/shape/size, color slots, font sizes, logo and launcher-icon uploads via `chatbotUploadAssetEndpoint` (disabled until the chatbot is saved — uploads need an ID).
- `Preview.vue` — the real widget (`human-webapp-chatbot-widget` Engine + WidgetUI) mounted contained, driven by `PreviewClient` against the admin preview API over SSE; forwards chatbot identity so preview sessions resolve back to their chatbot in the inbox; remount debounced 400ms per config change to avoid tearing down live SSE sessions on every keystroke.

## Data touched

- `$SystemAPI` asset uploads (Styling) and `/api/system/chatbot/preview/*` (Preview); no Pinia stores.
- Automation hooks are stored as full resource refs `corteza::automation:ng-automation/<id>`; the panels translate to/from bare IDs for `CInputTAQ`.

## When changing this

- Preserve the mutate-the-prop contract — switching a panel to emit-based flow silently breaks dirty tracking and live preview.
- Scenario id uniqueness/non-emptiness is the parent's save gate; keep ids stable across reorder (selection re-maps by id).
- The automation resource-ref prefix must match the backend's ng-automation resource type.
- Preview drops empty-string numeric IDs before POSTing — the server's `uint64,string,omitempty` decoding rejects `""`.
