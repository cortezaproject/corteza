---
kind: folder
covers: recursive
owner: fe
depends-on:
  - lib/vue/src/stores/useApplicationsStore.js
  - lib/vue/src/components/navigation/CAppList.vue
  - client/web/unify/src/utils/appReachable.js
  - client/web/unify/src/router/sectionAccess.js
  - server/system/rest/application.go
touched-by:
  - client/web/unify/src/sections/index.js
tests:
  - client/web/unify/src/sections/app/bridge.test.js
  - client/web/unify/e2e/sections/app/app.spec.ts
---

# App section — custom applications

## Intention

An application of kind `custom` is one HTML document, written outside Human
(typically by Claude through the MCP) and stored with the application. This
section shows it at `/app/:applicationID` inside a sandbox that gives it no
network, no storage and no session, and hands it data through a message bridge
that runs every call as the viewer, under the viewer's permissions.

The design goal is that a page previewed as a plain HTML artifact runs unchanged
in Human, and that nothing an app author writes can reach the viewer's token.

## Contract (section values declared by index.js)

- `id: 'app'`; one route `app` at `/app/:applicationID`, tagged
  `meta.section: 'app'`.
- `app: null` at the section level: the section is not gated by one registry
  application. The route guard asks `canAccessApplication` on the application
  named in the route instead (see Access).
- `sidebar: null`; the app owns the whole content area.

## Access

- A custom application's tile has `unify.url = app/<applicationID>`. Being
  local, `appReachable` would ask the section's own `app`; for this section it
  asks the named application's `canAccessApplication` instead, so each custom
  app is offered and admitted on its own `access` grant.
- `/app/:applicationID` for an application the user may not access, or that is
  not `kind: custom`, bounces to home with `denied` like any refused section.
- `enabled: false` lands on the disabled screen like every application.

## Sandbox

The page is loaded with `applicationSourceRead` (Bearer auth; an iframe `src`
could not carry it) and shown as:

    outer <iframe srcdoc> — same origin, carries
      <meta http-equiv="Content-Security-Policy" content="frame-src 'none'">
      and the host script below
        inner <iframe sandbox="allow-scripts" srcdoc> — the app, prefixed with
          <meta CSP: default-src 'none'; script-src 'unsafe-inline' cdnjs;
           style-src 'unsafe-inline'; img-src data:; connect-src 'none';
           form-action 'none'>
          and the in-page bridge script

Why two frames: the inner frame's opaque origin keeps it off the viewer's
storage and session; the outer frame's `frame-src 'none'` is what stops the
inner document navigating itself elsewhere (a CSP on the inner document cannot
block its own navigation). The outer frame's host script owns the port and the
API calls, and posts results to the shell only through `postMessage` with the
shell's own origin.

The host destroys the inner frame if it fires `load` a second time.

## Bridge

Transport: the app posts `{type: 'human:hello', v: 1}` to its parent; the host
answers with `{type: 'human:port'}` carrying one `MessagePort`. Every call after
that goes over the port as `{id, op, args}` and is answered with
`{id, result}` or `{id, error}`; `error` is a string naming what was refused.

Operations in this slice (all read):

| op               | args                                            | result                                          |
| ---------------- | ----------------------------------------------- | ----------------------------------------------- |
| `records.list`   | `{module, filter?, sort?, limit?, pageCursor?}` | `{records, refs, nextPageCursor}`               |
| `records.read`   | `{module, recordID}`                            | `{record, refs}`                                |
| `records.report` | `{module, metrics, dimensions, filter?}`        | what `recordReport` returns                     |
| `user`           | —                                               | `{userID, name, email}`                         |
| `theme`          | —                                               | `{dark: bool, colors: {primary, body-bg, ...}}` |
| `resize`         | `{height}`                                      | `true` — the shell sets the frame height        |

Rules the host enforces, whatever the app asks:

- `module` must be one of `sourceMeta.modules`, in `sourceMeta.namespace`;
  anything else is refused with `module "<x>" is not declared for this app`.
- `limit` is capped at 500.
- A record arrives as `{recordID, values: {Field: value | [values]}, ownedBy,
createdAt, updatedAt}` — `values` is a plain object keyed by field name,
  multi-value fields as arrays, missing fields absent. `refs` maps user and
  record IDs to labels the way the MCP tools do.
- No write, no delete, no attachment in this slice.

## Map

- `index.js` — section contract, route, per-app guard.
- `views/AppView.vue` — loads source + meta, builds the two frames, mounts the
  host, shows the disabled / refused / error states.
- `host.js` — the host script serialised into the outer frame: port handshake,
  operation dispatch, allowlist, value reshaping. Pure functions, unit-tested.
- `bridge.js` — the in-page bridge prefixed to the app source. The same text
  the MCP skill tells Claude to paste, so an app written to the skill needs
  nothing added; it is idempotent (`window.human = window.human || …`, sample
  data read lazily from `window.SAMPLE`) because the page usually carries its
  own copy too.
- `bridge.test.js` — host allowlist, reshaping, error texts.

## When changing this

- The two-frame shape and the CSP strings are the security boundary. Loosening
  `connect-src`, adding `allow-same-origin`, or dropping the outer frame's
  `frame-src 'none'` each reopen a way for app code to reach the viewer's
  session or ship data out; do it only with a ruling in this doc.
- New operations are added to the table above first. Every write operation
  needs its own permission story before it exists.
- The reshape of `values` is the app-facing contract; the MCP tool's
  `[{name, value}]` shape is not what apps see.
