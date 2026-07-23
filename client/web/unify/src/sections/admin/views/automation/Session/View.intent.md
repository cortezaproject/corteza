---
kind: file
covers: View.vue
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin/routes.js
tests: []
---

# Session View view

## Intention

Inspect a single workflow session — what ran, when, triggered by whom, and why
it failed — and cancel it while it is still active.

## UX capabilities

- Read-only detail: session/workflow IDs, status tag, created/completed timestamps, event and resource type.
- Initiating user resolved to a display name via `$SystemAPI.userRead`, best-effort — falls back to the raw ID on failure.
- Error output rendered as a dedicated panel when present.
- Actions: open the owning workflow's editor; cancel — offered only while `completedAt` is unset, then the session is reloaded.

## Routes

`automation.sessions.view` at `/automation/sessions/:sessionID`; back navigates to `automation.sessions`, "open workflow" to `automation.workflows.edit`. Fetch failure redirects back to the list.

## When changing this

- A session is cancelable only while `completedAt` is unset — keep `isActive` derived from that, not from status strings.
- User resolution must stay non-fatal; the screen renders fully without it.
