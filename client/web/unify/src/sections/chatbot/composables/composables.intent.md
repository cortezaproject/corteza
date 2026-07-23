---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/chatbot/views/Editor/Preview.vue
tests: []
---

# Chatbot composables

## Intention

Section-local non-visual logic. Currently one module: the admin-side client the
Editor preview uses to run a real chatbot session without a published widget
key.

## Map

- `usePreviewClient.ts` — `PreviewClient` class mirroring the public widget API surface against `/api/system/chatbot/preview/*`: open/start/close session, submit message/form/consent, advance step, request/complete handoff, operator message, and an SSE event stream.

## Data touched

- POSTs authenticate with the Corteza access token from `$SystemAPI.accessTokenFn()` (Bearer header); the SSE stream passes it as `?jwt=` because `EventSource` cannot set headers.
- Base URL is derived from `$SystemAPI.baseURL` by replacing the trailing `/system` segment.

## When changing this

- Keep endpoint/parity with the public widget client in `client/web/chatbot-widget` (`src/api.ts`) — the preview is only faithful while the surfaces match.
- Status codes are the contract: 204/202 success, 422 returns form field errors as data, anything else throws.
- Error text comes from the backend's plain-text `http.Error` with a `preview: ` prefix, which is stripped for display — coordinate prefix changes with the server.
