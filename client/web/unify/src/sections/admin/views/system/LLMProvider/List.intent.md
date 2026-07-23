---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/vue
touched-by: []
tests: []
---

# LLMProvider List view

## Intention

See every configured LLM provider at a glance — kind, handle, and live
status — and reach the editor or remove providers.

## UX capabilities

- Single full fetch, no search, no pagination (`hide-search`).
- Status column renders a severity tag: active / unauthorized / other.
- Provider column maps raw kind to display name (OpenAI/Anthropic/Mistral/Other).
- Delete via row action menu, gated by `canDeleteLlmProvider`.

## Routes

- `system.llmProviders` → `/system/llm-providers`.
- Header button → `system.llmProviders.create`; row click →
  `system.llmProviders.edit` with `llmProviderID`.

## When changing this

- The list is not paginated — if provider counts ever grow, switch to
  `useResourceList` rather than raising limits ad hoc.
- Status severities/icons must stay consistent with the tag shown in the
  editor header; both read the same `status` values.
