---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/js
  - lib/vue
touched-by: []
tests: []
---

# LLMProvider Editor view

## Intention

Create or reconfigure one LLM provider: identity, provider kind, prompt URL,
and the API key — the credential AI features depend on.

## UX capabilities

- Create mode requires provider + API key inline; edit mode moves the key to
  a dedicated "update key" dialog (write-only — entered fresh, never shown).
- Provider and prompt URL sync both ways: picking a provider prefills its
  default URL; typing a known default URL switches the provider, anything
  else flips to "other".
- Delete (gated by `canDeleteLlmProvider`); unsaved-changes guard on leave.

## Routes

- `system.llmProviders.create` → `/system/llm-providers/new`;
  `system.llmProviders.edit` → `/system/llm-providers/:llmProviderID`.
- Successful create redirects to `.edit`; back action → `system.llmProviders`.

## When changing this

- The API key must stay out of the regular save payload on edit — it rides
  only on the explicit key-update call and is cleared afterwards.
- The two provider/URL watchers are a pair; changing one without the other
  causes prefill loops or stuck "other" state.
