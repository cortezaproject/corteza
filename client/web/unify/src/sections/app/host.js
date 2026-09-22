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

// The app's policy: no network of any kind, inline script and style only,
// images from data: URLs.
export const CSP_INNER =
  "default-src 'none'; " +
  "script-src 'unsafe-inline' https://cdnjs.cloudflare.com; " +
  "style-src 'unsafe-inline'; " +
  'img-src data:; ' +
  "connect-src 'none'; " +
  "form-action 'none'"

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

// The most records one call may ask for.
export const MAX_LIMIT = 500

// The bridge contract a page gets when it names none. See app.intent.md, Bridge.
export const BRIDGE_VERSION = 2

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

  var frame = document.createElement('iframe')
  frame.setAttribute('sandbox', 'allow-scripts')
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
      if (e.origin !== SHELL || data.type !== 'human:result') return
      var waiting = pending[data.id]
      if (!waiting) return
      delete pending[data.id]
      waiting.port.postMessage({ id: waiting.id, result: data.result, error: data.error })
      return
    }

    if (e.source !== frame.contentWindow || data.type !== 'human:hello') return

    var channel = new MessageChannel()
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
export async function dispatch(op, args = {}, ctx) {
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
        recordID: args.recordID,
      })
      const fields = fieldsFor(ctx, moduleID)
      return {
        record: reshapeRecord(record, fields, ctx.version),
        refs: await refsFor(ctx, fields, [record]),
      }
    }

    case 'records.report': {
      const moduleID = moduleIDFor(ctx, args.module)
      return ctx.compose.recordReport({
        namespaceID: ctx.namespaceID,
        moduleID,
        metrics: args.metrics,
        dimensions: args.dimensions,
        filter: args.filter,
      })
    }

    case 'user':
      return ctx.user()

    case 'theme':
      return ctx.theme()

    case 'resize':
      ctx.resize(args.height)
      return true

    default:
      throw new Error(`operation "${op}" is not available to an app`)
  }
}
