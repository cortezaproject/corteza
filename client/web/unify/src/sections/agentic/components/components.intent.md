---
kind: folder
covers: recursive
owner: fe
depends-on:
  - client/web/unify/src/sections/agentic/toolAccess.js
  - lib/vue/src/composables/useAgentTurn.ts
touched-by:
  - client/web/unify/src/sections/agentic/views/Editor.vue
tests:
  - client/web/unify/src/sections/agentic/components/AgentToolDialog.test.ts
  - client/web/unify/src/sections/agentic/components/AiChat.approval.test.ts
---

# Agentic components

## Intention

The editor's two halves that are too big to live in it: the conversational test
bench (`AiChat` talks to the agent, `AiTrace` exposes how each answer was
produced) and `AgentToolDialog`, where what the agent is allowed to call is
chosen. All three are embedded only by `views/Editor.vue`.

## Map

- `AiChat.vue` — chat UI over `CChatMessages`; runs the turn through `useAgentTurn` (lib/vue), appends user/agent messages and a trace entry per exchange, renders the approval gate, and can embed an `AiTrace` side panel with two-way selection between messages and trace cards.
- `AiTrace.vue` — read-only trace inspector: injected context (rendered markdown), per-exchange cards of decisions (tool call vs response) with payloads and token usage; exposes `selectExternal`/`highlightLatest` and emits `select`.
- `AgentToolDialog.vue` — the tool picker: every tool the server offers, grouped by subject, each row carrying its permission and optional per-tool note and namespace/module narrowing. Takes `tools`/`grants`/`namespaces`/`modules`, emits `apply` with the whole grant list.

## Data touched

- `AiChat` — `$SystemAPI.agentExec`, passed into `useAgentTurn` as its `exec`; all conversation state lives in the editor.
- `AgentToolDialog` — no API of its own. The catalog and the agent's scope arrive as props, and what a grant resolves to is `toolAccess.js`'s to say.

## When changing this

- AiChat deliberately mutates the `conversation` prop in place (`{label, messages, conversationID/aiConversationID, traceHistory, context}`); AiTrace only reads it. That shape is an informal contract shared with the editor's history and tab handling — change all of them together.
- Errors are surfaced in-band: a failed exec pushes an agent error message plus an `{error}` trace entry; keep that so traces stay aligned with messages by index.
- Chat renders only `user`/`agent`/`assistant` roles with content — other roles are trace-only.
- The approval gate sits above the composer, not in the message list: it is a decision to make, not a turn that happened. Approving for the chat remembers the tool per agent and conversation; approving once sends it with that call alone. The memory is a convenience — the server asks again on a conversation it was not told about, and nothing here grants what the invoking user could not already do.
- The dialog edits copies. `draft` (named grants) and `families` are rebuilt from `grants` every time it opens, and only `apply` hands them back — closing any other way must leave the agent exactly as it was.
- **Blocking a tool nothing else grants deletes its entry** rather than storing a `deny`; a `deny` is kept only where a family would otherwise cover the tool. An entry that grants nothing is noise on the wire.
- **Choosing the mode the tool's risk would give anyway clears the override** instead of pinning it, so a tool that should simply follow the rule still follows it if the rule changes. An entry the dialog was opened with is left alone — `named` is what tells an override added this session from one that was already there.
- Family grants are handed back untouched. The dialog chooses tools and no longer writes families, but it must still know which tools an existing family covers.
