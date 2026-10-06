---
kind: file
covers: host.js
owner: fe
depends-on:
  - client/web/unify/src/sections/app/bridge.js
  - client/web/unify/src/sections/app/components/CustomAppFrame.vue
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/CustomBlock.vue
tests:
  - client/web/unify/src/sections/app/bridge.test.js
  - client/web/unify/e2e/sections/app/gotchas.spec.ts
  - client/web/unify/e2e/sections/compose/custom-block.spec.ts
---

# Custom app bridge — the host end

## Intention

Every call a custom page makes reaches Human here, and is answered or refused
here, whatever the page asks. The sandbox around it, and what a page declares,
are `app.intent.md`.

## Transport

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
An operation added to the table is not a new version: a page that never calls
it is unchanged, and one written to the skill's copy of the bridge carries the
call with it.

## Operations

| op               | args                                                   | result                                                         |
| ---------------- | ------------------------------------------------------ | -------------------------------------------------------------- |
| `records.list`   | `{module, filter?, sort?, limit?, pageCursor?}`        | `{records, refs, nextPageCursor}`                              |
| `records.read`   | `{module, recordID}`                                   | `{record, refs}`                                               |
| `records.report` | `{module, dimension, metrics?, filter?}`               | `{rows, refs}`, or the rows alone in contract 1                |
| `modules`        | —                                                      | the declared modules, with their fields today                  |
| `user`           | —                                                      | `{userID, name, email}`                                        |
| `theme`          | —                                                      | `{dark: bool, colors: {primary, body-bg, ...}}`                |
| `records.create` | `{module, values}`                                     | `{record}`                                                     |
| `records.update` | `{module, recordID, values}`                           | `{record}`                                                     |
| `download`       | `{name, text}`                                         | `true` — the shell saves the file                              |
| `resize`         | `{height}`                                             | `true` — the shell sets the frame height                       |
| `context`        | —                                                      | `{namespaceID, namespace, pageID, moduleID, recordID, params}` |
| `navigate`       | `{page}` or `{page, recordID}` or `{module, recordID}` | `true` — the shell opens it                                    |
| `files.list`     | `{module, recordID, field}`                            | `[{attachmentID, name, mimetype, size}]`                       |
| `files.read`     | `{module, recordID, field, attachmentID?}`             | `{attachmentID, name, mimetype, size, dataURL}`                |
| `automation.run` | `{automation, input?}`                                 | `{ok: true}`                                                   |
| `chatbot.open`   | `{chatbot}`                                            | `true` — the shell shows the chatbot                           |
| `chatbot.close`  | `{chatbot}`                                            | `true`                                                         |
| `refresh`        | —                                                      | `true` — the rest of the page catches up                       |
| `records.delete` | `{module, recordID}`                                   | `true`                                                         |
| `records.open`   | `{module, recordID, edit?}`                            | `true` — Human shows the record                                |
| `files.upload`   | `{module, recordID, field, name, dataURL}`             | the stored file, as `files.list` describes it                  |
| `users.search`   | `{query, limit?}`                                      | `[{userID, name, email}]`                                      |
| `toast`          | `{message, severity?, title?}`                         | `true` — Human shows it                                        |
| `confirm`        | `{message, title?, accept?, reject?}`                  | `true` or `false`                                              |
| `prompt`         | `{message, title?, value?}`                            | the text, or `null` when dismissed                             |
| `title`          | `{text}`                                               | `true` — the block or top bar shows it                         |

Rules the host enforces, whatever the app asks:

- `module` must be one of `sourceMeta.modules`, in `sourceMeta.namespace`;
  anything else is refused with `module "<x>" is not declared for this app`.
- `limit` is capped at 500 here, but the store answers with at most 200 records
  a call and sets `nextPageCursor` when there are more, so the cap is not the
  page size and a page showing everything has to follow the cursor.
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
- A change (`records.create`, `records.update`) needs three things, and the
  order is the permission story: the module is in `sourceMeta.writes` (a subset
  of `modules`, declared when the page was deployed); the viewer agreed, once
  per app while it is open, through a Human dialog the app cannot draw or
  answer; and the viewer's own permissions allow it, since every call runs as
  them. A value names a field the module has, and is written the way the
  webapp's own editors write it (`true` as `'1'`, `false` as empty).
- `records.update` changes exactly the one record it names and only the fields
  it lists; it refuses to run without a `recordID`, because the endpoint under
  it changes every record a filter matches.
- A `recordID` is digits, or the call is refused with
  `"<x>" is not a record ID`: it is written into the request path, where `../`
  would reach another module, the namespace, or any endpoint the viewer may
  call.
- A report groups by one field, and the app names it `dimension` — what the
  API and its tools call it. The plural is taken as well, because a bridge that
  took only the plural dropped what an author wrote and left the server
  answering `field dimensions is empty`; naming neither is refused here.
- A report grouped by a Record or User field is keyed by bare IDs; from
  contract 2 its `refs` name them, by the same rule a listing follows.
- `modules` answers with what the app declared, as the modules stand now:
  each field's label, kind and Select options. A page that bakes those in goes
  stale the day somebody renames an option. An option carries its wording under
  both `text` and `label`: `text` is what the module calls it and what an author
  reading the module through the API is shown, so a bridge that renamed it left
  pages printing nothing where the option should be.
- An app that throws tells the shell, which says so over the page. Silence is
  the worst answer: a page that stops after its headings reads as a working app
  holding no data, and neither the guard nor the API can see the difference.
  The reporter is the first script in the document, ahead of the bridge, or it
  misses the errors that stop a page before anything else runs.
- `theme` answers with the instance's own palette under the names the API and
  its tools use — `primary`, `secondary`, `success`, `warning`, `danger`,
  `black`, `white`, `light`, `extra-light`, `body-bg`, `sidebar-bg`,
  `topbar-bg` — with `content-bg`, `text`, `text-muted` and `border` read off
  the running document besides. A page is told to read that palette through the
  API before it is written, so sending back a different set of names left every
  page keeping the hex values its author copied in.
- `download` hands the viewer a file the sandbox could not save itself: at most
  5 MB of text, under a name reduced to one file name.
- `context` is a plain copy of where the page is shown (a reactive object
  cannot be posted, and a block's options are a reactive draft in the
  builder): a Custom block's namespace, page, record page module and record,
  and `params`; the app view hands nothing.
- `navigate` opens a page of the namespace the app reads, by handle or ID, or a
  module's record page with a record — Human routes only; the page asks its
  own permissions.
- `files.*` reach a file only through a File field of a record of a declared
  module, read as the viewer; `files.read` hands it as a `data:` URL the shell
  fetched by the attachment's signed address, at most 5 MB.
- `automation.run` follows the writes' story: declared; agreed once per open
  page in a Human dialog of its own; run as the viewer through the endpoints a
  page's automation button uses. Input is flat strings, numbers and booleans;
  on a record page the record goes along.
- `chatbot.open` mounts the declared chatbot's own widget in the shell
  (`mountChatbot`), under its settings: enabled, this site's origin allowed,
  and readable by the viewer. It goes when the page does.
- `records.delete` follows the writes' story with a list and a consent of its
  own: the module is in `deletes` (within `modules`); the viewer agreed, once
  per open page, in a Human dialog apart from the one for changes; and the
  viewer's permissions allow it. It soft-deletes the one record named.
- `files.upload` is a change: a module in `writes`, a File field, the change
  consent, as the viewer. The shell stores the file (at most 10 MB, decoded
  from a `data:` URL) and the field then names it — beside what a multi-value
  field held, instead of what a single one did.
- `records.open` shows a record of a declared module in Human's own record
  view: in the namespace's record modal, over the page, when the app is on a
  page of that namespace; on its own record page otherwise.
- `users.search` looks users up as the viewer, two characters at least, at
  most 50.
- `toast`, `confirm`, `prompt` and `title` put plain text, at most 500
  characters, in Human's own chrome. A confirm or prompt is one at a time and
  drawn by Human, since the sandbox suppresses the browser's.
- A succeeded change (`records.create`, `records.update`, `records.delete`,
  `files.upload`, `automation.run`) or `refresh` tells the host; a Custom
  block then sends `refetch-records`.
- No owner change.

## Events

Human also tells the app things: the shell posts `{type: 'human:event'}` to the
outer frame, which passes `{event, payload}` down the port, and the page
listens with `human.on(event, fn)` (it returns the call that stops listening).
One event so far, `refresh`: something else on the page changed data — another
block, an automation button. A Custom block forwards every `refetch-records`
except the one its own change sent.

## When changing this

- New operations are added to the table above first. Every write operation
  needs its own permission story before it exists; the one above is
  declaration, then the viewer's agreement, then the viewer's permissions, and
  a new one is measured against all three.
- The reshape of `values` is the app-facing contract; the MCP tool's
  `[{name, value}]` shape is not what apps see.
- An operation the page can call is in `bridge.js` and in the `custom_app`
  skill's copy of it, so a page written to the skill can call it in a preview.
