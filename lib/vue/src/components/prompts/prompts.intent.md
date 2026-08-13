---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useWorkflowPromptsStore.ts
  - lib/js
touched-by:
  - client/web/unify/src/App.vue
tests: []
---

# prompts/

## Intention

Client-side half of automation workflow prompts: when a server workflow
suspends on a prompt step, these components surface it to the user (toast or
modal), collect input, and resume/cancel the workflow through
`useWorkflowPromptsStore`.

## Map

- `CPrompts.vue` — mount-once host rendering CPromptToast + CPromptModal;
  `hideToasts` prop for surfaces that only want modals.
- `CPromptToast.vue` — toast list (only while the tab has focus): passive
  prompts always, plus the active ones whenever the modal is not showing one,
  so an active prompt is never invisible.
- `CPromptModal.vue` — active prompt dialog; resumes with collected input.
- `kinds/` — one component per prompt ref (alert, choice, input, options,
  notification, composeRecordPicker) extending `base.vue`.
- `definitions.ts` — exported as `promptDefinitions`: the catalog of prompt
  refs + parameter schemas the workflow editor's prompt-step config renders.
- `utils.ts` — `pVal`/`pType` (read `automation.Vars` `@value`/`@type`),
  variant→PrimeVue severity mapping.

## Contracts consumers rely on

- Apps mount `CPrompts` once; everything else is store-driven (websocket
  delivery, resume/cancel) — no per-screen wiring.
- `definitions.ts` must mirror the server's prompt-step definitions: `ref` and
  parameter names are the wire contract with workflow state.

## When changing this

Adding a prompt kind = server support + `definitions.ts` entry + `kinds/`
component (registered in `kinds/index.ts`). Prompt payloads are
`automation.Vars` — always unwrap via `pVal`, never index `@value` inline.
