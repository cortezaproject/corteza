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
tests:
  - lib/vue/src/stores/useWorkflowPromptsStore.scope.test.ts
---

# useWorkflowPromptsStore

## Intention

Queue of pending workflow user-prompts (workflow paused waiting for input)
plus which prompt dialog is currently shown.

## State owned

`prompts` (`automation.Prompt` list keyed by stateID), `active` — false | true
| the specific prompt being displayed. Alongside them: `ownedSessions` (the
sessions this screen started), `foreign` (prompts held back for a webapp that
can render them) and `skipped` (the ones stepped over, for the screen that
started them to report).

## API surface consumed

`$AutomationAPI.sessionListPrompts/ResumeState/Cancel`.

## Consumers

Unify shell (initial `update(webapp)` fetch + realtime-triggered refreshes),
prompt dialog components.

## Invariants

- `update(webapp)` only appends prompts that are FRESH (unknown stateID) and
  whose definition exists in `promptDefinitions` with a matching webapp — the
  UI never receives a prompt it cannot render. The one destructive case: an
  EMPTY server response clears the whole queue, so a prompt pushed by
  websocket can be dropped by a refresh that races it.
- A prompt whose definition names other webapps is never put in `prompts`, but
  it is no longer dropped either: it waits in `foreign`. Once its session is
  claimed with `ownSession` it is resumed with empty input and moved to
  `skipped` — the step's effect is not something this screen can perform, but
  the run must not be left holding a dialog nobody here can answer. The prompt
  push beats the exec response naming the session, so the decision cannot be
  taken on arrival.
- Only a claimed session is stepped over. The queue is the whole user's, so an
  unclaimed foreign prompt is left alone for the webapp it belongs to — another
  tab may be about to render it.
- `resume`/`cancel` remove the prompt locally even when the API call fails
  (finally-block) — a dead prompt must not wedge the queue.
- Removing the active prompt falls back to the picker list (`active = true`)
  when other prompts remain, and closes the dialog when none do.
