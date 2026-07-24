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

# Chatbot dialogs

## Intention

The chatbot resource-kind's implementation of the project resource dialog
standard: CreateDialog + DetailDialog, rendered once by the Wizard and opened
via `createResource('chatbot')` / `inspectResource('chatbot', id)`. Chatbot
configuration proper lives in the chatbot section (`chatbot.edit`), reached by
deep link; these dialogs own only the project-facing name + enabled flag.

## Data touched

- `useProjectsStore`: `chatbotsFor`, `addChatbot`, `updateChatbot`.
- Router: resolves `chatbot.edit` (`chatbotID`) for the "open builder" links.

## Map

- `ChatbotCreateDialog.vue` — name-only create with plain Create and
  Create-and-open-builder buttons (tab opened synchronously before the await
  to survive popup blockers; closed on failure). Never auto-opens the detail
  dialog.
- `ChatbotDetailDialog.vue` — staged name/enabled (committed on Save); embeds
  `ResourcePermissionsSection` (`kind="chatbot"`); footer deep-links to the
  chatbot editor (new tab).

## When changing this

- In the permissions matrix a chatbot is a single runtime "Use" (read) toggle;
  starting a session additionally needs a platform-wide grant the project
  cannot set — do not add UI here that pretends otherwise.
- Dialogs open only through the Wizard's provides (dialog standard); keep the
  synchronous `window.open` pattern in the create flow.
