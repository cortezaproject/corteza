---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/chatbot/views/Editor
  - lib/vue/src/stores/useChatbotStore.js
touched-by:
  - client/web/unify/src/sections/chatbot/index.js
tests: []
---

# Chatbot Editor view

## Intention

Create or edit a single chatbot with immediate visual feedback: the config form
sits beside a live, fully functional widget preview.

## UX capabilities

- Config tab composes the `Editor/` panels (General, Scenarios, Styling), all mutating one shared `system.Chatbot` instance in place.
- Sessions tab (edit mode only): `CChatbotInbox` scoped to this chatbot, 5s refresh.
- Permanent right-hand live Preview column (hidden below `md`).
- Save gated: a name or handle is required and every scenario needs a unique, non-empty id; create redirects to the edit route of the new chatbot.
- Unsaved-changes guard via deep compare against a post-load/post-save snapshot; delete gated by `canDeleteChatbot`; widget-key regeneration behind a confirm (`chatbotRegenerateWidgetKey`).

## Routes

`chatbot.create` at `/chatbot/create`, `chatbot.edit` at `/chatbot/:chatbotID/edit`; watches the `chatbotID` param so sidebar navigation between chatbots reloads in place.

## When changing this

- Sub-panels mutate the passed chatbot object directly — the dirty check and the Preview's remount watcher both rely on that single shared reference.
- Save strips derived styling fields (`logoURL`, `launcher.iconURL`, empty attachment IDs) before submit; any new derived/display-only field must be stripped there too.
- The agent list loaded here feeds the Scenarios panel's agent picker.
