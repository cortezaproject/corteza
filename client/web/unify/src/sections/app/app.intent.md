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
  - client/web/unify/e2e/sections/app/gotchas.spec.ts
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

Transport: the app posts `{type: 'human:hello', v: <n>}` to its parent; the
host answers with `{type: 'human:port'}` carrying one `MessagePort`. Every call
after that goes over the port as `{id, op, args}` and is answered with
`{id, result}` or `{id, error}`; `error` is a string naming what was refused.

Contract version: `n` is the contract the page was written against. The shell
prefixes its own copy of the bridge, so the handshake that arrives always
carries the shell's number; the version is therefore read from the source (the
`human:hello` its own copy sends). A handshake with no number is 1; a page with
no copy of its own gets the current contract, 2. A page keeps the contract it
was written for — changing what an existing version returns breaks deployed
apps, so a change to a value's shape is a new version, never an edit to an old
one.

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
  multi-value fields as arrays, an empty value absent.
  - Contract 1: every value is the string the store holds.
  - Contract 2: a Bool is `true`/`false` and always present (the store keeps
    false as nothing, which an app cannot tell from unset); a Number is a
    number. Everything else is as in 1.
- `refs` maps user and record IDs to labels, in every contract (it only gains
  keys). A record's label follows the MCP tools and the webapp's viewers: the
  reference's `labelField`, else the target module's first field, one further
  level when that field is itself a reference; at most 500 per call, one list
  call per target module, as the viewer.
- A reference's target module need not be declared to be labelled. The
  reference is part of a declared record and Human's own viewers show its label
  there; the app gets that label and nothing else from the target.
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

## Gotchas — where each kind is fixed

A custom app that breaks in Human is sorted by what broke, and each kind has
one home. A gotcha with no test in that home is not written down anywhere else.

| Kind                                                                     | Fixed in                                                                         | Pinned by                                                                       |
| ------------------------------------------------------------------------ | -------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| API shape — what the store returns (a false Bool as nothing, `"12.5"`)   | the bridge (`host.js`); a change to an existing shape is a new contract version  | a probe in `e2e/sections/app/gotchas.spec.ts`                                   |
| Sandbox — what the frame blocks (network, storage, dialogs, new windows) | the deploy guard (`checkApplicationSource`), refusing with what to write instead | a case in `application_handler_test.go`                                         |
| Data — what a namespace holds (empty modules, invented Select values)    | a rule in the `custom_app` skill, read before the page is written                | nothing automated yet — a fresh session run against the skill is the only check |

The skill carries only what neither the bridge nor the guard can absorb; when
one of them takes a gotcha over, its line leaves the skill.
