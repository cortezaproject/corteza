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

An application of kind `custom` is one HTML document, written outside Human —
by hand, or by an assistant over the MCP surface — and stored with the
application. This
section shows it at `/app/:applicationID` inside a sandbox that gives it no
network, no storage and no session, and hands it data through a message bridge
that runs every call as the viewer, under the viewer's permissions.

The same sandbox and bridge show a compose page's Custom block — an
application's page by `applicationID`, or the block's own `options.source` (see
`PageBlocks.intent.md`) — through one component, `components/CustomAppFrame.vue`.

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
- A deleted application opens for nobody: the tile goes from the menu and the
  route refuses it, so a link kept by someone stops working too.
- `enabled: false` lands on the disabled screen like every application, except
  for whoever may change the app's page (`canManageSourceOnApplication`): they
  get the app, with a banner saying it is a switched-off preview. The rule is
  `preview.js`; `App.vue` applies it.
- The app reads as its viewer, so a viewer needs the same rights to its
  declared namespace, modules and records that a page over them would need, and
  nothing beyond them: what the page declares is resolved to IDs when it is
  stored (`sourceMeta.namespaceID`, `sourceMeta.moduleIDs`) and read by ID
  afterwards. A page stored before that carries handles only, and the view
  falls back to looking them up, which does need `namespaces.search`. The same
  fallback carries an application brought in from another instance, whose
  stored IDs name nothing here: what it declared by name still does.

## What a page declares

A page names what it may reach, and the bridge holds it to that — an
application in `sourceMeta`, a Custom block in its options, both set with
`components/CustomAppDeclaration.vue`:

- `namespace` — an application picks it; a Custom block reads its page's.
- `modules`, `writes`, `deletes` — handles it reads, and those of them it
  may change and delete from; any outside `modules` is refused when stored. Each is stored beside its
  ID (`moduleIDs`) so a viewer reads by ID.
- `origins` — exact `https://host[:port]`, beyond cdnjs, it may load scripts,
  stylesheets, fonts and images from; no path, no wildcard, at most 20. Ruling:
  whoever may change the page sets them — for a Custom block, anyone who may
  update the page — accepted knowing a listed origin sees every URL requested
  from it. `connect-src` stays `'none'` whatever is listed.
- `automations`, `chatbots` — handles of what it may run and open.
- A Custom block's `params` — the author's settings, read as
  `human.context().params`, so one page serves several blocks.

## Sandbox

The page is loaded with `applicationSourceRead` (Bearer auth; an iframe `src`
could not carry it) and shown as:

    outer <iframe srcdoc> — same origin, carries
      <meta http-equiv="Content-Security-Policy" content="frame-src 'none'">
      and the host script below
        inner <iframe sandbox="allow-scripts allow-forms" srcdoc> — the app, prefixed with
          <meta CSP: default-src 'none';
           script-src 'unsafe-inline' cdnjs <origins>;
           style-src 'unsafe-inline' <origins>; font-src data: <origins>;
           img-src data: <origins>; connect-src 'none'; form-action 'none'>
          and the in-page bridge script

`<origins>` are the page's declared origins (`cspInner` in `host.js`), and a
value that is not a bare https origin is dropped before it is written into the
attribute.

Why two frames: the inner frame's opaque origin keeps it off the viewer's
storage and session; the outer frame's `frame-src 'none'` is what stops the
inner document navigating itself elsewhere (a CSP on the inner document cannot
block its own navigation). The outer frame's host script owns the port and the
API calls, and posts results to the shell only through `postMessage` with the
shell's own origin.

The inner frame is `sandbox="allow-scripts allow-forms"`. `allow-forms` is not
a way out: the inner CSP sets `form-action 'none'` and the script below cancels
every submit in the capture phase, so a form can reach nowhere. What it buys is
the submit event itself — without it the browser blocks submission before any
listener runs, and a page that saves the ordinary way, through a form's own
`submit` handler, renders perfectly and writes nothing.

The inner document starts with `<base href="about:srcdoc">` and a host script
that cancels every link and form navigation except a `#fragment`, before the
bridge and before any app script. Without the base a srcdoc document resolves
links against the shell's URL, so `#section` would load the shell again inside
the frame; without the script a stray link would end the app. The page's own
click handlers still run — only the navigation is cancelled.

The host destroys the inner frame if it fires `load` a second time: that is
now only a page navigating itself by script.

## Bridge

The app talks to Human only through the bridge: one `MessagePort`, calls
answered or refused by `host.js`. The transport, the contract version, every
operation and the rule each is held to are `host.intent.md`; an operation is
added there first.

## Map

- `index.js` — section contract, route, per-app guard.
- `views/AppView.vue` — loads an application's source + meta and shows it in
  the frame; the disabled / refused / empty / error states.
- `components/CustomAppFrame.vue` — the two frames and the host end of the
  bridge for one page: resolves what it declares, builds the documents, answers
  calls, asks consent, mounts chatbots. Used by the app view and the Custom
  block.
- `components/CustomAppDeclaration.vue` — the declaration fields, shared by
  the application editor and the block configurator; holds IDs in its pickers
  and gives names back.
- `host.js` — the host script serialised into the outer frame: port handshake,
  operation dispatch, allowlist, value reshaping. Pure functions, unit-tested.
- `bridge.js` — the in-page bridge prefixed to the app source. The same text
  the MCP skill hands an app's author, so a page written to the skill needs
  nothing added; it is idempotent (`window.human = window.human || …`, sample
  data read lazily from `window.SAMPLE`) because the page usually carries its
  own copy too.
- `bridge.test.js` — host allowlist, reshaping, error texts.
- `preview.js` — who sees a switched-off custom app; `preview.test.js`.

## When changing this

- The two-frame shape and the CSP strings are the security boundary. Loosening
  `connect-src`, adding `allow-same-origin`, or dropping the outer frame's
  `frame-src 'none'` each reopen a way for app code to reach the viewer's
  session or ship data out; do it only with a ruling in this doc. Declared
  origins are the one such ruling so far (see What a page declares).
- What must reach Human over the network for the page (a chatbot, a file) is
  done by the shell, never by opening the sandbox to it.

## Gotchas — where each kind is fixed

A custom app that breaks in Human is sorted by what broke, and each kind has
one home. A gotcha with no test in that home is not written down anywhere else.

| Kind                                                                     | Fixed in                                                                         | Pinned by                                                                                                                                                                   |
| ------------------------------------------------------------------------ | -------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| API shape — what the store returns (a false Bool as nothing, `"12.5"`)   | the bridge (`host.js`); a change to an existing shape is a new contract version  | a probe in `e2e/sections/app/gotchas.spec.ts`                                                                                                                               |
| Sandbox — what the frame blocks (network, storage, dialogs, new windows) | the deploy guard (`checkApplicationSource`), refusing with what to write instead | a case in `application_handler_test.go`                                                                                                                                     |
| Data — what a namespace holds (empty modules, invented Select values)    | a rule in the `custom_app` skill, read before the page is written                | the skill shape test (`server/system/agentic/skills`) and `server/tests/mcp/skills_test.go`, plus a hand review of the skill against the tools it names — never a model run |

The skill carries only what neither the bridge nor the guard can absorb; when
one of them takes a gotcha over, its line leaves the skill.
