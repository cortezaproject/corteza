---
name: corredor-script
description: Write, deploy and verify a Corredor automation script — pick the trigger from the generated reference, write it into an extension directory, confirm the runner loaded it, trigger it and read what it did. Use when a task needs an automation script (server or client), not a workflow or a TAQ.
---

# /corredor-script

Corredor runs user-written automation scripts. One script is an ES module
exporting `{ label, description, security, triggers|iterator, exec }`, loaded
from an extension directory on disk. Nothing about a script is stored in the
database — the file **is** the deployment.

## The reference answers the naming questions

`corredor/SCRIPTING.md` is generated from the files that define the scripting
surface: every resource with its event types, arguments and constraint names,
the trigger and iterator DSL, the security block, and what the execution
context exposes. Read it instead of guessing, and quote names from it verbatim —
a misspelled resource, event type or constraint is a script that loads fine and
never fires.

Regenerate it after touching the surface (an `events.yaml`, the parser, the
context, the helpers):

```sh
pnpm --filter @planetcrust/human-corredor docs:scripting
```

`pnpm --filter @planetcrust/human-corredor test:unit` fails when the committed
file no longer matches its sources, so the reference cannot go stale quietly.

## 1. Pick the trigger

Decide **where it runs** first. A server script runs inside Corredor with an API
token and can change data on the way through. A client script runs in the
browser, sees the page's in-memory objects, and can talk to the user.

| what it should do | binding |
| --- | --- |
| correct or refuse a change | `before('create'\|'update'\|…)` on the resource |
| react once a change is committed | `after(…)` — the record's `after*` events are immutable |
| a button on a page or in an editor | `on('manual')` plus `uiProp('app', …)`, `uiProp('page', …)`, `uiProp('slot', …)` |
| run on a clock | `every('* * * * *')` or `at('…')`, plus `security.runAs` |
| walk an existing set of records | the `iterator` form, with `action: 'update'` to save each one |
| answer an HTTP call | `on('request').for('system:sink').where('request.path', …)` |
| fix a form before it submits | client script on `ui:compose:record-page`, `before('formSubmit')` |

Then take the resource, its event types and its constraint names from the
reference, and narrow the trigger with `where(…)` so it stays inert on every
other namespace and module.

## 2. Write the script

Copy the nearest working example from `dev/fixtures/corredor` — one script per
shape lives there (before-create, manual button, interval with `runAs`,
iterator, sink, and client scripts for the `compose` and `admin` bundles) — and
change it. Do not add scripts to that directory: it is committed fixture data.

The layout inside an extension directory is fixed:

```
<extension>/server-scripts/**/*.js          run inside Corredor
<extension>/client-scripts/<bundle>/**/*.js  bundled and run in the browser
```

The bundles are `compose`, `admin` and `unify`; a script's name — what the API
and the logs call it — is its path from the extension root plus `:default`, e.g.
`/server-scripts/agent-sandbox/ContactActivate.js:default`.

## 3. Deploy it

Corredor loads every extension directory it finds under
`CORREDOR_EXT_SEARCH_PATHS` (colon-separated, globs allowed). In a deployment,
that variable is the whole deployment step: drop the file on a search path.

In this repo Corredor is part of the stack and serves three paths of the
checkout it belongs to — `dev/fixtures/corredor`, `corredor/usr` and
`corredor/usr/*`. So a new script goes to `corredor/usr/server-scripts/…` or
`corredor/usr/<extension>/server-scripts/…`.

```sh
dev/agent/corredor.sh up        # primary checkout: start it, print what server/.env needs
dev/agent/corredor.sh status    # is it listening, and on what port
dev/agent/worktree.sh up NAME   # a worktree: brings up corredor, server and webapp together
```

`corredor.sh` belongs to the primary checkout only; a worktree's Corredor is
already up whenever its stack is, on the worktree's API port plus 50000. The
server dials Corredor over gRPC and refuses to boot when Corredor is enabled and
unreachable, so Corredor comes up first and goes down last.

**Getting the runner to see the file.** Corredor watches the directories it
found when it started and reloads everything on any change under them. So:

- a script added or edited inside an extension directory Corredor already
  serves is picked up within a second;
- a **new** extension directory, or the first `server-scripts/` under one, was
  not there to be watched — restart Corredor, or force a reload by touching a
  script it already watches (`touch dev/fixtures/corredor/server-scripts/SystemPing.js`);
- deleting a file in an unwatched directory leaves the script in the runner's
  list until the next reload — force one, or it stays triggerable.

## 4. Confirm the runner loaded it

Three places say so, in increasing distance from the file:

1. **The runner log** — `.run/corredor.log` in the checkout. A reload prints
   `reloading server scripts` and then `processed` with `valid` and `total`
   counts; a script that failed to parse prints `script error: …` with its path.
   A script whose count never moves was never seen: check the search path.
2. **The API** — `GET /api/system/automation/` (and `/api/compose/automation/`)
   list what the server holds, each entry with its `name`, `triggers`,
   `security` and `errors`. `dev/agent/api.sh GET '/system/automation/'`.
   The server polls Corredor every few seconds, so allow a moment.
3. **The admin Scripts screen** — `/admin/automation/scripts` shows the same
   list with the Corredor connection state, the parse errors per script, and a
   refresh button.

A script with a non-empty `errors` is loaded but unrunnable — the server refuses
to execute it. Fix the file; the next reload clears it.

## 5. Trigger it and read what it did

| shape | how to fire it |
| --- | --- |
| `on('manual').for('system')` | `POST /system/automation/trigger` `{"script": …, "args": {…}}` |
| `on('manual')` on a record | `POST /compose/namespace/{nsID}/module/{mID}/record/{recID}/trigger` |
| `on('manual')` on a module / namespace | `POST …/module/{mID}/trigger`, `POST …/namespace/{nsID}/trigger` |
| `before`/`after` on a resource | do the thing — create, update or delete the resource |
| `every()` / `at()` / deferred iterator | wait for the tick |
| a client script | open the page in the webapp and use it |

The trigger endpoints answer `OK` and do not hand back what the script
returned, so the script has to leave evidence: `ctx.log.info(…)` lands in the
runner log (and in the gRPC response metadata), and anything it saved is
readable through the API. Read the log by tailing `.run/corredor.log`; each
execution prints `executing script …`, then any emitted log lines, then `done`
with its duration.

A script that throws prints its error and stack there instead — and **loses its
log lines**, because the buffer is only serialised once `exec` resolves. So when
a script is being investigated, put the call under `try`/`catch` and log the
message; the error alone will not tell you what the script had already worked
out.

## Traps

- **Triggers are parsed in isolation.** The parser lifts the default export out
  of the file and evaluates that object literal in a bare VM context, so
  `triggers`, `iterator`, `label`, `description` and `security` cannot reference
  an import or a module-scope constant. One that does is loaded with an error
  and never runs. `exec` is exempt — it may import freely.
- **An iterator's query string takes single quotes.** `query` is a server-side
  expression, not JSON: `query: "status = 'dormant'"`.
- **A manual server script on a record does not save it.** The server loads the
  record, runs the script, and returns the result to the page, which applies it
  to the form in memory. The user still has to save. A script that must persist
  should save through `ctx.Compose.saveRecord(…)` itself — and then it is
  changing a record behind the form the user is looking at.
- **`runAs` needs the server to allow it.** `CORREDOR_RUN_AS_ENABLED` defaults
  to on, but with it off a script declaring `runAs` errors instead of running.
  A deferred script has no triggering user, so without `runAs` it carries no
  token at all and every API call it makes fails.
- **Client scripts are bundled.** Corredor rebuilds the bundle on change, but
  the browser holds the one it loaded at sign-in: reload the page to pick up an
  edit.
- **Narrow every trigger.** A `compose:record` trigger without
  `.where('module', …).where('namespace', …)` fires for every record in the
  instance.

## Cleaning up

A script is a file, so removing it is removing the file — then force a reload
(see step 3) and confirm it is gone from `GET /api/system/automation/`. Records
or other data the script created on the dev server are cleaned up the usual
way, with `dev/agent/cleanup.sh`.
