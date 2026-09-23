---
name: custom_app
description: A custom app is one HTML document Human renders in a sandbox with no network and no storage — how to write one that runs both in a preview and inside Human, and how to deploy it.
importance: high
announce: Asked to build a page, dashboard, report view or small app for someone — even just a preview or an HTML artifact that uses Human data — the skill custom_app decides what that page may contain, and a page written without it cannot be deployed into Human.
triggers:
  - system_application_create
  - system_application_update
  - system_application_source_set
  - system_application_source_get
---

# Custom apps

A custom app is one HTML file shown inside a sandboxed frame. The frame has no
network, no storage and no login. Everything it knows about Human arrives
through the `human` bridge. The same file must also run in your own preview,
where there is no bridge — so every app carries sample data and says which mode
it is in.

## Before writing

1. Read the modules the app uses (`compose_module_lookup`) and ten or so real
   records (`compose_record_lookup`). Use the real field names. Real data is
   messy — empty values, duplicates, missing labels — and the app must render it
   without blank rows or "undefined".
2. Sample values must come from the module's real Select options, not from
   invented ones. A status the module cannot hold makes the preview a
   demonstration of something that will never appear.
3. Read the colours (`system_theme_lookup`) and use them as CSS custom
   properties. At run time `human.theme()` sends the live ones; yours are the
   fallback.

## The file

- One plain HTML document. No React, no JSX, no `import`, no `export`, no build
  step. Scripts run inline.
- Libraries only from `https://cdnjs.cloudflare.com`, pinned to an exact
  version. No web fonts, no external stylesheets, no external images — inline
  them or use `data:` URIs.
- No `fetch`, `XMLHttpRequest` or `WebSocket`. The sandbox is served with
  `connect-src 'none'`, so they fail with nothing to render. Data comes from
  `human.*` only.
- No `localStorage`, `sessionStorage`, cookies or IndexedDB — not even inside a
  `try` — the frame's origin is opaque and they throw. Keep state in variables.
- No `alert`, `confirm`, `prompt` (the sandbox suppresses them; `confirm` always
  answers false), no `window.open` or `target="_blank"`, no `<a download>`.
  Use in-page dialogs and in-page detail panes, and `human.download(name, text)`
  for an export — up to 5 MB of text, which Human saves for the person.
- No links that leave the page — `mailto:`, `tel:`, another site or another
  file. They work in a preview and do nothing in Human. Link only to `#anchors`
  within the page, and show an email address or phone number as text.

`system_application_source_set` refuses a document that breaks any rule above
and says which; writing to them from the start saves the round trip.

Asked for something in that list, build the page anyway with the nearest thing
that works inside Human and say in one line what you did instead — an address
shown as text rather than a mail link, a filter that resets when the page is
closed. Do not stop to ask; a page that does most of it is worth more than a
question.

- Keep it under 60 KB. Deploying sends the whole file as one tool argument, and
  the hard limit is 256 KB.

## The bridge

Human prefixes this same block to every app it renders, and the snippet is
written to be idempotent, so include it and the second copy is a no-op:

    window.human = window.human || (() => {
      let port = null, seq = 0; const waiting = new Map()
      const ready = new Promise(resolve => {
        const t = setTimeout(() => resolve(false), 600)
        addEventListener('message', e => {
          if (e.source !== parent || e.data?.type !== 'human:port' || !e.ports[0]) return
          clearTimeout(t); port = e.ports[0]
          port.onmessage = m => { const w = waiting.get(m.data.id); if (!w) return
            waiting.delete(m.data.id); m.data.error ? w.reject(new Error(m.data.error)) : w.resolve(m.data.result) }
          resolve(true)
        })
        try { parent.postMessage({ type: 'human:hello', v: 2 }, '*') } catch { resolve(false) }
      })
      const call = (op, args) => ready.then(live => {
        if (!live) {
          const sample = window.SAMPLE || {}
          return sample[op] ? sample[op](args) : Promise.reject(new Error('no sample data for ' + op))
        }
        return new Promise((resolve, reject) => { const id = ++seq
          waiting.set(id, { resolve, reject }); port.postMessage({ id, op, args }) })
      })
      return { ready, call,
        records: { list: a => call('records.list', a), read: a => call('records.read', a),
                   report: a => call('records.report', a) },
        user: () => call('user'), theme: () => call('theme'),
        resize: height => call('resize', { height }) }
    })()

It must stay a global property and never a `const` or `let` declaration: two
top-level `const human` blocks in one document is a syntax error and the page
never runs.

Define the fallback as `window.SAMPLE = { ... }` — an assignment, not
`const SAMPLE`, because `call` reads it off `window` at call time. It is an
object keyed by operation name, each entry a function returning what the live
call would, built from the real records you read. Five to fifteen rows; keep one
or two of the messy ones.

Operations, all read-only:

| call                   | args                                            | result                                       |
| ---------------------- | ----------------------------------------------- | -------------------------------------------- |
| `human.records.list`   | `{module, filter?, sort?, limit?, pageCursor?}` | `{records, refs, nextPageCursor}`            |
| `human.records.read`   | `{module, recordID}`                            | `{record, refs}`                             |
| `human.records.report` | `{module, dimension, metrics?, filter?}`        | `{rows, refs}`                               |
| `human.modules`        | —                                               | the declared modules and their fields today  |
| `human.user`           | —                                               | `{userID, name, email}`                      |
| `human.theme`          | —                                               | `{dark, colors}`                             |
| `human.records.create` | `{module, values}`                              | `{record}`                                   |
| `human.records.update` | `{module, recordID, values}`                    | `{record}`                                   |
| `human.download`       | `(name, text)`                                  | `true` — Human saves the file for the person |
| `human.resize`         | `height`                                        | `true` — the shell resizes the frame         |

- `module` is a handle and carries no namespace: the app reads the namespace it
  declared when it was deployed. A module it did not declare is refused with
  `module "<x>" is not declared for this app`.
- A record arrives as `{recordID, values: {Field: value | [values]}, ownedBy,
createdAt, updatedAt}`. `values` is an object keyed by field name, multi-value
  fields as arrays — not the `[{name, value}]` list the MCP tools return, and
  not their string-only values:
  - a Bool field is `true` or `false` and always present;
  - a Number field is a number;
  - a DateTime is an ISO string (`YYYY-MM-DD` for a date-only field);
  - a Record or User field is the ID, and `refs` maps it to the label a person
    would see — look it up there rather than showing the ID;
  - any other field is a string, and absent when it is empty.
- Write `SAMPLE` in that same shape: `true`/`false` and numbers, not the `"1"`
  and `"12.50"` the MCP tools show you.
- The snippet's `v: 2` is the contract a page is written against. Keep it as it
  is; a page with `v: 1` gets the older, string-only values.
- `limit` is capped at 500.
- A report groups by one field, named `dimension`, and every row comes back as
  `{dimension_0, count}` — the record count is there without asking for it:

      const { rows } = await human.records.report({ module: 'contacts', dimension: 'contact_type' })
      // rows: [{ dimension_0: 'billing', count: 1 }, { dimension_0: null, count: 3 }]

  `dimension_0` is the stored value, so a Select needs its option's wording
  looked up, and `null` means the records that hold nothing there — count them
  as "Not set" rather than dropping them. `metrics` is only for a sum or an
  average, and each one must say so: `'SUM(amount) AS total'`.

- A report grouped by a Record or User field comes back keyed by bare IDs, with
  `refs` naming them — show the name, never the ID.
- `human.modules()` answers `[{handle, moduleID, writable, fields}]`, each field
  `{name, label, kind, multi, options?}` as the module stands now. A Select's
  options are `[{value, text, label}]`, where `text` and `label` are the same
  wording. Take a field's label and a Select's options from there rather than
  writing them into the page: an option renamed next month should not leave the
  page lying. A value the options do not cover still has to render — fall back
  to the value itself rather than printing nothing.
- Changing records is off unless the app was deployed with `writes` naming the
  modules it may change; every one of them must be in `modules` too. Leave
  `writes` out for a read-only app.
  - `human.records.create({module, values})` and
    `human.records.update({module, recordID, values})`. `values` is an object
    keyed by field name, in the same shape reading gives you — `true`/`false`
    and numbers, not strings. An update changes only the fields it lists, and
    always exactly one record.
  - Human asks the person once, while the app is open, before the first change.
    A refusal comes back as an error; say so in the page rather than retrying.
  - Every change runs as the person using the app, under their permissions, so
    a refusal from the server is theirs and its message is worth showing.
  - Give `window.SAMPLE` a `records.create` and `records.update` too, returning
    the record as if it had been saved, so the preview works.
- Deleting a record and changing its owner are not available. Do not write an
  app that needs them.

Two modes, and the app always says which one it is in. `await human.ready` is
`true` when the bridge answered within 600 ms and `false` in a preview. Show a
small badge — "live data" or "sample data" — in the same place either way, and
never hide it: an app quietly showing invented rows is how a decision gets made
on numbers nobody has.

## Deploying

When the user is happy with the preview:

1. `system_application_create` with the title and
   `unify: {"kind":"custom","listed":true}`. Leave `url` out — it is set to
   `app/<applicationID>` for you. Pass `enabled: true` for it to reach anyone.
2. `system_application_source_set` with `source` (the file exactly as
   previewed), `namespace` and `modules` — the allowlist the bridge enforces, so
   a module missing here makes the app refuse its own data — plus `writes` if
   the page changes records.
3. Compare the returned `size` and `hash` with what you sent, then give the user
   the returned `url`.

To change an app later, read it with `system_application_source_get` and send a
patch: `old_string` copied exactly from what came back, with `new_string`.
Somebody may have edited the page in Human since you wrote it, so patch what
`system_application_source_get` returns rather than what you remember sending.
`old_string` must match exactly once, so include the surrounding lines. Sending
the whole file again is right only for a rewrite.
