// The sandbox around a custom app: the two documents it runs in, and the rules
// every call it makes goes through. See app.intent.md.
//
// The outer document is a `srcdoc` string, so it can import nothing — its
// script is generated here as self-contained text. It holds no token and makes
// no API call: it owns the port to the app and relays each call to the shell,
// which answers it by running `dispatch` with the real clients.

// The outer document's own policy. `frame-src 'none'` is what stops the inner
// document navigating itself somewhere else; a policy on the inner document
// cannot block its own navigation.
export const CSP_OUTER = "frame-src 'none'"

// What an allowed origin may be, and nothing else: it is written into a
// policy inside an attribute, so a value that could close either is dropped.
const ORIGIN = /^https:\/\/[a-z0-9.-]+(:[0-9]{1,5})?$/

// The app's policy: no network calls of any kind; script from inline, cdnjs
// and the origins the page lists; style, font and image from inline or data:
// and those origins.
export function cspInner(origins = []) {
  const listed = (origins || []).filter(o => ORIGIN.test(o)).join(' ')
  const from = base => (listed ? `${base} ${listed}` : base)

  return (
    "default-src 'none'; " +
    `script-src ${from("'unsafe-inline' https://cdnjs.cloudflare.com")}; ` +
    `style-src ${from("'unsafe-inline'")}; ` +
    `font-src ${from('data:')}; ` +
    `img-src ${from('data:')}; ` +
    "connect-src 'none'; " +
    "form-action 'none'"
  )
}

export const CSP_INNER = cspInner()

// A link or form in the app that would load another document does nothing
// instead. The page's own handlers still run; only the navigation is cancelled,
// because a navigating frame is one the outer host has to remove.
export const INERT_NAVIGATION = `(function () {
  document.addEventListener('click', function (e) {
    var a = e.target && e.target.closest && e.target.closest('a[href]')
    if (!a) return
    var href = a.getAttribute('href') || ''
    if (href.charAt(0) === '#') return
    e.preventDefault()
  }, true)
  document.addEventListener('submit', function (e) { e.preventDefault() }, true)
})()`

// An app that throws draws nothing, or worse, draws its headings and stops —
// which reads as a working page holding no data. The frame cannot say so
// itself, so it tells the shell, which puts it in front of the viewer.
export const REPORT_FAILURE = `(function () {
  function tell(message) {
    try { parent.postMessage({ type: 'human:failed', message: String(message).slice(0, 300) }, '*') } catch (e) {}
  }
  addEventListener('error', function (e) { tell((e && e.message) || 'the page stopped') })
  addEventListener('unhandledrejection', function (e) {
    var reason = e && e.reason
    tell((reason && reason.message) || reason || 'the page stopped')
  })
})()`

// A peer connection reaches the network through STUN and TURN, which no
// policy directive covers (`connect-src` does not, and browsers ignore
// `webrtc 'block'`), so the app's window has none. A child frame cannot hand
// one back: inside the sandbox every frame is an opaque origin of its own.
export const NO_PEER_CONNECTIONS = `(function () {
  ['RTCPeerConnection', 'webkitRTCPeerConnection', 'mozRTCPeerConnection'].forEach(function (name) {
    try { delete window[name] } catch (e) {}
  })
})()`

// The most records one call may ask for.
export const MAX_LIMIT = 500

// The bridge contract a page gets when it names none. See app.intent.md, Bridge.
export const BRIDGE_VERSION = 2

// The largest file an app may hand the viewer.
export const MAX_DOWNLOAD = 5 * 1024 * 1024

// A file name of the app's choosing, reduced to one: no path of its own, no
// leading dot, and a name when it offers none.
export function downloadName(name) {
  const cleaned = String(name || '')
    .replace(/[\\/]/g, '-')
    .replace(/^\.+/, '')
    .trim()
    .slice(0, 120)
  return cleaned || 'download.txt'
}

// The most referenced records one call resolves into `refs`.
export const MAX_REFS = 500

// The contract a page was written against, read from the handshake its own copy
// of the bridge sends. The shell prefixes its copy to every page, so the
// handshake that arrives always carries the shell's number; the source is the
// only place that says what the author wrote. A page with no copy of its own
// relied on the shell's, which is the current one.
export function bridgeVersion(source) {
  const text = String(source || '')
  const hello = text.match(/human:hello['"]?\s*,\s*v\s*:\s*(\d+)/)
  if (hello) return Number(hello[1])
  return /human:hello/.test(text) ? 1 : BRIDGE_VERSION
}

// A string as a JS literal safe to sit inside a `<script>` body: `</script`
// anywhere in it would end the element it is written into.
function embed(value) {
  return JSON.stringify(value).replace(/<\//g, '<\\/')
}

// The outer document's script, with the shell's origin baked in — every
// message it sends the shell, and every one it accepts, names that origin.
export function hostScriptSource({ origin }) {
  return `(function () {
  var SHELL = ${embed(origin)}
  var seq = 0
  var pending = {}
  var loads = 0
  var appPort = null

  var frame = document.createElement('iframe')
  // allow-forms is what lets a form fire its own submit event. Without it the
  // browser blocks the submission before any handler runs, so a page that
  // saves the ordinary way — a form and a submit listener — does nothing at
  // all and says nothing about it. Submitting anywhere stays impossible twice
  // over: the inner CSP sets form-action 'none', and the script above cancels
  // the navigation in the capture phase.
  frame.setAttribute('sandbox', 'allow-scripts allow-forms')
  frame.srcdoc = window.__humanInner

  // A second load is the app navigating itself: the document that answered the
  // handshake is gone, so the frame goes with it.
  frame.addEventListener('load', function () {
    if (++loads < 2) return
    frame.remove()
    parent.postMessage({ type: 'human:navigated' }, SHELL)
  })

  addEventListener('message', function (e) {
    var data = e.data || {}

    if (e.source === parent) {
      if (e.origin !== SHELL) return
      // Something happened in Human the app may want to know about.
      if (data.type === 'human:event') {
        if (appPort) appPort.postMessage({ event: data.event, payload: data.payload })
        return
      }
      if (data.type !== 'human:result') return
      var waiting = pending[data.id]
      if (!waiting) return
      delete pending[data.id]
      waiting.port.postMessage({ id: waiting.id, result: data.result, error: data.error })
      return
    }

    if (e.source !== frame.contentWindow) return

    if (data.type === 'human:failed') {
      parent.postMessage({ type: 'human:failed', message: data.message }, SHELL)
      return
    }

    if (data.type !== 'human:hello') return

    var channel = new MessageChannel()
    appPort = channel.port1
    channel.port1.onmessage = function (m) {
      var call = m.data || {}
      var id = ++seq
      pending[id] = { port: channel.port1, id: call.id }
      parent.postMessage({ type: 'human:call', id: id, op: call.op, args: call.args }, SHELL)
    }
    // The app's origin is opaque, so '*' is the only target there is.
    e.source.postMessage({ type: 'human:port' }, '*', [channel.port2])
  })

  document.body.appendChild(frame)
})()`
}

// The whole outer document: the inner document as a string, then the host that
// puts it in a sandboxed frame.
//
// The inner document's base is its own address: a srcdoc document otherwise
// resolves links against the shell's URL, so `#section` would load the shell
// again inside the frame instead of scrolling.
export function buildOuterDocument({ source, hostScript, bridgeScript, cspInner = CSP_INNER }) {
  const inner = `<!doctype html>
<meta charset="utf-8">
<base href="about:srcdoc">
<meta http-equiv="Content-Security-Policy" content="${cspInner}">
<script>${REPORT_FAILURE}</script>
<script>${NO_PEER_CONNECTIONS}</script>
<script>${INERT_NAVIGATION}</script>
<script>${bridgeScript}</script>
${source}`

  // An explicit body: the host script appends to it as it runs, and a document
  // written without one has no body element until the parser reaches content.
  return `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="${CSP_OUTER}">
<style>html,body{margin:0;padding:0;height:100%;overflow:hidden}iframe{display:block;width:100%;height:100%;border:0}</style>
</head>
<body>
<script>window.__humanInner = ${embed(inner)}</script>
<script>${hostScript}</script>
</body>
</html>`
}

// The largest file `files.read` hands an app, as a data: URL.
export const MAX_FILE = 5 * 1024 * 1024

// The one record a call names. It is written into the request path, so
// anything but digits could walk that path to another endpoint (`../`).
export function recordIDFor(args, refusal) {
  const id =
    args.recordID === undefined || args.recordID === null ? '' : String(args.recordID).trim()
  if (!id) throw new Error(refusal)
  if (!/^[0-9]+$/.test(id)) throw new Error(`"${id}" is not a record ID`)
  return id
}

// A File field of a declared module, which is the only way an app names a
// file: what it may read is what a record it may read holds.
export function fileFieldFor(ctx, args = {}) {
  const recordID = recordIDFor(args, 'a file is read from a record; name it with recordID')
  const moduleID = moduleIDFor(ctx, args.module)
  const field = fieldsFor(ctx, moduleID).find(f => f.name === args.field)
  if (!field) throw new Error(`module "${args.module}" has no field "${args.field}"`)
  if (field.kind !== 'File')
    throw new Error(`field "${args.field}" is a ${field.kind} field, not a File field`)
  return { moduleID, field: field.name, recordID }
}

async function fileIDs(ctx, moduleID, recordID, field) {
  const record = await ctx.compose.recordRead({ namespaceID: ctx.namespaceID, moduleID, recordID })
  return (record.values || [])
    .filter(v => v.name === field && v.value && v.value !== '0')
    .map(v => String(v.value))
}

export function describeFile(attachment = {}) {
  const original = attachment.meta?.original || {}
  return {
    attachmentID: attachment.attachmentID,
    name: attachment.name || '',
    mimetype: original.mimetype || '',
    size: Number(original.size) || 0,
  }
}

// What an app hands an automation, as the typed envelopes the exec endpoints
// read: a flat object of strings, numbers and booleans, nothing nested.
export function automationInput(input) {
  if (input === undefined || input === null) return {}
  if (typeof input !== 'object' || Array.isArray(input)) {
    throw new Error('an automation input is an object of named values')
  }

  const out = {}
  for (const [name, value] of Object.entries(input)) {
    switch (typeof value) {
      case 'string':
        out[name] = { '@type': 'String', '@value': value }
        break
      case 'number':
        out[name] = { '@type': Number.isInteger(value) ? 'Integer' : 'Float', '@value': value }
        break
      case 'boolean':
        out[name] = { '@type': 'Boolean', '@value': value }
        break
      default:
        throw new Error(`input "${name}" is not a string, number or boolean`)
    }
  }
  return out
}

const TOAST_SEVERITIES = ['info', 'success', 'warn', 'error']

// The longest text an app puts in front of the viewer in Human's own chrome.
export const MAX_TEXT = 500

// The most users one search hands back.
export const MAX_USERS = 50

// The largest file an app uploads, decoded.
export const MAX_UPLOAD = 10 * 1024 * 1024

// Text an app shows in Human's chrome: plain, one value, bounded. Optional
// text may be empty; required text is refused when it is.
function shortText(value, what, optional = false) {
  const text = value === undefined || value === null ? '' : String(value).trim()
  if (!text && !optional) throw new Error(`${what} needs a message`)
  return text.slice(0, MAX_TEXT)
}

// Whether the app may delete records of a module: declared in `deletes`, and
// so one it reads.
export function allowDelete(meta, module) {
  if (!(meta?.deletes || []).includes(module)) {
    return `module "${module}" is not one this app may delete records from`
  }
  return allowModule(meta, module)
}

// A file handed over as a data: URL, decoded and bounded, with a name reduced
// to one file name.
export function dataURLFile(dataURL, name) {
  const m = /^data:([^;,]*)(;base64)?,(.*)$/s.exec(String(dataURL || ''))
  if (!m) throw new Error('a file is uploaded as a data: URL')
  const bytes = m[2] ? atob(m[3]) : decodeURIComponent(m[3])
  if (bytes.length > MAX_UPLOAD) {
    throw new Error(`an upload is at most ${MAX_UPLOAD} bytes; this one is ${bytes.length}`)
  }
  const data = new Uint8Array(bytes.length)
  for (let i = 0; i < bytes.length; i++) data[i] = bytes.charCodeAt(i)
  return { name: downloadName(name || 'file'), type: m[1] || 'application/octet-stream', data }
}

// Where `navigate` may take the viewer, as the app named it: a page of the
// namespace the app reads, by handle or ID, or the record page of one of its
// modules; a record page always with the record to open.
export function navigationTarget(args = {}) {
  const page = args.page === undefined || args.page === null ? '' : String(args.page).trim()
  const module = String(args.module || '').trim()
  const recordID =
    args.recordID === undefined || args.recordID === null ? '' : String(args.recordID).trim()

  if (page && module) throw new Error('navigate takes a page or a module, not both')
  if (!page && !module)
    throw new Error('navigate needs a page (handle or ID), or a module and a recordID')
  if (module && !recordID)
    throw new Error('navigate to a module opens one of its records; name it with recordID')
  if (recordID && !/^[0-9]+$/.test(recordID)) throw new Error(`"${recordID}" is not a record ID`)

  return { page, module, recordID }
}

// The app-facing shape of a record: `values` keyed by field name, a repeated
// field as an array, a field with nothing in it absent.
//
// From contract 2 a value also takes its field's type: a Bool is `true` or
// `false` and always present — the store keeps false as nothing at all, which
// an app cannot tell from unset — and a Number is a number.
export function reshapeRecord(record, fields = [], version = 1) {
  const values = {}
  const kinds = version >= 2 ? kindsOf(fields) : {}

  for (const { name, value } of record?.values || []) {
    if (value === null || value === undefined || value === '') continue
    const typed = typedValue(kinds[name], value)
    if (typed === undefined) continue
    if (name in values) {
      values[name] = Array.isArray(values[name]) ? [...values[name], typed] : [values[name], typed]
    } else {
      values[name] = typed
    }
  }

  for (const [name, kind] of Object.entries(kinds)) {
    if (kind === 'Bool' && !(name in values)) values[name] = false
  }

  return {
    recordID: record?.recordID,
    values,
    ownedBy: record?.ownedBy,
    createdAt: record?.createdAt,
    updatedAt: record?.updatedAt,
  }
}

function kindsOf(fields) {
  const kinds = {}
  for (const field of fields || []) kinds[field.name] = field.kind
  return kinds
}

function typedValue(kind, value) {
  switch (kind) {
    case 'Bool':
      return value === '1' || value === 'true' || value === true
    case 'Number': {
      const n = Number(value)
      return Number.isFinite(n) ? n : undefined
    }
    default:
      return value
  }
}

// The field whose value stands in for a referenced record: the one the
// reference names, else the target's first field — what the webapp's viewers
// and the MCP tools show in its place.
export function labelFieldOf(module, named) {
  const fields = module?.fields || []
  return (named && fields.find(f => f.name === named)) || fields[0] || null
}

// Every record a set of records points at, grouped by the module it lives in.
export function referenceTargets(fields, records) {
  const targets = {}

  for (const field of fields || []) {
    if (field.kind !== 'Record') continue
    const moduleID = field.options?.moduleID
    if (!moduleID || moduleID === '0') continue

    const target = (targets[moduleID] ||= {
      ids: new Set(),
      labelField: field.options?.labelField || '',
      recordLabelField: field.options?.recordLabelField || '',
    })

    for (const record of records || []) {
      for (const { name, value } of record?.values || []) {
        if (name === field.name && value && value !== '0') target.ids.add(value)
      }
    }
  }

  return targets
}

// Labels for the records `targets` names, two levels deep when a label is
// itself a reference. Runs as the viewer: a record they may not read stays
// unlabelled, which reads the same as it did before.
export async function recordLabels(compose, namespaceID, targets) {
  const out = {}
  let budget = MAX_REFS

  const load = async (moduleID, ids) => {
    const wanted = [...ids].slice(0, budget)
    if (!wanted.length) return { module: null, records: [] }
    budget -= wanted.length

    const module = await compose.moduleRead({ namespaceID, moduleID }).catch(() => null)
    if (!module) return { module: null, records: [] }

    const { set = [] } = await compose
      .recordList({ namespaceID, moduleID, recordID: wanted, limit: wanted.length })
      .catch(() => ({}))
    return { module, records: set }
  }

  const valueOf = (record, name) =>
    (record.values || []).find(v => v.name === name && v.value !== '' && v.value != null)?.value

  for (const [moduleID, target] of Object.entries(targets || {})) {
    if (!target.ids.size || budget <= 0) continue

    const { module, records } = await load(moduleID, target.ids)
    const field = labelFieldOf(module, target.labelField)
    if (!field) continue

    let nested = {}
    if (field.kind === 'Record' && field.options?.moduleID) {
      const ids = new Set(records.map(r => valueOf(r, field.name)).filter(Boolean))
      const inner = await load(field.options.moduleID, ids)
      const innerField = labelFieldOf(inner.module, target.recordLabelField)
      if (innerField && innerField.kind !== 'Record') {
        for (const r of inner.records) {
          const label = valueOf(r, innerField.name)
          if (label) nested[r.recordID] = label
        }
      }
    }

    for (const record of records) {
      const value = valueOf(record, field.name)
      if (value) out[record.recordID] = nested[value] || value
    }
  }

  return out
}

// Whether the app may change records in this module: the refusal text, or null.
// Reading is not enough — changing is declared on its own.
export function allowWrite(meta, module) {
  const refusal = allowModule(meta, module)
  if (refusal) return refusal
  if (!(meta?.writes || []).includes(module)) {
    return `module "${module}" is not declared as one this app may change`
  }
  return null
}

// The app's values as the store takes them: one entry per value, a repeated
// field once per item, `true`/`false` as the webapp's own editors write them.
export function toStoreValues(values, fields) {
  const known = new Set((fields || []).map(f => f.name))
  const out = []

  for (const [name, value] of Object.entries(values || {})) {
    if (!known.has(name)) {
      throw new Error(`"${name}" is not a field of this module`)
    }
    for (const one of Array.isArray(value) ? value : [value]) {
      out.push({ name, value: storeValue(one) })
    }
  }

  return out
}

function storeValue(value) {
  if (value === true) return '1'
  if (value === false || value === null || value === undefined) return ''
  return String(value)
}

// Whether the app may touch this module: the refusal text, or null.
export function allowModule(meta, module) {
  if (!module) return 'no module was named'
  if (!(meta?.modules || []).includes(module)) {
    return `module "${module}" is not declared for this app`
  }
  return null
}

export function capLimit(limit) {
  const asked = Number(limit)
  if (!Number.isFinite(asked) || asked <= 0) return undefined
  return Math.min(Math.floor(asked), MAX_LIMIT)
}

function moduleIDFor(ctx, module) {
  const refusal = allowModule(ctx.meta, module)
  if (refusal) throw new Error(refusal)

  const moduleID = ctx.moduleIDs?.[module]
  if (!moduleID) {
    throw new Error(`module "${module}" was not found in namespace "${ctx.meta?.namespace}"`)
  }
  return moduleID
}

// A breakdown grouped by a reference comes back keyed by the target's bare ID;
// these are the labels a person would see in their place.
export async function dimensionRefs(ctx, moduleID, dimension, rows) {
  const field = (ctx.fields?.(moduleID) || []).find(f => f.name === dimension)
  if (!field) return {}

  const ids = new Set(
    (Array.isArray(rows) ? rows : []).map(row => row?.dimension_0).filter(id => id && id !== '0'),
  )
  if (!ids.size) return {}

  if (field.kind === 'User') return ctx.userLabels ? ctx.userLabels([...ids]) : {}

  if (field.kind === 'Record' && field.options?.moduleID) {
    return recordLabels(ctx.compose, ctx.namespaceID, {
      [field.options.moduleID]: {
        ids,
        labelField: field.options.labelField || '',
        recordLabelField: field.options.recordLabelField || '',
      },
    })
  }

  return {}
}

// A field as an app needs to read it: what to call it, what it holds, and the
// options it may hold.
export function describeField(field) {
  const out = {
    name: field.name,
    label: field.label || field.name,
    kind: field.kind,
    multi: !!(field.isMulti ?? field.multi),
  }

  // An option carries its wording twice. `text` is what the module calls it and
  // what an author reading the module through the API sees; `label` is the word
  // every other field here uses. They are always the same string.
  const options = field.options?.options
  if (Array.isArray(options)) {
    out.options = options.map(o => {
      const text = o.text || o.value
      return { value: o.value, text, label: text }
    })
  }

  return out
}

function writableModuleIDFor(ctx, module) {
  const refusal = allowWrite(ctx.meta, module)
  if (refusal) throw new Error(refusal)
  return moduleIDFor(ctx, module)
}

// The viewer is asked, in Human, before this app's first change; the app can
// neither draw that question nor answer it.
async function agreed(ctx, module) {
  if (!ctx.consent) return
  if (!(await ctx.consent(module))) {
    throw new Error('the person using this app did not agree to it changing records')
  }
}

function fieldsFor(ctx, moduleID) {
  return ctx.fields?.(moduleID) || []
}

// User and record labels together. A reference is part of the declared record
// it sits on, and Human's own viewers show its label there, so a target module
// the app did not declare is still labelled — only the label, and only what the
// viewer may read.
async function refsFor(ctx, fields, records) {
  const users = ctx.refs ? await ctx.refs(records) : {}
  const labels = await recordLabels(ctx.compose, ctx.namespaceID, referenceTargets(fields, records))
  return { ...users, ...labels }
}

// Every operation an app can reach. Runs in the shell, as the viewer, under
// the viewer's permissions.
// Operations after which what the rest of the page shows may be stale.
const CHANGES = [
  'records.create',
  'records.update',
  'records.delete',
  'files.upload',
  'automation.run',
  'refresh',
]

// Runs one operation and, once one that changes data has succeeded, tells the
// host, which brings the rest of the page up to date.
export async function dispatch(op, args = {}, ctx) {
  const result = await runOperation(op, args, ctx)
  if (CHANGES.includes(op)) ctx.changed?.()
  return result
}

async function runOperation(op, args = {}, ctx) {
  switch (op) {
    case 'records.list': {
      const moduleID = moduleIDFor(ctx, args.module)
      const { set = [], filter = {} } = await ctx.compose.recordList({
        namespaceID: ctx.namespaceID,
        moduleID,
        query: args.filter,
        sort: args.sort,
        limit: capLimit(args.limit),
        pageCursor: args.pageCursor,
      })
      const fields = fieldsFor(ctx, moduleID)
      return {
        records: set.map(r => reshapeRecord(r, fields, ctx.version)),
        refs: await refsFor(ctx, fields, set),
        nextPageCursor: filter.nextPage || null,
      }
    }

    case 'records.read': {
      const moduleID = moduleIDFor(ctx, args.module)
      const record = await ctx.compose.recordRead({
        namespaceID: ctx.namespaceID,
        moduleID,
        recordID: recordIDFor(args, 'records.read needs the recordID of the record to read'),
      })
      const fields = fieldsFor(ctx, moduleID)
      return {
        record: reshapeRecord(record, fields, ctx.version),
        refs: await refsFor(ctx, fields, [record]),
      }
    }

    case 'records.report': {
      const moduleID = moduleIDFor(ctx, args.module)
      // One field is grouped by, and the API calls it `dimension` — which is
      // also what an author reading the module through the tools is shown. The
      // plural is taken too, because the operation reads like it takes several.
      const dimension = args.dimension ?? args.dimensions
      if (!dimension) throw new Error('a report needs a `dimension`: the field to group by')

      const rows = await ctx.compose.recordReport({
        namespaceID: ctx.namespaceID,
        moduleID,
        metrics: args.metrics,
        dimensions: dimension,
        filter: args.filter,
      })

      // Contract 1 was handed the rows as they come; from 2 a breakdown by a
      // reference carries the labels too, the way a listing does.
      if (ctx.version < 2) return rows

      return { rows, refs: await dimensionRefs(ctx, moduleID, dimension, rows) }
    }

    case 'records.create': {
      const moduleID = writableModuleIDFor(ctx, args.module)
      await agreed(ctx, args.module)
      const record = await ctx.compose.recordCreate({
        namespaceID: ctx.namespaceID,
        moduleID,
        values: toStoreValues(args.values, fieldsFor(ctx, moduleID)),
      })
      return { record: reshapeRecord(record, fieldsFor(ctx, moduleID), ctx.version) }
    }

    case 'records.update': {
      const moduleID = writableModuleIDFor(ctx, args.module)
      // The endpoint behind this changes every record a filter matches, so one
      // record is named here and nothing else can widen it.
      const recordID = recordIDFor(
        args,
        'records.update needs the recordID of the one record to change',
      )
      await agreed(ctx, args.module)

      const values = toStoreValues(args.values, fieldsFor(ctx, moduleID))
      if (!values.length) throw new Error('records.update needs at least one value to change')

      await ctx.compose.recordPatch({
        namespaceID: ctx.namespaceID,
        moduleID,
        recordID: [recordID],
        values,
      })

      const record = await ctx.compose.recordRead({
        namespaceID: ctx.namespaceID,
        moduleID,
        recordID,
      })
      return { record: reshapeRecord(record, fieldsFor(ctx, moduleID), ctx.version) }
    }

    case 'download': {
      const text = String(args.text ?? '')
      if (text.length > MAX_DOWNLOAD) {
        throw new Error(
          `a download is at most ${MAX_DOWNLOAD} characters; this one is ${text.length}`,
        )
      }
      ctx.download(downloadName(args.name), text)
      return true
    }

    case 'files.list': {
      const { moduleID, field, recordID } = fileFieldFor(ctx, args)
      const ids = await fileIDs(ctx, moduleID, recordID, field)
      return Promise.all(ids.map(async id => describeFile(await ctx.attachment(id))))
    }

    case 'files.read': {
      const { moduleID, field, recordID } = fileFieldFor(ctx, args)
      const ids = await fileIDs(ctx, moduleID, recordID, field)
      const id = args.attachmentID ? String(args.attachmentID) : ids[0]
      if (!id) throw new Error(`field "${field}" of record ${recordID} holds no file`)
      if (!ids.includes(id)) {
        throw new Error(`file ${id} is not one field "${field}" of record ${recordID} holds`)
      }

      const file = describeFile(await ctx.attachment(id))
      if (file.size > MAX_FILE) {
        throw new Error(`a file read is at most ${MAX_FILE} bytes; ${file.name} is ${file.size}`)
      }
      return { ...file, dataURL: await ctx.fileData(id) }
    }

    case 'automation.run': {
      const name = String(args.automation || '').trim()
      if (!name)
        throw new Error('automation.run needs the handle of an automation the app declared')
      if (!(ctx.meta?.automations || []).includes(name)) {
        throw new Error(`automation "${name}" is not declared for this app`)
      }
      const input = automationInput(args.input)
      if (!(await ctx.consentToRun())) {
        throw new Error('the viewer did not allow this app to run automations')
      }
      return ctx.runAutomation(name, input)
    }

    case 'refresh':
      return true

    case 'toast': {
      const severity = TOAST_SEVERITIES.includes(args.severity) ? args.severity : 'info'
      ctx.toast(severity, shortText(args.message, 'a toast'), shortText(args.title, '', true))
      return true
    }

    case 'confirm':
      return !!(await ctx.ask('confirm', {
        message: shortText(args.message, 'a confirmation'),
        title: shortText(args.title, '', true),
        accept: shortText(args.accept, '', true),
        reject: shortText(args.reject, '', true),
      }))

    case 'prompt': {
      const answer = await ctx.ask('prompt', {
        message: shortText(args.message, 'a prompt'),
        title: shortText(args.title, '', true),
        value: args.value === undefined || args.value === null ? '' : String(args.value),
      })
      return answer === null || answer === undefined ? null : String(answer)
    }

    case 'title':
      ctx.setTitle(shortText(args.text, '', true))
      return true

    case 'records.open': {
      if (!args.recordID || !/^[0-9]+$/.test(String(args.recordID))) {
        throw new Error('records.open needs the recordID of the record to open')
      }
      moduleIDFor(ctx, args.module)
      await ctx.openRecord({
        module: args.module,
        recordID: String(args.recordID),
        edit: !!args.edit,
      })
      return true
    }

    case 'records.delete': {
      const refusal = allowDelete(ctx.meta, args.module)
      if (refusal) throw new Error(refusal)
      const recordID = recordIDFor(args, 'records.delete removes one record; name it with recordID')
      const moduleID = moduleIDFor(ctx, args.module)
      if (!(await ctx.consentToDelete())) {
        throw new Error('the person using this app did not agree to it deleting records')
      }
      await ctx.compose.recordDelete({
        namespaceID: ctx.namespaceID,
        moduleID,
        recordID,
      })
      return true
    }

    case 'files.upload': {
      const moduleID = writableModuleIDFor(ctx, args.module)
      const { field, recordID } = fileFieldFor(ctx, args)
      const file = dataURLFile(args.dataURL, args.name)
      await agreed(ctx, args.module)

      const attachment = await ctx.uploadFile({ moduleID, recordID, field, file })

      // The upload stores the file; the record holds it only once the field
      // names it — beside what a multi-value field held, instead of what a
      // single one did.
      const spec = fieldsFor(ctx, moduleID).find(f => f.name === field)
      const multi = !!(spec?.isMulti ?? spec?.multi)
      const held = multi ? await fileIDs(ctx, moduleID, recordID, field) : []
      await ctx.compose.recordPatch({
        namespaceID: ctx.namespaceID,
        moduleID,
        recordID: [recordID],
        values: [...held, attachment.attachmentID].map(value => ({ name: field, value })),
      })
      return describeFile(attachment)
    }

    case 'users.search': {
      const query = String(args.query || '').trim()
      if (query.length < 2)
        throw new Error('users.search needs at least two characters to look for')
      const limit = Math.min(Math.max(Number(args.limit) || 10, 1), MAX_USERS)
      return ctx.searchUsers(query, limit)
    }

    case 'chatbot.open':
    case 'chatbot.close': {
      const name = String(args.chatbot || '').trim()
      if (!name) throw new Error(`${op} needs the handle of a chatbot the app declared`)
      if (!(ctx.meta?.chatbots || []).includes(name)) {
        throw new Error(`chatbot "${name}" is not declared for this app`)
      }
      if (op === 'chatbot.open') await ctx.openChatbot(name)
      else ctx.closeChatbot(name)
      return true
    }

    case 'navigate':
      await ctx.navigate(navigationTarget(args))
      return true

    case 'modules': {
      // What the app declared, as it stands now: the labels and the options a
      // Select holds today, rather than the ones its author copied in.
      return (ctx.meta?.modules || []).map(handle => {
        const moduleID = ctx.moduleIDs?.[handle]
        return {
          handle,
          moduleID,
          writable: (ctx.meta?.writes || []).includes(handle),
          fields: fieldsFor(ctx, moduleID).map(describeField),
        }
      })
    }

    case 'user':
      return ctx.user()

    case 'theme':
      return ctx.theme()

    case 'context':
      return ctx.context ? ctx.context() : {}

    case 'resize':
      ctx.resize(args.height)
      return true

    default:
      throw new Error(`operation "${op}" is not available to an app`)
  }
}
