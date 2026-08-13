---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/agentic/views/Editor.vue
tests: []
---

# Agentic components

## Intention

The editor's conversational test bench: `AiChat` talks to the agent, `AiTrace`
exposes how each answer was produced. Both are embedded only by
`views/Editor.vue` and operate on a conversation object the editor owns.

## Map

- `AiChat.vue` — chat UI over `CChatMessages`; sends input via `$SystemAPI.agentExec` (threading `conversationID` for follow-ups), appends user/agent messages and a trace entry per exchange, and can embed an `AiTrace` side panel with two-way selection between messages and trace cards.
- `AiTrace.vue` — read-only trace inspector: injected context (rendered markdown), per-exchange cards of decisions (tool call vs response) with payloads and token usage; exposes `selectExternal`/`highlightLatest` and emits `select`.

## Data touched

- Only `$SystemAPI.agentExec` (AiChat); all conversation state lives in the editor.

## When changing this

- AiChat deliberately mutates the `conversation` prop in place (`{label, messages, conversationID/aiConversationID, traceHistory, context}`); AiTrace only reads it. That shape is an informal contract shared with the editor's history and tab handling — change all of them together.
- Errors are surfaced in-band: a failed exec pushes an agent error message plus an `{error}` trace entry; keep that so traces stay aligned with messages by index.
- Chat renders only `user`/`agent`/`assistant` roles with content — other roles are trace-only.
