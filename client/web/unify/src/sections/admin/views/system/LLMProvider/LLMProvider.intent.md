---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - lib/js
touched-by: []
tests: []
---

# LLMProvider

## Intention

Configure LLM providers the platform's AI features can call: provider kind,
connection config, status, and the API key.

## Data touched

- `$SystemAPI.llmProvider*`; class-based resource `system.LlmProvider`.

## Map

- `List.vue` — unpaginated provider list with status (see sidecar).
- `Editor.vue` — provider form + write-only API-key dialog (see sidecar).

## When changing this

- The API key is a secret: it rides along only on the explicit key-update
  call and is cleared from state afterwards — never echo it back into forms.
