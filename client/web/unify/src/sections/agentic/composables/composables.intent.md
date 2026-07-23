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

# Agentic composables

## Intention

Section-local non-visual logic. Currently one module: the editor's split-pane
behavior, kept out of the (already large) Editor view.

## Map

- `useEditorSplit.js` — drag-resizable chat column: returns `chatWidth`, `showChat` and `startChatResize` for the editor's resize handle.

## Data touched

- `localStorage` key `agent-editor-chat-width` — the persisted column width, clamped 320–800px (default 480) on both load and drag.

## When changing this

- The chat column auto-hides below a 1024px viewport (`showChat`); the editor renders neither the column nor the handle then — keep the two gated together.
- Dragging attaches document-level mousemove/mouseup listeners and overrides body cursor/user-select; all of it must be torn down on mouseup and on unmount to avoid leaking global state.
