---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/chatbot/views/Editor/Preview.vue
tests: []
---

# Chatbot widget

## Intention

The embeddable chatbot: a single self-contained `widget.js` (IIFE, no
framework) that third-party sites drop into any page via one `<script>` tag.
Everything renders inside an open shadow root appended to `document.body`, so
host-page CSS can never touch it and vice versa. `public/` (demo page) and
`dist/` (build output) are excluded from the intent system.

## Map

- `src/entry.ts` — boot: locates its own `<script>` tag, reads the widget key from `data-widget-key` (fallback `?k=` in src), fetches config, wires Engine/UI/API together and opens the session immediately so the server drives the first scenario over SSE. Also owns static-message auto-advance/auto-close timers.
- `src/api.ts` — `WidgetAPI`: the public HTTP client contract against `/api/widget/v1` (origin derived from the script's own src). Session-token Bearer auth, `credentials: 'omit'`.
- `src/engine.ts` — thin pub/sub holding messages, current step, and handoff state; only reflects what the server reports via SSE, never orchestrates.
- `src/ui.ts` — DOM layer (no framework): launcher + panel in the shadow root, per-scenario-type rendering, callback hooks (`onUserInput`, `onFormSubmit`, `onConsentDecision`, handoff, close) that entry.ts fills in. Sanitizes author rich-HTML (consent bodies and static messages) against tag/style allow-lists.
- `src/styles.ts` — base CSS string injected into the shadow root; `applyStyling` maps the config's styling block onto `--hb-*` custom properties.
- `src/types.ts` — config, session, and SSE payload shapes; `SSEEventName` enumerates the event vocabulary. It must stay exhaustive: `agent_error` is listened for in `entry.ts` and the admin preview but is currently missing from the union.
- `vite.config.ts` — IIFE build to `dist/widget.js`, es2019 target, CSS inlined (no separate file emitted).

## Data touched

- Widget API only; no Corteza JS client, no cookies. Scenario steps are
  server-driven: the widget submits and the server answers over SSE.
- SSE stream auths via query (`?widgetKey=…&token=…`) because `EventSource`
  cannot set headers; POSTs use the same session token as a Bearer header.

## When changing this

- `src/api.ts` is mirrored by the unify chatbot section's `PreviewClient`
  (`client/web/unify/src/sections/chatbot/composables/usePreviewClient.ts`) —
  endpoint paths, payload shapes, and status-code semantics must stay in
  parity or the editor preview lies.
- Status codes are the contract: `sendMessage` expects 202 (async agent
  reply), `submitConsent`/`advanceStep`/`closeSession`/`closeHandoff` expect
  204, `submitForm` treats 422 as data (field→error map), anything else throws.
- Error text is the backend's plain-text `http.Error` with a `widget: ` prefix
  that the client strips — coordinate prefix changes with the server.
- Must stay dependency-light and framework-free; the only workspace import is
  `renderMarkdown` from `lib/js`. Keep the build a single file — no code
  splitting, no emitted CSS.
- Host pages are hostile territory: keep `:host { all: initial }`, absolute
  z-index, and the sanitizer allow-lists intact.
