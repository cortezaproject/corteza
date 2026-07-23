---
kind: file
covers: useWorkflowPromptsStore.ts
backfilled: true
owner: fe
depends-on:
  - lib/js
  - lib/vue/src/components/prompts
touched-by:
  - client/web/unify/src/App.vue
tests: []
---

# useWorkflowPromptsStore

## Intention

Queue of pending workflow user-prompts (workflow paused waiting for input)
plus which prompt dialog is currently shown.

## State owned

`prompts` (`automation.Prompt` list keyed by stateID), `active` — false | true
| the specific prompt being displayed.

## API surface consumed

`$AutomationAPI.sessionListPrompts/ResumeState/Cancel`.

## Consumers

Unify shell (initial `update(webapp)` fetch + realtime-triggered refreshes),
prompt dialog components.

## Invariants

- `update(webapp)` only appends prompts that are FRESH (unknown stateID) and
  whose definition exists in `promptDefinitions` with a matching webapp — the
  UI never receives a prompt it cannot render.
- `resume`/`cancel` remove the prompt locally even when the API call fails
  (finally-block) — a dead prompt must not wedge the queue.
- Removing the active prompt re-activates only if other prompts remain.
